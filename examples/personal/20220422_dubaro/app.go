package main

import (
	"fmt"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	bhvrCommon "github.com/pikomonde/gogeta/behaviour/behaviour_common"
	bhvrRoom "github.com/pikomonde/gogeta/behaviour/behaviour_room"
	"github.com/pikomonde/gogeta/gm"
)

const (
	WindowWidth  = 600
	WindowHeight = 600
	CanvasWidth  = 300
	CanvasHeight = 300
	// CanvasWidth  = 186
	// CanvasHeight = 300
	// CanvasWidth  = 93
	// CanvasHeight = 150
)

func main() {
	// Initialize objects
	gm.Init(WindowWidth, WindowHeight, CanvasWidth, CanvasHeight)
	gm.InitObject(&room01{})

	// Run game
	if err := gm.Run(); err != nil {
		log.Fatal("error game run: ", err)
	}
}

type room01 struct {
	gm.Objecter
	BhvrRoom bhvrRoom.Room
}

func (obj *room01) Init() {
	obj.BhvrRoom.InitObject(&obj01{}, bhvrRoom.InstanceData{})
	obj.BhvrRoom.InitObject(&obj02{}, bhvrRoom.InstanceData{})
}

func (obj *room01) Update() {
}

func (obj *room01) Draw(screen *ebiten.Image) {
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
