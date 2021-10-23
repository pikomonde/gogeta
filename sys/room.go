package sys

import (
	"github.com/golang/geo/r2"
)

func Bound(pos r2.Point, x1, y1, x2, y2 float64) r2.Point {
	if pos.X < x1 {
		pos.X = x1
	}
	if pos.Y < y1 {
		pos.Y = y1
	}
	if pos.X > x2 {
		pos.X = x2
	}
	if pos.Y > y2 {
		pos.Y = y2
	}
	return pos
}
