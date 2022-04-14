package gm

import (
	"github.com/hajimehoshi/ebiten/v2"
)

var gm game

type game struct {
	objects objects // all objects in the game, indexed by "object type" -> "object interface pointer"
	// parentObjects behaviourObjects // all objects in the game, indexed by "behvaiour interface pointer"
	layoutW, layoutH int
}

func GetObjectDB() objects {
	return gm.objects
}

func (g *game) Update() error {
	g.objects.update()
	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	// Draw object
	g.objects.draw(screen)
}

func (g *game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.layoutW, g.layoutH
}

func Init(windowW, windowH, layoutW, layoutH int) error {
	gm.objects = make(objects)
	ebiten.SetWindowSize(windowW, windowH)
	gm.layoutW, gm.layoutH = layoutW, layoutH
	return nil
}

func Run() error {
	return ebiten.RunGame(&gm)
}
