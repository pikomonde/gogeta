package gm

import "github.com/hajimehoshi/ebiten/v2"

// TODO: Consider different interface between Behaviour and Object
type Object interface {
	Init()
	Update()
	Draw(*ebiten.Image)
}

// Set an Instance to Game and initialize it.
func InitObject(obj Object) Object {
	gm.objects.setObject(obj)
	initBehaviours(obj)
	obj.Init()
	return obj
}
