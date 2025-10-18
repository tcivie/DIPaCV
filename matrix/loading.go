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
	width := bounds.Dx()
	height := bounds.Dy()

	imageMatrix := make([][]color.RGBA, height)
	for y := 0; y < height; y++ {
		imageMatrix[y] = make([]color.RGBA, width)
		for x := 0; x < width; x++ {
			r, g, b, a := img.At(bounds.Min.X+x, bounds.Min.Y+y).RGBA()
			imageMatrix[y][x] = color.RGBA{
				R: uint8(r >> 8),
				G: uint8(g >> 8),
				B: uint8(b >> 8),
				A: uint8(a >> 8),
			}
		}
	}
	return &Matrix{mat: imageMatrix}
}

func NewEmptyMatrix(width, height int) *Matrix {
	imageMatrix := make([][]color.RGBA, height)
	for i := 0; i < height; i++ {
		imageMatrix[i] = make([]color.RGBA, width)
	}
	return &Matrix{mat: imageMatrix}
}

func NewSolidColorMatrix(width, height int, col color.RGBA) *Matrix {
	imageMatrix := make([][]color.RGBA, height)
	for i := 0; i < height; i++ {
		imageMatrix[i] = make([]color.RGBA, width)
		for j := 0; j < width; j++ {
			imageMatrix[i][j] = col
		}
	}
	return &Matrix{mat: imageMatrix}
}
