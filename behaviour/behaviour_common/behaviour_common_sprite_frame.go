package behaviour_common

import (
	"github.com/golang/geo/r2"
	"github.com/hajimehoshi/ebiten/v2"
)

type Frame struct {
	Image  *ebiten.Image
	Anchor r2.Point
}

func (bhvr *Frame) SetAnchorToggle(pos FrameAnchorToggle) *Frame {
	w, h := bhvr.Image.Size()
	switch pos {
	case Sprite_FrameAnchor_ToggleTopLeft:
		bhvr.Anchor = r2.Point{0, 0}
	case Sprite_FrameAnchor_ToggleTopCenter:
		bhvr.Anchor = r2.Point{float64(w) / 2, 0}
	case Sprite_FrameAnchor_ToggleTopRight:
		bhvr.Anchor = r2.Point{float64(w), 0}
	case Sprite_FrameAnchor_ToggleMiddleLeft:
		bhvr.Anchor = r2.Point{0, float64(h) / 2}
	case Sprite_FrameAnchor_ToggleMiddleCenter:
		bhvr.Anchor = r2.Point{float64(w) / 2, float64(h) / 2}
	case Sprite_FrameAnchor_ToggleMiddleRight:
		bhvr.Anchor = r2.Point{float64(w), float64(h) / 2}
	case Sprite_FrameAnchor_ToggleBottomLeft:
		bhvr.Anchor = r2.Point{0, float64(h)}
	case Sprite_FrameAnchor_ToggleBottomCenter:
		bhvr.Anchor = r2.Point{float64(w) / 2, float64(h)}
	case Sprite_FrameAnchor_ToggleBottomRight:
		bhvr.Anchor = r2.Point{float64(w), float64(h)}
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
