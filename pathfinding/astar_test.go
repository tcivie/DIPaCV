package pathfinding

import (
	"image"
	"reflect"
	"testing"
)

func TestAStar(t *testing.T) {
	type args struct {
		start        image.Point
		goal         image.Point
		isValidPoint ValidPointChecker
	}
	makeValidPointChecker := func(obstacles []image.Point) ValidPointChecker {
		return func(p image.Point) bool {
			for _, o := range obstacles {
				if p == o {
					return false
				}
			}
			return p.X >= 0 && p.Y >= 0 && p.X < 10 && p.Y < 10
		}
	}

	tests := []struct {
		name  string
		args  args
		want  []image.Point
		want1 bool
	}{
		{
			name: "direct path",
			args: args{
				start:        image.Point{},
				goal:         image.Point{X: 2, Y: 2},
				isValidPoint: makeValidPointChecker(nil),
			},
			want: []image.Point{
				{0, 0},
				{1, 1},
				{2, 2},
			},
			want1: true,
		},
		{
			name: "path with obstacle",
			args: args{
				start:        image.Point{},
				goal:         image.Point{X: 2, Y: 2},
				isValidPoint: makeValidPointChecker([]image.Point{{1, 1}}),
			},
			want: []image.Point{
				{0, 0},
				{1, 0},
				{2, 1},
				{2, 2},
			},
			want1: true,
		},
		{
			name: "no path possible",
			args: args{
				start: image.Point{},
				goal:  image.Point{X: 2, Y: 2},
				isValidPoint: makeValidPointChecker(
					[]image.Point{
						{1, 1},
						{1, 0},
						{0, 1},
						{2, 1},
						{1, 2}}),
			},
			want:  nil,
			want1: false,
		},
		{
			name: "start point is goal",
			args: args{
				start:        image.Point{X: 1, Y: 1},
				goal:         image.Point{X: 1, Y: 1},
				isValidPoint: makeValidPointChecker(nil),
			},
			want:  []image.Point{{1, 1}},
			want1: true,
		},
		{
			name: "invalid start point",
			args: args{
				start:        image.Point{X: -1},
				goal:         image.Point{X: 2, Y: 2},
				isValidPoint: makeValidPointChecker(nil),
			},
			want:  nil,
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := AStar(tt.args.start, tt.args.goal, tt.args.isValidPoint)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AStar() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("AStar() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
