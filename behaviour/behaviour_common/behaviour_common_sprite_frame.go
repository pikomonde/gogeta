package behaviour_common

import (
	"log"

	"github.com/golang/geo/r2"
	"github.com/hajimehoshi/ebiten/v2"
)

type Frame struct {
	image  *ebiten.Image
	mask   Mask
	anchor r2.Point
}

func (bhvr *Frame) Image() *ebiten.Image { return bhvr.image }

func (bhvr *Frame) Mask() *Mask { return &bhvr.mask }

func (bhvr *Frame) MaskType() MaskType { return bhvr.mask.maskType }

func (bhvr *Frame) Anchor() r2.Point { return bhvr.anchor }

func (bhvr *Frame) SetImage(image *ebiten.Image) *Frame {
	bhvr.image = image
	bhvr.mask.maskType = Sprite_MaskType_NoMask
	bhvr.SetAnchorToggle(Sprite_FrameAnchor_ToggleTopLeft)
	return bhvr
}

func (bhvr *Frame) RemoveMask() *Frame {
	bhvr.mask.maskType = Sprite_MaskType_NoMask
	bhvr.mask.vectors = nil
	return bhvr
}

func (bhvr *Frame) SetMaskRectangle(mask r2.Rect) *Frame {
	bhvr.mask.maskType = Sprite_MaskType_Recatangle
	w, h := bhvr.image.Size()
	if (mask.X.Lo < 0) || (mask.Y.Lo < 0) || (mask.X.Hi > float64(w)) || (mask.Y.Hi > float64(h)) {
		log.Panicf("Mask is out of image boundary. Image {width: %d, height: %d}, mask {%s}", w, h, mask.String())
	}
	bhvr.mask.vectors = []r2.Point{
		{X: mask.X.Lo, Y: mask.Y.Lo},
		{X: mask.X.Hi, Y: mask.Y.Lo},
		{X: mask.X.Hi, Y: mask.Y.Hi},
		{X: mask.X.Lo, Y: mask.Y.Hi},
	}
	return bhvr
}

func (bhvr *Frame) SetMaskFill() *Frame {
	bhvr.mask.maskType = Sprite_MaskType_Recatangle
	w, h := bhvr.image.Size()
	bhvr.mask.vectors = []r2.Point{
		{X: 0, Y: 0},
		{X: float64(w), Y: 0},
		{X: float64(w), Y: float64(h)},
		{X: 0, Y: float64(h)},
	}
	return bhvr
}

func (bhvr *Frame) SetAnchor(p r2.Point) *Frame {
	w, h := bhvr.image.Size()
	if (p.X < 0) || (p.Y < 0) || (p.X > float64(w)) || (p.Y > float64(h)) {
		log.Panicf("Anchor is out of image boundary. Image {width: %d, height: %d}, anchor {X: %2f, Y: %2f}", w, h, p.X, p.Y)
	}
	return bhvr
}

func (bhvr *Frame) SetAnchorToggle(pos FrameAnchorToggle) *Frame {
	w, h := bhvr.image.Size()
	switch pos {
	case Sprite_FrameAnchor_ToggleTopLeft:
		bhvr.anchor = r2.Point{0, 0}
	case Sprite_FrameAnchor_ToggleTopCenter:
		bhvr.anchor = r2.Point{float64(w) / 2, 0}
	case Sprite_FrameAnchor_ToggleTopRight:
		bhvr.anchor = r2.Point{float64(w), 0}
	case Sprite_FrameAnchor_ToggleMiddleLeft:
		bhvr.anchor = r2.Point{0, float64(h) / 2}
	case Sprite_FrameAnchor_ToggleMiddleCenter:
		bhvr.anchor = r2.Point{float64(w) / 2, float64(h) / 2}
	case Sprite_FrameAnchor_ToggleMiddleRight:
		bhvr.anchor = r2.Point{float64(w), float64(h) / 2}
	case Sprite_FrameAnchor_ToggleBottomLeft:
		bhvr.anchor = r2.Point{0, float64(h)}
	case Sprite_FrameAnchor_ToggleBottomCenter:
		bhvr.anchor = r2.Point{float64(w) / 2, float64(h)}
	case Sprite_FrameAnchor_ToggleBottomRight:
		bhvr.anchor = r2.Point{float64(w), float64(h)}
	}
	return bhvr
}

// === Variable and Constant ===

type FrameAnchorToggle int

const (
	Sprite_FrameAnchor_ToggleTopLeft FrameAnchorToggle = iota
	Sprite_FrameAnchor_ToggleTopCenter
	Sprite_FrameAnchor_ToggleTopRight
	Sprite_FrameAnchor_ToggleMiddleLeft
	Sprite_FrameAnchor_ToggleMiddleCenter
	Sprite_FrameAnchor_ToggleMiddleRight
	Sprite_FrameAnchor_ToggleBottomLeft
	Sprite_FrameAnchor_ToggleBottomCenter
	Sprite_FrameAnchor_ToggleBottomRight
)
