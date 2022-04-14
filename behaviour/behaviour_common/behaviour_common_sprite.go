package behaviour_common

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Sprite struct {
	Animations       Animations
	CurrentAnimation string
	CurrentFrame     int
}

func (bhvr *Sprite) PreInit() {
	// create Amintaions if not exist
	if bhvr.Animations == nil {
		bhvr.Animations = make(Animations)
		return
	}

	// use available animation as default if exist and not defined yet
	if bhvr.CurrentAnimation == "" {
		if len(bhvr.Animations) >= 1 {
			for k, _ := range bhvr.Animations {
				bhvr.CurrentAnimation = k
			}
		}
	}
}

func (bhvr *Sprite) PostInit() {
	// create Amintaions if not exist
	if bhvr.Animations == nil {
		bhvr.Animations = make(Animations)
		return
	}

	// use available animation as default if exist and not defined yet
	if bhvr.CurrentAnimation == "" {
		if len(bhvr.Animations) >= 1 {
			for k, _ := range bhvr.Animations {
				bhvr.CurrentAnimation = k
			}
		}
	}
}

func (bhvr *Sprite) GetCurrentFrame() *Frame {
	if _, ok := bhvr.Animations[bhvr.CurrentAnimation]; !ok {
		return defaultFrame
	}
	if bhvr.CurrentFrame >= len(*bhvr.Animations[bhvr.CurrentAnimation]) {
		return defaultFrame
	}
	return (*bhvr.Animations[bhvr.CurrentAnimation])[bhvr.CurrentFrame]
}

// === Behaviour specific method ===

// InsertFrame insert frame(s) at the end of the Animation named
// "animationName" by Image.
func (bhvr *Sprite) InsertFrame(animationName string, newFrames ...*Frame) {
	animation := bhvr.createAnimationIfNotExist(animationName)
	*animation = append(*animation, newFrames...)
}

// InsertFrameByImage insert frame(s) at the end of the Animation named
// "animationName" by Image.
func (bhvr *Sprite) InsertFrameByImage(AnimationName string, images ...*ebiten.Image) {
	newFrames := make([]*Frame, 0)

	for _, image := range images {
		newFrame := &Frame{Image: image}
		newFrames = append(newFrames, newFrame)
	}

	bhvr.InsertFrame(AnimationName, newFrames...)
}

// createAnimationIfNotExist creates new animation for the object.
func (bhvr *Sprite) createAnimationIfNotExist(animationName string) *Animation {
	if _, ok := bhvr.Animations[animationName]; !ok {
		animation := make(Animation, 0)
		bhvr.Animations[animationName] = &animation
	}
	return bhvr.Animations[animationName]
}

// func (bhvr *Sprite) SetSize(x, y int) {
// 	bhvr.Image = ebiten.NewImage(x, y)
// }

// func (bhvr *Sprite) FillColor(colorNRGBA color.NRGBA) {
// 	bhvr.Image.Fill(colorNRGBA)
// }

// func (bhvr *Sprite) SetSizeAndFillColor(x, y int, colorNRGBA color.NRGBA) {
// 	bhvr.SetSize(x, y)
// 	bhvr.FillColor(colorNRGBA)
// }

// func (bhvr *Sprite) SetSprite(image *ebiten.Image) {
// 	bhvr.Image = image
// }

// === Variable and Constant ===

var (
	defaultFrame = (&Frame{
		Image: ebiten.NewImage(int(32), int(32)),
	})
)

const (
	DefaultAnimationName = "default"
)
