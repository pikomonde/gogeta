package gm

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/pikomonde/game-210530-theMacaronChef/constant"
)

var gm game

type game struct {
	objects objects // all objects in the game, indexed by "object type" -> "object interface pointer"
	// parentObjects behaviourObjects // all objects in the game, indexed by "behvaiour interface pointer"
}

func (g *game) Update() error {
	// btn.Update()
	g.objects.Update()
	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	// Color the background
	screen.Fill(color.NRGBA{0x00, 0x40, 0x80, 0xff})

	// Draw object
	g.objects.Draw(Screen(screen))
}

func (g *game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return constant.CanvasWidth, constant.CanvasHeight
}

func Init() error {
	gm.objects = make(objects)
	ebiten.SetWindowSize(constant.WindowWidth, constant.WindowHeight)
	return nil
}

func Run() error {
	return ebiten.RunGame(&gm)
}
