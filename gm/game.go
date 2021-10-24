package gm

import (
	"github.com/hajimehoshi/ebiten/v2"
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
	// Draw object
	g.objects.Draw(screen)
}

func (g *game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func Init(w, h int) error {
	gm.objects = make(objects)
	ebiten.SetWindowSize(w, h)
	return nil
}

func Run() error {
	return ebiten.RunGame(&gm)
}
