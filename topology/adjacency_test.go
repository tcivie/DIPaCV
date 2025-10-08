package topology

import (
	"image"
	"testing"
)

func TestAre4Adjacent(t *testing.T) {
	tests := []struct {
		name string
		p    image.Point
		q    image.Point
		want bool
	}{
		{"right", image.Point{5, 5}, image.Point{6, 5}, true},
		{"left", image.Point{5, 5}, image.Point{4, 5}, true},
		{"up", image.Point{5, 5}, image.Point{5, 4}, true},
		{"down", image.Point{5, 5}, image.Point{5, 6}, true},
		{"diagonal", image.Point{5, 5}, image.Point{6, 6}, false},
		{"same point", image.Point{5, 5}, image.Point{5, 5}, false},
		{"two away", image.Point{5, 5}, image.Point{7, 5}, false},
		{"negative coords right", image.Point{-3, -3}, image.Point{-2, -3}, true},
		{"negative coords up", image.Point{-3, -3}, image.Point{-3, -4}, true},
		{"negative to positive", image.Point{-1, 0}, image.Point{0, 0}, true},
		{"negative diagonal", image.Point{-5, -5}, image.Point{-4, -4}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Are4Adjacent(tt.p, tt.q); got != tt.want {
				t.Errorf("Are4Adjacent(%v, %v) = %v, want %v", tt.p, tt.q, got, tt.want)
			}
		})
	}
}

func TestAre8Adjacent(t *testing.T) {
	tests := []struct {
		name string
		p    image.Point
		q    image.Point
		want bool
	}{
		{"right", image.Point{5, 5}, image.Point{6, 5}, true},
		{"diagonal", image.Point{5, 5}, image.Point{6, 6}, true},
		{"diagonal up-left", image.Point{5, 5}, image.Point{4, 4}, true},
		{"same point", image.Point{5, 5}, image.Point{5, 5}, false},
		{"two away", image.Point{5, 5}, image.Point{7, 5}, false},
		{"knight move", image.Point{5, 5}, image.Point{6, 7}, false},
		{"negative diagonal", image.Point{-2, -2}, image.Point{-1, -1}, true},
		{"negative right", image.Point{-5, -3}, image.Point{-4, -3}, true},
		{"cross zero", image.Point{-1, -1}, image.Point{0, 0}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Are8Adjacent(tt.p, tt.q); got != tt.want {
				t.Errorf("Are8Adjacent(%v, %v) = %v, want %v", tt.p, tt.q, got, tt.want)
			}
		})
	}
}

func TestAre24Adjacent(t *testing.T) {
	tests := []struct {
		name string
		p    image.Point
		q    image.Point
		want bool
	}{
		{"direct neighbor", image.Point{5, 5}, image.Point{6, 5}, true},
		{"diagonal", image.Point{5, 5}, image.Point{6, 6}, true},
		{"2nd level straight", image.Point{5, 5}, image.Point{7, 5}, true},
		{"2nd level diagonal", image.Point{5, 5}, image.Point{7, 7}, true},
		{"knight move", image.Point{5, 5}, image.Point{6, 7}, true},
		{"same point", image.Point{5, 5}, image.Point{5, 5}, false},
		{"three away", image.Point{5, 5}, image.Point{8, 5}, false},
		{"negative 2nd level", image.Point{-3, -3}, image.Point{-5, -4}, true},
		{"negative cross zero", image.Point{-1, -1}, image.Point{1, 0}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Are24Adjacent(tt.p, tt.q); got != tt.want {
				t.Errorf("Are24Adjacent(%v, %v) = %v, want %v", tt.p, tt.q, got, tt.want)
			}
		})
	}
}

func TestChebyshevDistance(t *testing.T) {
	tests := []struct {
		p, q image.Point
		want int
	}{
		{image.Point{0, 0}, image.Point{0, 0}, 0},
		{image.Point{0, 0}, image.Point{1, 0}, 1},
		{image.Point{0, 0}, image.Point{1, 1}, 1},
		{image.Point{0, 0}, image.Point{3, 2}, 3},
		{image.Point{5, 5}, image.Point{2, 3}, 3},
		{image.Point{-3, -3}, image.Point{-1, -2}, 2},
		{image.Point{-5, -5}, image.Point{-8, -7}, 3},
		{image.Point{-2, 3}, image.Point{1, -1}, 4},
		{image.Point{-1, -1}, image.Point{1, 1}, 2},
	}

	for _, tt := range tests {
		if got := ChebyshevDistance(tt.p, tt.q); got != tt.want {
			t.Errorf("ChebyshevDistance(%v, %v) = %v, want %v", tt.p, tt.q, got, tt.want)
		}
	}
}

func TestManhattanDistance(t *testing.T) {
	tests := []struct {
		p, q image.Point
		want int
	}{
		{image.Point{0, 0}, image.Point{0, 0}, 0},
		{image.Point{0, 0}, image.Point{1, 0}, 1},
		{image.Point{0, 0}, image.Point{1, 1}, 2},
		{image.Point{0, 0}, image.Point{3, 2}, 5},
		{image.Point{5, 5}, image.Point{2, 3}, 5},
		{image.Point{-3, -3}, image.Point{-1, -2}, 3},
		{image.Point{-5, -5}, image.Point{-8, -7}, 5},
		{image.Point{-2, 3}, image.Point{1, -1}, 7},
		{image.Point{-1, -1}, image.Point{1, 1}, 4},
	}

	for _, tt := range tests {
		if got := ManhattanDistance(tt.p, tt.q); got != tt.want {
			t.Errorf("ManhattanDistance(%v, %v) = %v, want %v", tt.p, tt.q, got, tt.want)
		}
	}
}

func TestGetNeighbors4(t *testing.T) {
	tests := []struct {
		name     string
		p        image.Point
		expected []image.Point
	}{
		{
			"positive coords",
			image.Point{5, 5},
			[]image.Point{{5, 4}, {4, 5}, {6, 5}, {5, 6}},
		},
		{
			"negative coords",
			image.Point{-3, -3},
			[]image.Point{{-3, -4}, {-4, -3}, {-2, -3}, {-3, -2}},
		},
		{
			"origin",
			image.Point{0, 0},
			[]image.Point{{0, -1}, {-1, 0}, {1, 0}, {0, 1}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			neighbors := GetNeighbors4(tt.p)

			if len(neighbors) != 4 {
				t.Fatalf("expected 4 neighbors, got %d", len(neighbors))
			}

			for _, exp := range tt.expected {
				found := false
				for _, n := range neighbors {
					if n == exp {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("expected neighbor %v not found in %v", exp, neighbors)
				}
			}
		})
	}
}

func TestGetNeighbors8(t *testing.T) {
	tests := []struct {
		name string
		p    image.Point
	}{
		{"positive coords", image.Point{5, 5}},
		{"negative coords", image.Point{-3, -3}},
		{"origin", image.Point{0, 0}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			neighbors := GetNeighbors8(tt.p)

			if len(neighbors) != 8 {
				t.Fatalf("expected 8 neighbors, got %d", len(neighbors))
			}

			// All neighbors should be at Chebyshev distance 1
			for _, n := range neighbors {
				if ChebyshevDistance(tt.p, n) != 1 {
					t.Errorf("neighbor %v is not at distance 1 from %v", n, tt.p)
				}
			}
		})
	}
}

func TestGetNeighbors24(t *testing.T) {
	tests := []struct {
		name string
		p    image.Point
	}{
		{"positive coords", image.Point{5, 5}},
		{"negative coords", image.Point{-3, -3}},
		{"origin", image.Point{0, 0}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			neighbors := GetNeighbors24(tt.p)

			if len(neighbors) != 24 {
				t.Fatalf("expected 24 neighbors, got %d", len(neighbors))
			}

			// All neighbors should be at Chebyshev distance <= 2
			for _, n := range neighbors {
				dist := ChebyshevDistance(tt.p, n)
				if dist > 2 || dist == 0 {
					t.Errorf("neighbor %v is at invalid distance %d from %v", n, dist, tt.p)
				}
			}
		})
	}
}

func TestGetKthLevelNeighbors(t *testing.T) {
	tests := []struct {
		name string
		p    image.Point
		k    int
		want int
	}{
		{"k=0", image.Point{5, 5}, 0, 0},
		{"k=1 positive", image.Point{5, 5}, 1, 8},
		{"k=2 positive", image.Point{5, 5}, 2, 16},
		{"k=3 positive", image.Point{5, 5}, 3, 24},
		{"k=1 negative", image.Point{-3, -3}, 1, 8},
		{"k=2 negative", image.Point{-3, -3}, 2, 16},
		{"k=1 origin", image.Point{0, 0}, 1, 8},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			neighbors := GetKthLevelNeighbors(tt.p, tt.k)
			if len(neighbors) != tt.want {
				t.Errorf("GetKthLevelNeighbors(%v, %d) returned %d neighbors, want %d",
					tt.p, tt.k, len(neighbors), tt.want)
			}

			// Verify all neighbors are at exactly distance k
			for _, n := range neighbors {
				if dist := ChebyshevDistance(tt.p, n); dist != tt.k {
					t.Errorf("neighbor %v is at distance %d from %v, want %d", n, dist, tt.p, tt.k)
				}
			}
		})
	}
}

func TestGetNeighborsWithinK(t *testing.T) {
	tests := []struct {
		name string
		p    image.Point
		k    int
		want int
	}{
		{"k=0", image.Point{5, 5}, 0, 0},
		{"k=1 positive", image.Point{5, 5}, 1, 8},
		{"k=2 positive", image.Point{5, 5}, 2, 24},
		{"k=1 negative", image.Point{-3, -3}, 1, 8},
		{"k=2 negative", image.Point{-3, -3}, 2, 24},
		{"k=1 origin", image.Point{0, 0}, 1, 8},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			neighbors := GetNeighborsWithinK(tt.p, tt.k)
			if len(neighbors) != tt.want {
				t.Errorf("GetNeighborsWithinK(%v, %d) returned %d neighbors, want %d",
					tt.p, tt.k, len(neighbors), tt.want)
			}

			// Verify all neighbors are within distance k
			for _, n := range neighbors {
				if dist := ChebyshevDistance(tt.p, n); dist > tt.k || dist == 0 {
					t.Errorf("neighbor %v is at distance %d from %v, should be 1-%d", n, dist, tt.p, tt.k)
				}
			}
		})
	}
}
