package geometry

import (
	"image"
	"testing"
)

func TestEuclideanLength(t *testing.T) {
	type args struct {
		p image.Point
		q image.Point
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "zero distance",
			args: args{p: image.Point{0, 0}, q: image.Point{0, 0}},
			want: 0,
		},
		{
			name: "horizontal line",
			args: args{p: image.Point{0, 0}, q: image.Point{3, 0}},
			want: 3,
		},
		{
			name: "vertical line",
			args: args{p: image.Point{0, 0}, q: image.Point{0, 4}},
			want: 4,
		},
		{
			name: "diagonal line",
			args: args{p: image.Point{0, 0}, q: image.Point{3, 4}},
			want: 5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EuclideanLength(tt.args.p, tt.args.q); got != tt.want {
				t.Errorf("EuclideanLength() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManhattanLength(t *testing.T) {
	type args struct {
		p image.Point
		q image.Point
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "zero distance",
			args: args{p: image.Point{0, 0}, q: image.Point{0, 0}},
			want: 0,
		},
		{
			name: "horizontal line",
			args: args{p: image.Point{0, 0}, q: image.Point{3, 0}},
			want: 3,
		},
		{
			name: "vertical line",
			args: args{p: image.Point{0, 0}, q: image.Point{0, 4}},
			want: 4,
		},
		{
			name: "diagonal line",
			args: args{p: image.Point{0, 0}, q: image.Point{3, 4}},
			want: 7,
		},
		{
			name: "negative coordinates",
			args: args{p: image.Point{-2, -3}, q: image.Point{2, 1}},
			want: 8,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ManhattanLength(tt.args.p, tt.args.q); got != tt.want {
				t.Errorf("ManhattanLength() = %v, want %v", got, tt.want)
			}
		})
	}
}
