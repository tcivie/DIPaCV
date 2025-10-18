package matrix

import "image/color"

type Matrix struct {
	mat [][]color.RGBA
}

func (m *Matrix) applyWithOther(other *Matrix, op func(c1, c2 color.RGBA) color.RGBA) {
	for y := 0; y < m.Dy(); y++ {
		for x := 0; x < m.Dx(); x++ {
			m.mat[y][x] = op(m.mat[y][x], other.mat[y][x])
		}
	}
}

func (m *Matrix) Dx() int {
	if len(m.mat) == 0 {
		return 0
	}
	return len(m.mat[0])
}

func (m *Matrix) Dy() int {
	return len(m.mat)
}

func (m *Matrix) Dims() (int, int) {
	return m.Dx(), m.Dy()
}

func (m *Matrix) At(x, y int) color.RGBA {
	return m.mat[y][x]
}

func (m *Matrix) Set(x, y int, c color.RGBA) {
	m.mat[y][x] = c
}

// Arithmetic operations

func (m *Matrix) Add(other *Matrix) {
	m.applyWithOther(other, func(c1, c2 color.RGBA) color.RGBA {
		return color.RGBA{
			R: clampAdd(c1.R, c2.R),
			G: clampAdd(c1.G, c2.G),
			B: clampAdd(c1.B, c2.B),
			A: clampAdd(c1.A, c2.A),
		}
	})
}

// clampAdd adds two uint8 values and clamps to 255 to prevent overflow
func clampAdd(a, b uint8) uint8 {
	if int(a)+int(b) > 255 {
		return 255
	}
	return a + b
}

func (m *Matrix) Sub(other *Matrix) {
	m.applyWithOther(other, func(c1, c2 color.RGBA) color.RGBA {
		return color.RGBA{
			R: c1.R - c2.R,
			G: c1.G - c2.G,
			B: c1.B - c2.B,
			A: c1.A - c2.A,
		}
	})
}

func (m *Matrix) Mul(other *Matrix) {
	m.applyWithOther(other, func(c1, c2 color.RGBA) color.RGBA {
		return color.RGBA{
			R: c1.R * c2.R,
			G: c1.G * c2.G,
			B: c1.B * c2.B,
			A: c1.A * c2.A,
		}
	})
}

func (m *Matrix) Div(other *Matrix) {
	m.applyWithOther(other, func(c1, c2 color.RGBA) color.RGBA {
		return color.RGBA{
			R: c1.R / c2.R,
			G: c1.G / c2.G,
			B: c1.B / c2.B,
			A: c1.A / c2.A,
		}
	})
}
