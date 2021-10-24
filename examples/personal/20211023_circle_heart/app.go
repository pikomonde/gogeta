package main

import (
	"image/color"
	_ "image/png"
	"log"

	"github.com/golang/geo/r2"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/pikomonde/gogeta/behaviour"
	bhvrCommon "github.com/pikomonde/gogeta/behaviour/behaviour_common"
	"github.com/pikomonde/gogeta/gm"
)

const (
	WindowWidth  = 372
	WindowHeight = 600
	CanvasWidth  = 186
	CanvasHeight = 300
)

func main() {
	// Initialize objects
	gm.Init(WindowWidth, WindowHeight)
	gm.InitObject(&room01{})

	// Run game
	if err := gm.Run(); err != nil {
		log.Fatal("error game run: ", err)
	}
}

type room01 struct{}

func (obj *room01) Init() {
	heart := objHeart{BhvrCommon: behaviour.Common{Position: r2.Point{100, 100}}}
	gm.InitObject(&heart)
	// gm.InitObject(&Snail{
	// 	// Common: behaviour.Common{
	// 	// 	Position: r2.Point{100, 100},
	// 	// },
	// })
	// gm.InitObject(&Tile{})
}

func (obj *room01) Update() {
}

func (obj *room01) Draw(screen *ebiten.Image) {
}

/*
  objHeart
*/

type objHeart struct {
	BhvrCommon behaviour.Common
}

func (obj *objHeart) Init() {
	frame := &bhvrCommon.Frame{Image: ebiten.NewImage(int(32), int(32))}
	frame.Image.Fill(color.NRGBA{0x00, 0x80, 0x00, 0xff})
	frame.SetAnchorToggle(bhvrCommon.Sprite_FrameAnchor_ToggleDefault)
	obj.BhvrCommon.Sprite.InsertFrame("", frame)

}

func (obj *objHeart) Update() {
}

func (obj *objHeart) Draw(screen *ebiten.Image) {
	// Color the background
	// screen.Fill(color.NRGBA{0x00, 0x40, 0x80, 0xff})
}
