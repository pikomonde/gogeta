package gm

// TODO: Consider different interface between Behaviour and Object
type Object interface {
	Init()
	Update()
	Draw(Screen)
}

// Set an Instance to Game and initialize it.
func SetAndInitObject(obj Object) Object {
	gm.objects.setObject(obj)
	initBehaviours(obj)
	obj.Init()
	return obj
}
