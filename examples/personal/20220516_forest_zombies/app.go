package main

import (
	"fmt"
	_ "image/png"
	"log"

	"github.com/golang/geo/r2"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	bhvrCommon "github.com/pikomonde/gogeta/behaviour/behaviour_common"
	bhvrRoom "github.com/pikomonde/gogeta/behaviour/behaviour_room"
	"github.com/pikomonde/gogeta/gm"
	"github.com/pikomonde/gogeta/gogetautil/txt"
)

const (
	CanvasWidth  = 372
	CanvasHeight = 600
)

var fontBocil57 = txt.MustNewFontFromFile("asset/sprite/font_bocil_57_0020_007F.png", 5, 7, txt.CharSet_0020_007F, txt.Font{Size: 24})

func main() {
	// Initialize objects
	gm.SetLayoutType(gm.LayoutType_SnapOutside)
	gm.InitObject(&roomMain{})
	ebiten.SetWindowSize(CanvasWidth, CanvasHeight)

	// Run game
	if err := gm.Run(); err != nil {
		log.Fatal("error game run: ", err)
	}
}

type roomMain struct {
	gm.Objecter
	BhvrRoom bhvrRoom.Room
}

func (obj *roomMain) Init() {
	obj.BhvrRoom.Size = r2.Point{X: float64(CanvasWidth), Y: float64(CanvasHeight)}
	obj.BhvrRoom.InitObject(&obj01{}, bhvrRoom.InstanceData{})
	obj.BhvrRoom.InitObject(&obj02{}, bhvrRoom.InstanceData{})
}

func (obj *roomMain) Update() {
}

func (obj *roomMain) Draw(screen *ebiten.Image) {
	x, y := ebiten.CursorPosition()
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("mouse loc: %d %d", x, y), 8, 8)
}

type obj01 struct {
	gm.Objecter
	BhvrCommon bhvrCommon.Common
}

func (obj *obj01) Init() {
}

func (obj *obj01) Update() {
}

func (obj *obj01) Draw(screen *ebiten.Image) {
}

type obj02 struct {
	gm.Objecter
	BhvrCommon bhvrCommon.Common
}

func (obj *obj02) Init() {
}

func (obj *obj02) Update() {
}

func (obj *obj02) Draw(screen *ebiten.Image) {
}
