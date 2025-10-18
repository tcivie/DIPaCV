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
	bounds := image.Rect(0, 0, m.Dx(), m.Dy())
	img := image.NewRGBA(bounds)

	for y := 0; y < m.Dy(); y++ {
		for x := 0; x < m.Dx(); x++ {
			img.SetRGBA(x, y, m.mat[y][x])
		}
	}

	return img
}

// Save matrix as image file
func (m *Matrix) SaveToFile(filename string, quality ...int) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	img := m.ToImage()
	ext := strings.ToLower(filepath.Ext(filename))

	switch ext {
	case ".png":
		return png.Encode(file, img)
	case ".jpg", ".jpeg":
		// Default quality is 95
		q := 95
		if len(quality) > 0 && quality[0] > 0 && quality[0] <= 100 {
			q = quality[0]
		}
		return jpeg.Encode(file, img, &jpeg.Options{Quality: q})
	default:
		return fmt.Errorf("unsupported format: %s (use .png or .jpg)", ext)
	}
}

// Display matrix info and optionally save a preview
func (m *Matrix) Display() string {
	return fmt.Sprintf("Matrix [%dx%d]", m.Dx(), m.Dy())
}

// Get a small preview of the matrix (for debugging)
func (m *Matrix) Preview(maxSize int) string {
	if maxSize <= 0 {
		maxSize = 5
	}

	var preview strings.Builder
	preview.WriteString(fmt.Sprintf("Matrix Preview [%dx%d]:\n", m.Dx(), m.Dy()))

	rowsToShow := min(m.Dy(), maxSize)
	colsToShow := min(m.Dx(), maxSize)

	for y := 0; y < rowsToShow; y++ {
		for x := 0; x < colsToShow; x++ {
			c := m.mat[y][x]
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

// Helper function
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Show opens the image in the system's default image viewer
func (m *Matrix) Show() error {
	// Create temp file
	tmpfile, err := os.CreateTemp("", "matrix-*.png")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpname := tmpfile.Name()
	tmpfile.Close()

	// Save to temp file
	if err := m.SaveToFile(tmpname); err != nil {
		os.Remove(tmpname)
		return fmt.Errorf("failed to save temp file: %w", err)
	}

	// Open with system default viewer
	return openFile(tmpname)
}

// openFile opens a file with the system's default application
func openFile(filename string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin": // macOS
		cmd = exec.Command("open", filename)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", filename)
	case "linux":
		// Try different commands that might be available
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

// Helper to check if a command exists
func isCommandAvailable(name string) bool {
	cmd := exec.Command("which", name)
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

// Alternative: Show with auto-cleanup after viewing
func (m *Matrix) ShowAndCleanup() error {
	tmpfile, err := os.CreateTemp("", "matrix-*.png")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpname := tmpfile.Name()
	tmpfile.Close()

	// Ensure cleanup happens
	defer os.Remove(tmpname)

	if err := m.SaveToFile(tmpname); err != nil {
		return fmt.Errorf("failed to save temp file: %w", err)
	}

	// Open and wait for the process to start
	if err := openFile(tmpname); err != nil {
		return err
	}

	// Give the viewer time to open the file
	// You might want to make this configurable
	fmt.Println("Press Enter to clean up the temporary file...")
	fmt.Scanln()

	return nil
}
