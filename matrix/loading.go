package matrix

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"log/slog"
	"os"
)

func LoadFileToMatrix(imagePath string) (*Matrix, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	img, format, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}
	slog.Info("Image loaded successfully", "format", format)

	return NewImageMatrix(img), nil
}

func NewImageMatrix(img image.Image) *Matrix {
	bounds := img.Bounds()
	m := NewMatrix(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			m.Set(x, y, color.RGBA{
				R: uint8(r >> 8),
				G: uint8(g >> 8),
				B: uint8(b >> 8),
				A: uint8(a >> 8),
			})
		}
	}
	return m
}

// Creates empty matrix starting at (0,0)
func NewEmptyMatrix(width, height int) *Matrix {
	return NewMatrixFromDimensions(width, height)
}

// Creates solid color matrix starting at (0,0)
func NewSolidColorMatrix(width, height int, col color.RGBA) *Matrix {
	m := NewMatrixFromDimensions(width, height)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			m.Set(x, y, col)
		}
	}
	return m
}

// Creates empty matrix with custom bounds
func NewEmptyMatrixWithBounds(bounds image.Rectangle) *Matrix {
	return NewMatrix(bounds)
}

// Creates solid color matrix with custom bounds
func NewSolidColorMatrixWithBounds(bounds image.Rectangle, col color.RGBA) *Matrix {
	m := NewMatrix(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			m.Set(x, y, col)
		}
	}
	return m
}

func NewMatrix(bounds image.Rectangle) *Matrix {
	width := bounds.Dx()
	height := bounds.Dy()
	mat := make([][]color.RGBA, height)
	for i := range mat {
		mat[i] = make([]color.RGBA, width)
	}
	return &Matrix{mat: mat, bounds: bounds}
}

// Convenience constructor for matrices starting at (0,0)
func NewMatrixFromDimensions(width, height int) *Matrix {
	return NewMatrix(image.Rect(0, 0, width, height))
}
