package behaviour_common

import (
	"github.com/golang/geo/r1"
	"github.com/golang/geo/r2"
	"github.com/hajimehoshi/ebiten/v2"
)

type Mask struct {
	maskType MaskType
	vectors  []r2.Point
}

func (bhvr *Mask) MaskType() MaskType { return bhvr.maskType }

func (bhvr *Mask) Vectors() []r2.Point { return bhvr.vectors }

// === Variable and Constant ===

type MaskType int

const (
	Sprite_MaskType_NoMask MaskType = iota
	Sprite_MaskType_Circle
	Sprite_MaskType_Recatangle
	Sprite_MaskType_Capsule
	Sprite_MaskType_ConvexHull
)

// === Package functions ===

// OuterRectangle returns rectangle that contain this mask
func (bhvr Mask) OuterRectangle() r2.Rect {
	switch bhvr.maskType {
	case Sprite_MaskType_Circle:
		// TODO:
		return r2.EmptyRect()
	case Sprite_MaskType_Recatangle:
		if bhvr.vectors[0].X >= bhvr.vectors[1].X {
			if bhvr.vectors[1].X >= bhvr.vectors[2].X {
				// Quadrant I
				return r2.Rect{
					X: r1.Interval{Lo: bhvr.vectors[0].X, Hi: bhvr.vectors[2].X},
					Y: r1.Interval{Lo: bhvr.vectors[1].Y, Hi: bhvr.vectors[3].Y},
				}
			}
			// Quadrant II
			return r2.Rect{
				X: r1.Interval{Lo: bhvr.vectors[3].X, Hi: bhvr.vectors[1].X},
				Y: r1.Interval{Lo: bhvr.vectors[0].Y, Hi: bhvr.vectors[2].Y},
			}
		}
		if bhvr.vectors[1].X <= bhvr.vectors[2].X {
			// Quadrant III
			return r2.Rect{
				X: r1.Interval{Lo: bhvr.vectors[2].X, Hi: bhvr.vectors[0].X},
				Y: r1.Interval{Lo: bhvr.vectors[3].Y, Hi: bhvr.vectors[1].Y},
			}
		}
		// Quadrant IV
		return r2.Rect{
			X: r1.Interval{Lo: bhvr.vectors[1].X, Hi: bhvr.vectors[3].X},
			Y: r1.Interval{Lo: bhvr.vectors[2].Y, Hi: bhvr.vectors[0].Y},
		}

	case Sprite_MaskType_Capsule:
		// TODO:
		return r2.EmptyRect()
	case Sprite_MaskType_ConvexHull:
		// TODO:
		return r2.EmptyRect()
	}
	return r2.EmptyRect()
}

// GeoTransform returns new Mask with already transformed by geoM
func (bhvr Mask) GeoTransform(geoM ebiten.GeoM) Mask {
	newVectors := make([]r2.Point, 0)
	for _, vector := range bhvr.vectors {
		newX, newY := geoM.Apply(vector.X, vector.Y)
		newVectors = append(newVectors, r2.Point{X: newX, Y: newY})
	}
	return Mask{
		maskType: bhvr.maskType,
		vectors:  newVectors,
	}
}
