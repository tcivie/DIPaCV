package matrix

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// Convert Matrix back to image.Image
func (m *Matrix) ToImage() *image.RGBA {
	return image.NewRGBA(m.bounds)
}

// Save matrix as image file
func (m *Matrix) SaveToFile(filename string, quality ...int) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	img := m.ToImage()

	// Copy matrix data to image
	for y := m.bounds.Min.Y; y < m.bounds.Max.Y; y++ {
		for x := m.bounds.Min.X; x < m.bounds.Max.X; x++ {
			img.SetRGBA(x, y, m.At(x, y))
		}
	}

	ext := strings.ToLower(filepath.Ext(filename))

	switch ext {
	case ".png":
		return png.Encode(file, img)
	case ".jpg", ".jpeg":
		q := 95
		if len(quality) > 0 && quality[0] > 0 && quality[0] <= 100 {
			q = quality[0]
		}
		return jpeg.Encode(file, img, &jpeg.Options{Quality: q})
	default:
		return fmt.Errorf("unsupported format: %s (use .png or .jpg)", ext)
	}
}

// Display matrix info
func (m *Matrix) Display() string {
	return fmt.Sprintf("Matrix [%dx%d] @ (%d,%d)",
		m.Dx(), m.Dy(),
		m.bounds.Min.X, m.bounds.Min.Y)
}

// Get a small preview of the matrix
func (m *Matrix) Preview(maxSize int) string {
	if maxSize <= 0 {
		maxSize = 5
	}

	var preview strings.Builder
	preview.WriteString(m.Display() + "\n")

	rowsToShow := min(m.Dy(), maxSize)
	colsToShow := min(m.Dx(), maxSize)

	startY := m.bounds.Min.Y
	startX := m.bounds.Min.X

	for y := 0; y < rowsToShow; y++ {
		for x := 0; x < colsToShow; x++ {
			c := m.At(startX+x, startY+y)
			preview.WriteString(fmt.Sprintf("(%3d,%3d,%3d) ", c.R, c.G, c.B))
		}
		if m.Dx() > maxSize {
			preview.WriteString("...")
		}
		preview.WriteString("\n")
	}

	if m.Dy() > maxSize {
		preview.WriteString("...\n")
	}

	return preview.String()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Show opens the image in the system's default image viewer
func (m *Matrix) Show() error {
	tmpfile, err := os.CreateTemp("", "matrix-*.png")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpname := tmpfile.Name()
	tmpfile.Close()

	if err := m.SaveToFile(tmpname); err != nil {
		os.Remove(tmpname)
		return fmt.Errorf("failed to save temp file: %w", err)
	}

	return openFile(tmpname)
}

func openFile(filename string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", filename)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", filename)
	case "linux":
		if isCommandAvailable("xdg-open") {
			cmd = exec.Command("xdg-open", filename)
		} else if isCommandAvailable("gnome-open") {
			cmd = exec.Command("gnome-open", filename)
		} else if isCommandAvailable("kde-open") {
			cmd = exec.Command("kde-open", filename)
		} else {
			return fmt.Errorf("no suitable image viewer found on Linux")
		}
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}

	fmt.Printf("Opening image: %s\n", filename)
	return nil
}

func isCommandAvailable(name string) bool {
	cmd := exec.Command("which", name)
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

func (m *Matrix) ShowAndCleanup() error {
	tmpfile, err := os.CreateTemp("", "matrix-*.png")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpname := tmpfile.Name()
	tmpfile.Close()

	defer os.Remove(tmpname)

	if err := m.SaveToFile(tmpname); err != nil {
		return fmt.Errorf("failed to save temp file: %w", err)
	}

	if err := openFile(tmpname); err != nil {
		return err
	}

	fmt.Println("Press Enter to clean up the temporary file...")
	fmt.Scanln()

	return nil
}
