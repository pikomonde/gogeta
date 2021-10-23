package sys

import (
	"math"
	"math/rand"

	"github.com/golang/geo/r2"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	FullPi = 2 * math.Pi
)

func AngleFromPoints(x1, y1, x2, y2 float64) float64 {
	return math.Atan2(y2-y1, x2-x1)
}

func AngleFromVector(vx, vy float64) float64 {
	return math.Atan2(vy, vx)
}

func RandUnitCircle() r2.Point {
	angle := (rand.Float64() - 0.5) * FullPi

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(1, 0)
	op.GeoM.Rotate(angle)
	xx, yy := op.GeoM.Apply(0, 0)
	return r2.Point{xx, yy}
}
