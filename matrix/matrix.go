package matrix

import (
	"image"
	"image/color"
)

type Matrix struct {
	mat    [][]color.RGBA
	bounds image.Rectangle
}

// Getters
func (m *Matrix) Dx() int                 { return m.bounds.Dx() }
func (m *Matrix) Dy() int                 { return m.bounds.Dy() }
func (m *Matrix) Dims() (int, int)        { return m.Dx(), m.Dy() }
func (m *Matrix) Bounds() image.Rectangle { return m.bounds }

// At returns color at absolute coordinates
func (m *Matrix) At(x, y int) color.RGBA {
	if !image.Pt(x, y).In(m.bounds) {
		return color.RGBA{}
	}
	matX := x - m.bounds.Min.X
	matY := y - m.bounds.Min.Y
	return m.mat[matY][matX]
}

// Set sets color at absolute coordinates
func (m *Matrix) Set(x, y int, c color.RGBA) {
	if !image.Pt(x, y).In(m.bounds) {
		return
	}
	matX := x - m.bounds.Min.X
	matY := y - m.bounds.Min.Y
	m.mat[matY][matX] = c
}

// Internal helpers
func (m *Matrix) applyWithOtherMatrix(other *Matrix, op func(c1, c2 color.RGBA) color.RGBA) {
	intersection := m.bounds.Intersect(other.bounds)
	if intersection.Empty() {
		return
	}
	for y := intersection.Min.Y; y < intersection.Max.Y; y++ {
		for x := intersection.Min.X; x < intersection.Max.X; x++ {
			c1 := m.At(x, y)
			c2 := other.At(x, y)
			m.Set(x, y, op(c1, c2))
		}
	}
}

func (m *Matrix) applyWithOtherScalar(scalar uint8, op func(c1 color.RGBA, s uint8) color.RGBA) {
	for y := m.bounds.Min.Y; y < m.bounds.Max.Y; y++ {
		for x := m.bounds.Min.X; x < m.bounds.Max.X; x++ {
			c := m.At(x, y)
			m.Set(x, y, op(c, scalar))
		}
	}
}

func (m *Matrix) checkWithOther(other *Matrix, op func(c1, c2 color.RGBA) bool) bool {
	intersection := m.bounds.Intersect(other.bounds)
	if intersection.Empty() {
		return true
	}

	for y := intersection.Min.Y; y < intersection.Max.Y; y++ {
		for x := intersection.Min.X; x < intersection.Max.X; x++ {
			if !op(m.At(x, y), other.At(x, y)) {
				return false
			}
		}
	}
	return true
}

// Arithmetic operations
func (m *Matrix) Add(other *Matrix) {
	m.applyWithOtherMatrix(other, func(c1, c2 color.RGBA) color.RGBA {
		return color.RGBA{
			R: clamp32(uint32(c1.R) + uint32(c2.R)),
			G: clamp32(uint32(c1.G) + uint32(c2.G)),
			B: clamp32(uint32(c1.B) + uint32(c2.B)),
			A: clamp32(uint32(c1.A) + uint32(c2.A)),
		}
	})
}

func (m *Matrix) Mul(scalar uint8) {
	m.applyWithOtherScalar(scalar, func(c color.RGBA, s uint8) color.RGBA {
		return color.RGBA{
			R: clamp32(uint32(c.R) * uint32(s)),
			G: clamp32(uint32(c.G) * uint32(s)),
			B: clamp32(uint32(c.B) * uint32(s)),
			A: clamp32(uint32(c.A) * uint32(s)),
		}
	})
}

func (m *Matrix) Div(scalar uint8) {
	if scalar == 0 {
		return // avoid divide by zero
	}
	m.applyWithOtherScalar(scalar, func(c color.RGBA, s uint8) color.RGBA {
		return color.RGBA{
			R: uint8(uint32(c.R) / uint32(s)),
			G: uint8(uint32(c.G) / uint32(s)),
			B: uint8(uint32(c.B) / uint32(s)),
			A: uint8(uint32(c.A) / uint32(s)),
		}
	})
}

func (m *Matrix) Row(row int) []color.RGBA { return m.mat[row] }

func (m *Matrix) Col(col int) []color.RGBA {
	colSlice := make([]color.RGBA, m.bounds.Dy())
	for y := m.bounds.Min.Y; y < m.bounds.Max.Y; y++ {
		colSlice[y-m.bounds.Min.Y] = m.At(col, y)
	}
	return colSlice
}

// Matrix multiplication (in-place)
func (m *Matrix) MulMat(other *Matrix) *Matrix {
	if !checkMulBounds(m, other) {
		panic("Invalid matrix dimensions for multiplication")
	}

	resBounds := image.Rect(0, 0, other.Dx(), m.Dy())
	result := NewMatrix(resBounds)

	for i := 0; i < m.Dy(); i++ {
		for j := 0; j < other.Dx(); j++ {
			var sumR, sumG, sumB, sumA uint32

			for k := 0; k < m.Dx(); k++ {
				a := m.mat[i][k]
				b := other.mat[k][j]

				sumR += uint32(a.R) * uint32(b.R)
				sumG += uint32(a.G) * uint32(b.G)
				sumB += uint32(a.B) * uint32(b.B)
				sumA += uint32(a.A) * uint32(b.A)
			}

			// Normalize (approximate) and clamp
			result.mat[i][j] = color.RGBA{
				R: clamp32(sumR / 255),
				G: clamp32(sumG / 255),
				B: clamp32(sumB / 255),
				A: clamp32(sumA / 255),
			}
		}
	}

	m.bounds = resBounds
	m.mat = result.mat
	return m
}

// Set operations
func (m *Matrix) Union(other *Matrix) *Matrix {
	unionBounds := m.bounds.Union(other.bounds)
	result := NewMatrix(unionBounds)

	for y := m.bounds.Min.Y; y < m.bounds.Max.Y; y++ {
		for x := m.bounds.Min.X; x < m.bounds.Max.X; x++ {
			result.Set(x, y, m.At(x, y))
		}
	}
	for y := other.bounds.Min.Y; y < other.bounds.Max.Y; y++ {
		for x := other.bounds.Min.X; x < other.bounds.Max.X; x++ {
			c := other.At(x, y)
			if c.A > 0 {
				result.Set(x, y, c)
			}
		}
	}
	return result
}

func (m *Matrix) Intersection(other *Matrix) *Matrix {
	intersection := m.bounds.Intersect(other.bounds)
	if intersection.Empty() {
		return NewMatrix(image.Rectangle{})
	}

	result := NewMatrix(intersection)
	for y := intersection.Min.Y; y < intersection.Max.Y; y++ {
		for x := intersection.Min.X; x < intersection.Max.X; x++ {
			c1 := m.At(x, y)
			c2 := other.At(x, y)
			result.Set(x, y, color.RGBA{
				R: uint8((uint32(c1.R) + uint32(c2.R)) / 2),
				G: uint8((uint32(c1.G) + uint32(c2.G)) / 2),
				B: uint8((uint32(c1.B) + uint32(c2.B)) / 2),
				A: uint8((uint32(c1.A) + uint32(c2.A)) / 2),
			})
		}
	}
	return result
}

func (m *Matrix) Inclusion(other *Matrix) bool {
	if !m.bounds.In(other.bounds) {
		return false
	}
	return m.checkWithOther(other, func(c1, c2 color.RGBA) bool {
		return c1 == c2
	})
}

// Helper functions
func clamp32(v uint32) uint8 {
	if v > 255 {
		return 255
	}
	return uint8(v)
}

func checkMulBounds(m1, m2 *Matrix) bool {
	// For A(m×n) * B(n×p), A.Dx()==n, B.Dy()==n
	return m1.Dx() == m2.Dy()
}

func clamp(v uint8) uint8 {
	if v > 255 {
		return 255
	} else if v < 0 {
		return 0
	}
	return v
}
