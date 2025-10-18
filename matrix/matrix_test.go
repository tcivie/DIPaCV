package matrix

import (
	"image"
	"image/color"
	"testing"
)

// Test basic matrix creation and dimensions
func TestMatrixCreation(t *testing.T) {
	tests := []struct {
		name   string
		width  int
		height int
	}{
		{"Small matrix", 10, 10},
		{"Wide matrix", 100, 10},
		{"Tall matrix", 10, 100},
		{"Large matrix", 500, 500},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewEmptyMatrix(tt.width, tt.height)
			if m == nil {
				t.Fatal("NewEmptyMatrix returned nil")
			}
			if m.Dx() != tt.width {
				t.Errorf("Expected width %d, got %d", tt.width, m.Dx())
			}
			if m.Dy() != tt.height {
				t.Errorf("Expected height %d, got %d", tt.height, m.Dy())
			}
		})
	}
}

// Test matrix with custom bounds
func TestMatrixWithBounds(t *testing.T) {
	bounds := image.Rect(10, 20, 50, 70)
	m := NewEmptyMatrixWithBounds(bounds)

	if m.Dx() != 40 {
		t.Errorf("Expected width 40, got %d", m.Dx())
	}
	if m.Dy() != 50 {
		t.Errorf("Expected height 50, got %d", m.Dy())
	}
	if m.Bounds() != bounds {
		t.Errorf("Expected bounds %v, got %v", bounds, m.Bounds())
	}
}

// Test solid color matrix
func TestSolidColorMatrix(t *testing.T) {
	red := color.RGBA{R: 255, G: 0, B: 0, A: 255}
	m := NewSolidColorMatrix(10, 10, red)

	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			c := m.At(x, y)
			if c != red {
				t.Errorf("Expected color %v at (%d,%d), got %v", red, x, y, c)
			}
		}
	}
}

// Test Get/Set operations
func TestGetSet(t *testing.T) {
	m := NewEmptyMatrix(5, 5)
	testColor := color.RGBA{R: 100, G: 150, B: 200, A: 255}

	m.Set(2, 3, testColor)
	got := m.At(2, 3)

	if got != testColor {
		t.Errorf("Expected %v, got %v", testColor, got)
	}
}

// Test Add operation with overflow protection
func TestAdd(t *testing.T) {
	tests := []struct {
		name     string
		c1       color.RGBA
		c2       color.RGBA
		expected color.RGBA
	}{
		{
			name:     "Normal addition",
			c1:       color.RGBA{R: 100, G: 50, B: 25, A: 200},
			c2:       color.RGBA{R: 50, G: 100, B: 50, A: 50},
			expected: color.RGBA{R: 150, G: 150, B: 75, A: 250},
		},
		{
			name:     "Overflow red channel",
			c1:       color.RGBA{R: 255, G: 0, B: 0, A: 255},
			c2:       color.RGBA{R: 100, G: 0, B: 0, A: 0},
			expected: color.RGBA{R: 255, G: 0, B: 0, A: 255},
		},
		{
			name:     "Overflow all channels",
			c1:       color.RGBA{R: 200, G: 200, B: 200, A: 200},
			c2:       color.RGBA{R: 100, G: 100, B: 100, A: 100},
			expected: color.RGBA{R: 255, G: 255, B: 255, A: 255},
		},
		{
			name:     "Red + Green = Yellow",
			c1:       color.RGBA{R: 255, G: 0, B: 0, A: 255},
			c2:       color.RGBA{R: 0, G: 255, B: 0, A: 0},
			expected: color.RGBA{R: 255, G: 255, B: 0, A: 255},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m1 := NewSolidColorMatrix(10, 10, tt.c1)
			m2 := NewSolidColorMatrix(10, 10, tt.c2)

			m1.Add(m2)

			got := m1.At(0, 0)
			if got != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, got)
			}

			// Verify all pixels have the same color
			for y := 0; y < 10; y++ {
				for x := 0; x < 10; x++ {
					if m1.At(x, y) != tt.expected {
						t.Errorf("Pixel (%d,%d): expected %v, got %v", x, y, tt.expected, m1.At(x, y))
					}
				}
			}
		})
	}
}

// Test Mul operation with scalar
func TestMul(t *testing.T) {
	tests := []struct {
		name     string
		color    color.RGBA
		scalar   uint8
		expected color.RGBA
	}{
		{
			name:     "Double intensity",
			color:    color.RGBA{R: 100, G: 50, B: 25, A: 128},
			scalar:   2,
			expected: color.RGBA{R: 200, G: 100, B: 50, A: 255},
		},
		{
			name:     "Half intensity",
			color:    color.RGBA{R: 200, G: 100, B: 50, A: 200},
			scalar:   1, // Divide by 2 in uint8 math
			expected: color.RGBA{R: 200, G: 100, B: 50, A: 200},
		},
		{
			name:     "Overflow protection",
			color:    color.RGBA{R: 200, G: 200, B: 200, A: 200},
			scalar:   2,
			expected: color.RGBA{R: 255, G: 255, B: 255, A: 255},
		},
		{
			name:     "Zero scalar",
			color:    color.RGBA{R: 255, G: 255, B: 255, A: 255},
			scalar:   0,
			expected: color.RGBA{R: 0, G: 0, B: 0, A: 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewSolidColorMatrix(10, 10, tt.color)
			m.Mul(tt.scalar)

			got := m.At(0, 0)
			if got != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, got)
			}
		})
	}
}

// Test Div operation with scalar
func TestDiv(t *testing.T) {
	tests := []struct {
		name     string
		color    color.RGBA
		scalar   uint8
		expected color.RGBA
	}{
		{
			name:     "Divide by 2",
			color:    color.RGBA{R: 200, G: 100, B: 50, A: 200},
			scalar:   2,
			expected: color.RGBA{R: 100, G: 50, B: 25, A: 100},
		},
		{
			name:     "Divide by zero (no-op)",
			color:    color.RGBA{R: 200, G: 100, B: 50, A: 200},
			scalar:   0,
			expected: color.RGBA{R: 200, G: 100, B: 50, A: 200},
		},
		{
			name:     "Divide by 1 (identity)",
			color:    color.RGBA{R: 255, G: 128, B: 64, A: 255},
			scalar:   1,
			expected: color.RGBA{R: 255, G: 128, B: 64, A: 255},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewSolidColorMatrix(10, 10, tt.color)
			m.Div(tt.scalar)

			got := m.At(0, 0)
			if got != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, got)
			}
		})
	}
}

// Test Union operation
func TestUnion(t *testing.T) {
	red := color.RGBA{R: 255, G: 0, B: 0, A: 255}
	blue := color.RGBA{R: 0, G: 0, B: 255, A: 255}

	// Create two matrices at different positions
	m1 := NewSolidColorMatrixWithBounds(image.Rect(0, 0, 10, 10), red)
	m2 := NewSolidColorMatrixWithBounds(image.Rect(5, 5, 15, 15), blue)

	result := m1.Union(m2)

	// Check union bounds
	expectedBounds := image.Rect(0, 0, 15, 15)
	if result.Bounds() != expectedBounds {
		t.Errorf("Expected bounds %v, got %v", expectedBounds, result.Bounds())
	}

	// Check that red region has red color
	if result.At(0, 0) != red {
		t.Errorf("Expected red at (0,0), got %v", result.At(0, 0))
	}

	// Check that blue region has blue color (overlapping region will be blue due to overwrite)
	if result.At(10, 10) != blue {
		t.Errorf("Expected blue at (10,10), got %v", result.At(10, 10))
	}
}

// Test Intersection operation
func TestIntersection(t *testing.T) {
	red := color.RGBA{R: 255, G: 0, B: 0, A: 255}
	blue := color.RGBA{R: 0, G: 0, B: 255, A: 255}

	m1 := NewSolidColorMatrixWithBounds(image.Rect(0, 0, 10, 10), red)
	m2 := NewSolidColorMatrixWithBounds(image.Rect(5, 5, 15, 15), blue)

	result := m1.Intersection(m2)

	// Check intersection bounds
	expectedBounds := image.Rect(5, 5, 10, 10)
	if result.Bounds() != expectedBounds {
		t.Errorf("Expected bounds %v, got %v", expectedBounds, result.Bounds())
	}

	// Check that color is averaged
	expectedColor := color.RGBA{R: 127, G: 0, B: 127, A: 255}
	got := result.At(7, 7)
	if got != expectedColor {
		t.Errorf("Expected averaged color %v at (7,7), got %v", expectedColor, got)
	}
}

// Test Intersection with no overlap
func TestIntersectionEmpty(t *testing.T) {
	red := color.RGBA{R: 255, G: 0, B: 0, A: 255}
	blue := color.RGBA{R: 0, G: 0, B: 255, A: 255}

	m1 := NewSolidColorMatrixWithBounds(image.Rect(0, 0, 10, 10), red)
	m2 := NewSolidColorMatrixWithBounds(image.Rect(20, 20, 30, 30), blue)

	result := m1.Intersection(m2)

	if !result.Bounds().Empty() {
		t.Errorf("Expected empty intersection, got bounds %v", result.Bounds())
	}
}

// Test Matrix multiplication
func TestMulMat(t *testing.T) {
	// Create simple 2x2 matrices
	m1 := NewMatrixFromDimensions(2, 2)
	m1.Set(0, 0, color.RGBA{R: 1, G: 2, B: 3, A: 255})
	m1.Set(1, 0, color.RGBA{R: 4, G: 5, B: 6, A: 255})
	m1.Set(0, 1, color.RGBA{R: 7, G: 8, B: 9, A: 255})
	m1.Set(1, 1, color.RGBA{R: 10, G: 11, B: 12, A: 255})

	m2 := NewMatrixFromDimensions(2, 2)
	m2.Set(0, 0, color.RGBA{R: 2, G: 0, B: 1, A: 255})
	m2.Set(1, 0, color.RGBA{R: 1, G: 3, B: 2, A: 255})
	m2.Set(0, 1, color.RGBA{R: 0, G: 1, B: 0, A: 255})
	m2.Set(1, 1, color.RGBA{R: 4, G: 2, B: 3, A: 255})

	result := m1.MulMat(m2)

	// Just verify dimensions and that it doesn't crash
	if result.Dx() != 2 || result.Dy() != 2 {
		t.Errorf("Expected 2x2 result, got %dx%d", result.Dx(), result.Dy())
	}
}

// Test bounds checking
func TestBoundsChecking(t *testing.T) {
	m := NewMatrixFromDimensions(10, 10)
	testColor := color.RGBA{R: 255, G: 0, B: 0, A: 255}

	// Out of bounds access should return zero color
	got := m.At(20, 20)
	expected := color.RGBA{}
	if got != expected {
		t.Errorf("Expected zero color for out-of-bounds access, got %v", got)
	}

	// Out of bounds set should not panic
	m.Set(20, 20, testColor)

	// Verify the set didn't affect the matrix
	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			if m.At(x, y) == testColor {
				t.Errorf("Out-of-bounds Set affected matrix at (%d,%d)", x, y)
			}
		}
	}
}

// Test Row and Col accessors
func TestRowCol(t *testing.T) {
	m := NewMatrixFromDimensions(3, 3)
	red := color.RGBA{R: 255, G: 0, B: 0, A: 255}

	// Set middle row to red
	for x := 0; x < 3; x++ {
		m.Set(x, 1, red)
	}

	row := m.Row(1)
	if len(row) != 3 {
		t.Errorf("Expected row length 3, got %d", len(row))
	}

	for _, c := range row {
		if c != red {
			t.Errorf("Expected red in row, got %v", c)
		}
	}

	// Set middle column to red
	m2 := NewMatrixFromDimensions(3, 3)
	for y := 0; y < 3; y++ {
		m2.Set(1, y, red)
	}

	col := m2.Col(1)
	if len(col) != 3 {
		t.Errorf("Expected col length 3, got %d", len(col))
	}

	for _, c := range col {
		if c != red {
			t.Errorf("Expected red in col, got %v", c)
		}
	}
}

// ============================================================================
// Benchmark Tests
// ============================================================================

// Benchmark matrix creation for different sizes
func BenchmarkMatrixCreation_Small(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewEmptyMatrix(10, 10)
	}
}

func BenchmarkMatrixCreation_Medium(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewEmptyMatrix(100, 100)
	}
}

func BenchmarkMatrixCreation_Large(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewEmptyMatrix(1000, 1000)
	}
}

func BenchmarkMatrixCreation_HD(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewEmptyMatrix(1920, 1080)
	}
}

// Benchmark solid color matrix creation
func BenchmarkSolidColorMatrix_Small(b *testing.B) {
	c := color.RGBA{R: 128, G: 128, B: 128, A: 255}
	for i := 0; i < b.N; i++ {
		_ = NewSolidColorMatrix(10, 10, c)
	}
}

func BenchmarkSolidColorMatrix_Medium(b *testing.B) {
	c := color.RGBA{R: 128, G: 128, B: 128, A: 255}
	for i := 0; i < b.N; i++ {
		_ = NewSolidColorMatrix(100, 100, c)
	}
}

func BenchmarkSolidColorMatrix_Large(b *testing.B) {
	c := color.RGBA{R: 128, G: 128, B: 128, A: 255}
	for i := 0; i < b.N; i++ {
		_ = NewSolidColorMatrix(1000, 1000, c)
	}
}

// Benchmark Add operation for different sizes
func BenchmarkAdd_Small(b *testing.B) {
	m1 := NewSolidColorMatrix(10, 10, color.RGBA{R: 100, G: 100, B: 100, A: 255})
	m2 := NewSolidColorMatrix(10, 10, color.RGBA{R: 50, G: 50, B: 50, A: 255})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m1.Add(m2)
	}
}

func BenchmarkAdd_Medium(b *testing.B) {
	m1 := NewSolidColorMatrix(100, 100, color.RGBA{R: 100, G: 100, B: 100, A: 255})
	m2 := NewSolidColorMatrix(100, 100, color.RGBA{R: 50, G: 50, B: 50, A: 255})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m1.Add(m2)
	}
}

func BenchmarkAdd_Large(b *testing.B) {
	m1 := NewSolidColorMatrix(500, 500, color.RGBA{R: 100, G: 100, B: 100, A: 255})
	m2 := NewSolidColorMatrix(500, 500, color.RGBA{R: 50, G: 50, B: 50, A: 255})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m1.Add(m2)
	}
}

func BenchmarkAdd_HD(b *testing.B) {
	m1 := NewSolidColorMatrix(1920, 1080, color.RGBA{R: 100, G: 100, B: 100, A: 255})
	m2 := NewSolidColorMatrix(1920, 1080, color.RGBA{R: 50, G: 50, B: 50, A: 255})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m1.Add(m2)
	}
}

// Benchmark Mul operation for different sizes
func BenchmarkMul_Small(b *testing.B) {
	m := NewSolidColorMatrix(10, 10, color.RGBA{R: 100, G: 100, B: 100, A: 255})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Mul(2)
	}
}

func BenchmarkMul_Medium(b *testing.B) {
	m := NewSolidColorMatrix(100, 100, color.RGBA{R: 100, G: 100, B: 100, A: 255})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Mul(2)
	}
}

func BenchmarkMul_Large(b *testing.B) {
	m := NewSolidColorMatrix(500, 500, color.RGBA{R: 100, G: 100, B: 100, A: 255})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Mul(2)
	}
}

func BenchmarkMul_HD(b *testing.B) {
	m := NewSolidColorMatrix(1920, 1080, color.RGBA{R: 100, G: 100, B: 100, A: 255})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Mul(2)
	}
}

// Benchmark Div operation for different sizes
func BenchmarkDiv_Small(b *testing.B) {
	m := NewSolidColorMatrix(10, 10, color.RGBA{R: 200, G: 200, B: 200, A: 255})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Div(2)
	}
}

func BenchmarkDiv_Medium(b *testing.B) {
	m := NewSolidColorMatrix(100, 100, color.RGBA{R: 200, G: 200, B: 200, A: 255})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Div(2)
	}
}

func BenchmarkDiv_Large(b *testing.B) {
	m := NewSolidColorMatrix(500, 500, color.RGBA{R: 200, G: 200, B: 200, A: 255})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Div(2)
	}
}

func BenchmarkDiv_HD(b *testing.B) {
	m := NewSolidColorMatrix(1920, 1080, color.RGBA{R: 200, G: 200, B: 200, A: 255})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Div(2)
	}
}

// Benchmark Matrix multiplication
func BenchmarkMulMat_2x2(b *testing.B) {
	m1 := NewMatrixFromDimensions(2, 2)
	m2 := NewMatrixFromDimensions(2, 2)

	// Fill with some values
	for y := 0; y < 2; y++ {
		for x := 0; x < 2; x++ {
			m1.Set(x, y, color.RGBA{R: 1, G: 2, B: 3, A: 255})
			m2.Set(x, y, color.RGBA{R: 2, G: 3, B: 4, A: 255})
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m1.MulMat(m2)
	}
}

func BenchmarkMulMat_10x10(b *testing.B) {
	m1 := NewMatrixFromDimensions(10, 10)
	m2 := NewMatrixFromDimensions(10, 10)

	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			m1.Set(x, y, color.RGBA{R: 1, G: 2, B: 3, A: 255})
			m2.Set(x, y, color.RGBA{R: 2, G: 3, B: 4, A: 255})
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m1.MulMat(m2)
	}
}

func BenchmarkMulMat_50x50(b *testing.B) {
	m1 := NewMatrixFromDimensions(50, 50)
	m2 := NewMatrixFromDimensions(50, 50)

	for y := 0; y < 50; y++ {
		for x := 0; x < 50; x++ {
			m1.Set(x, y, color.RGBA{R: 1, G: 2, B: 3, A: 255})
			m2.Set(x, y, color.RGBA{R: 2, G: 3, B: 4, A: 255})
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m1.MulMat(m2)
	}
}

// Benchmark Union operation
func BenchmarkUnion_Small(b *testing.B) {
	m1 := NewSolidColorMatrixWithBounds(image.Rect(0, 0, 10, 10), color.RGBA{R: 255, G: 0, B: 0, A: 255})
	m2 := NewSolidColorMatrixWithBounds(image.Rect(5, 5, 15, 15), color.RGBA{R: 0, G: 255, B: 0, A: 255})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m1.Union(m2)
	}
}

func BenchmarkUnion_Medium(b *testing.B) {
	m1 := NewSolidColorMatrixWithBounds(image.Rect(0, 0, 100, 100), color.RGBA{R: 255, G: 0, B: 0, A: 255})
	m2 := NewSolidColorMatrixWithBounds(image.Rect(50, 50, 150, 150), color.RGBA{R: 0, G: 255, B: 0, A: 255})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m1.Union(m2)
	}
}

func BenchmarkUnion_Large(b *testing.B) {
	m1 := NewSolidColorMatrixWithBounds(image.Rect(0, 0, 500, 500), color.RGBA{R: 255, G: 0, B: 0, A: 255})
	m2 := NewSolidColorMatrixWithBounds(image.Rect(250, 250, 750, 750), color.RGBA{R: 0, G: 255, B: 0, A: 255})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m1.Union(m2)
	}
}

// Benchmark Intersection operation
func BenchmarkIntersection_Small(b *testing.B) {
	m1 := NewSolidColorMatrixWithBounds(image.Rect(0, 0, 10, 10), color.RGBA{R: 255, G: 0, B: 0, A: 255})
	m2 := NewSolidColorMatrixWithBounds(image.Rect(5, 5, 15, 15), color.RGBA{R: 0, G: 255, B: 0, A: 255})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m1.Intersection(m2)
	}
}

func BenchmarkIntersection_Medium(b *testing.B) {
	m1 := NewSolidColorMatrixWithBounds(image.Rect(0, 0, 100, 100), color.RGBA{R: 255, G: 0, B: 0, A: 255})
	m2 := NewSolidColorMatrixWithBounds(image.Rect(50, 50, 150, 150), color.RGBA{R: 0, G: 255, B: 0, A: 255})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m1.Intersection(m2)
	}
}

func BenchmarkIntersection_Large(b *testing.B) {
	m1 := NewSolidColorMatrixWithBounds(image.Rect(0, 0, 500, 500), color.RGBA{R: 255, G: 0, B: 0, A: 255})
	m2 := NewSolidColorMatrixWithBounds(image.Rect(250, 250, 750, 750), color.RGBA{R: 0, G: 255, B: 0, A: 255})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m1.Intersection(m2)
	}
}

// Benchmark Get/Set operations
func BenchmarkGet(b *testing.B) {
	m := NewSolidColorMatrix(100, 100, color.RGBA{R: 128, G: 128, B: 128, A: 255})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m.At(50, 50)
	}
}

func BenchmarkSet(b *testing.B) {
	m := NewEmptyMatrix(100, 100)
	c := color.RGBA{R: 255, G: 0, B: 0, A: 255}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Set(50, 50, c)
	}
}

// Benchmark Row and Col accessors
func BenchmarkRow_Small(b *testing.B) {
	m := NewSolidColorMatrix(10, 10, color.RGBA{R: 128, G: 128, B: 128, A: 255})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m.Row(5)
	}
}

func BenchmarkRow_Medium(b *testing.B) {
	m := NewSolidColorMatrix(100, 100, color.RGBA{R: 128, G: 128, B: 128, A: 255})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m.Row(50)
	}
}

func BenchmarkRow_Large(b *testing.B) {
	m := NewSolidColorMatrix(1000, 1000, color.RGBA{R: 128, G: 128, B: 128, A: 255})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m.Row(500)
	}
}

func BenchmarkCol_Small(b *testing.B) {
	m := NewSolidColorMatrix(10, 10, color.RGBA{R: 128, G: 128, B: 128, A: 255})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m.Col(5)
	}
}

func BenchmarkCol_Medium(b *testing.B) {
	m := NewSolidColorMatrix(100, 100, color.RGBA{R: 128, G: 128, B: 128, A: 255})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m.Col(50)
	}
}

func BenchmarkCol_Large(b *testing.B) {
	m := NewSolidColorMatrix(1000, 1000, color.RGBA{R: 128, G: 128, B: 128, A: 255})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m.Col(500)
	}
}

// Benchmark ToImage conversion
func BenchmarkToImage_Small(b *testing.B) {
	m := NewSolidColorMatrix(10, 10, color.RGBA{R: 128, G: 128, B: 128, A: 255})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m.ToImage()
	}
}

func BenchmarkToImage_Medium(b *testing.B) {
	m := NewSolidColorMatrix(100, 100, color.RGBA{R: 128, G: 128, B: 128, A: 255})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m.ToImage()
	}
}

func BenchmarkToImage_Large(b *testing.B) {
	m := NewSolidColorMatrix(500, 500, color.RGBA{R: 128, G: 128, B: 128, A: 255})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m.ToImage()
	}
}

func BenchmarkToImage_HD(b *testing.B) {
	m := NewSolidColorMatrix(1920, 1080, color.RGBA{R: 128, G: 128, B: 128, A: 255})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m.ToImage()
	}
}
