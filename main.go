package main

import (
	"DIPaCV/matrix"
	"image/color"
)

func main() {
	im := matrix.NewSolidColorMatrix(100, 100, color.RGBA{R: 255, G: 0, B: 0, A: 255})

	if im == nil {
		return
	}

	im.Add(matrix.NewSolidColorMatrix(im.Dx(), im.Dy(), color.RGBA{R: 0, G: 255, B: 0, A: 255}))
	im.Show()
}
