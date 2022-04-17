package behaviour_common

import (
	"github.com/golang/geo/r2"
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
