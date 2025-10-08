package topology

import (
	"image"
	"reflect"
	"testing"
)

func TestAreConnected(t *testing.T) {
	type args struct {
		p image.Point
		q image.Point
		S Region
	}
	tests := []struct {
		name          string
		args          args
		wantFoundPath []image.Point
		wantHasPath   bool
	}{
		{
			name: "same point in region",
			args: args{
				p: image.Point{X: 1, Y: 1},
				q: image.Point{X: 1, Y: 1},
				S: Region{
					{1, 1}: true,
					{2, 2}: true,
				},
			},
			wantFoundPath: []image.Point{{1, 1}},
			wantHasPath:   true,
		},
		{
			name: "same point not in region",
			args: args{
				p: image.Point{X: 1, Y: 1},
				q: image.Point{X: 1, Y: 1},
				S: Region{
					{2, 2}: true,
					{3, 3}: true,
				},
			},
			wantFoundPath: nil,
			wantHasPath:   false,
		},
		{
			name: "directly adjacent points",
			args: args{
				p: image.Point{X: 0, Y: 0},
				q: image.Point{X: 1, Y: 1},
				S: Region{
					{0, 0}: true,
					{1, 1}: true,
				},
			},
			wantFoundPath: []image.Point{{0, 0}, {1, 1}},
			wantHasPath:   true,
		},
		{
			name: "connected through path",
			args: args{
				p: image.Point{X: 0, Y: 0},
				q: image.Point{X: 2, Y: 0},
				S: Region{
					{0, 0}: true,
					{1, 0}: true,
					{2, 0}: true,
				},
			},
			wantFoundPath: []image.Point{{0, 0}, {1, 0}, {2, 0}},
			wantHasPath:   true,
		},
		{
			name: "not connected - gap in region",
			args: args{
				p: image.Point{X: 0, Y: 0},
				q: image.Point{X: 2, Y: 0},
				S: Region{
					{0, 0}: true,
					{2, 0}: true,
				},
			},
			wantFoundPath: nil,
			wantHasPath:   false,
		},
		{
			name: "start point not in region",
			args: args{
				p: image.Point{X: 0, Y: 0},
				q: image.Point{X: 1, Y: 1},
				S: Region{
					{1, 1}: true,
					{2, 2}: true,
				},
			},
			wantFoundPath: nil,
			wantHasPath:   false,
		},
		{
			name: "end point not in region",
			args: args{
				p: image.Point{X: 0, Y: 0},
				q: image.Point{X: 2, Y: 2},
				S: Region{
					{0, 0}: true,
					{1, 1}: true,
				},
			},
			wantFoundPath: nil,
			wantHasPath:   false,
		},
		{
			name: "empty region",
			args: args{
				p: image.Point{X: 0, Y: 0},
				q: image.Point{X: 1, Y: 1},
				S: Region{},
			},
			wantFoundPath: nil,
			wantHasPath:   false,
		},
		{
			name: "complex connected path",
			args: args{
				p: image.Point{X: 0, Y: 0},
				q: image.Point{X: 2, Y: 2},
				S: Region{
					{0, 0}: true,
					{1, 0}: true,
					{2, 0}: true,
					{2, 1}: true,
					{2, 2}: true,
				},
			},
			wantFoundPath: []image.Point{{0, 0}, {1, 0}, {2, 1}, {2, 2}},
			wantHasPath:   true,
		},
		{
			name: "surrounded but not connected",
			args: args{
				p: image.Point{X: 1, Y: 1},
				q: image.Point{X: 3, Y: 3},
				S: Region{
					{1, 1}: true,
					{3, 3}: true,
					{0, 0}: true,
					{0, 1}: true,
					{1, 0}: true,
					{4, 3}: true,
					{3, 4}: true,
				},
			},
			wantFoundPath: nil,
			wantHasPath:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFoundPath, gotHasPath := AreConnected(tt.args.p, tt.args.q, tt.args.S)
			if !reflect.DeepEqual(gotFoundPath, tt.wantFoundPath) {
				t.Errorf("AreConnected() gotFoundPath = %v, want %v", gotFoundPath, tt.wantFoundPath)
			}
			if gotHasPath != tt.wantHasPath {
				t.Errorf("AreConnected() gotHasPath = %v, want %v", gotHasPath, tt.wantHasPath)
			}
		})
	}
}

func TestIsDiscretePath(t *testing.T) {
	type args struct {
		path *[]image.Point
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "nil path",
			args: args{path: nil},
			want: false,
		},
		{
			name: "empty path",
			args: args{path: &[]image.Point{}},
			want: false,
		},
		{
			name: "single point path",
			args: args{path: &[]image.Point{{0, 0}}},
			want: true,
		},
		{
			name: "valid adjacent horizontal path",
			args: args{path: &[]image.Point{{0, 0}, {1, 0}, {2, 0}}},
			want: true,
		},
		{
			name: "valid adjacent vertical path",
			args: args{path: &[]image.Point{{0, 0}, {0, 1}, {0, 2}}},
			want: true,
		},
		{
			name: "valid diagonal path",
			args: args{path: &[]image.Point{{0, 0}, {1, 1}, {2, 2}}},
			want: true,
		},
		{
			name: "valid mixed adjacent path",
			args: args{path: &[]image.Point{{0, 0}, {1, 0}, {1, 1}, {2, 1}}},
			want: true,
		},
		{
			name: "non-adjacent points",
			args: args{path: &[]image.Point{{0, 0}, {2, 0}}},
			want: false,
		},
		{
			name: "non-adjacent in middle",
			args: args{path: &[]image.Point{{0, 0}, {1, 0}, {3, 0}}},
			want: false,
		},
		{
			name: "repeated point",
			args: args{path: &[]image.Point{{0, 0}, {1, 0}, {0, 0}}},
			want: false,
		},
		{
			name: "repeated consecutive points",
			args: args{path: &[]image.Point{{0, 0}, {0, 0}}},
			want: false,
		},
		{
			name: "repeated point in middle",
			args: args{path: &[]image.Point{{0, 0}, {1, 0}, {2, 0}, {1, 0}}},
			want: false,
		},
		{
			name: "valid L-shaped path",
			args: args{path: &[]image.Point{{0, 0}, {1, 0}, {2, 0}, {2, 1}}},
			want: true,
		},
		{
			name: "valid complex path",
			args: args{path: &[]image.Point{{0, 0}, {0, 1}, {1, 1}, {1, 2}, {2, 2}}},
			want: true,
		},
		{
			name: "Manhattan distance 2 adjacency",
			args: args{path: &[]image.Point{{0, 0}, {2, 0}}},
			want: false,
		},
		{
			name: "Manhattan distance 3 - not adjacent",
			args: args{path: &[]image.Point{{0, 0}, {3, 0}}},
			want: false,
		},
		{
			name: "negative coordinates valid path",
			args: args{path: &[]image.Point{{-1, -1}, {0, -1}, {0, 0}}},
			want: true,
		},
		{
			name: "negative coordinates with repetition",
			args: args{path: &[]image.Point{{-1, -1}, {0, -1}, {-1, -1}}},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsDiscretePath(tt.args.path); got != tt.want {
				t.Errorf("IsDiscretePath() = %v, want %v", got, tt.want)
			}
		})
	}
}
