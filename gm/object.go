package gm

import (
	"log"
	"reflect"

	"github.com/hajimehoshi/ebiten/v2"
)

// TODO: Consider different interface between Behaviour and Object
type Object interface {
	IDerInterface
	Zidx() int
	SetZidx(int)
	setZidx(int)
	IsUpdate() bool
	SetIsUpdate(bool)
	IsDraw() bool
	SetIsDraw(bool)

	Init()
	Update()
	Draw(*ebiten.Image)
}

type ObjectData struct {
	IDer
	zidx        int
	isNotUpdate bool
	isNotDraw   bool
}

func (inst *ObjectData) Zidx() int          { return inst.zidx }
func (inst *ObjectData) SetZidx(d int)      { UpdateInstancesZidx(GetInstancesByObjInst()[inst.id], d) }
func (inst *ObjectData) setZidx(d int)      { inst.zidx = d }
func (inst *ObjectData) IsUpdate() bool     { return !inst.isNotUpdate }
func (inst *ObjectData) SetIsUpdate(d bool) { inst.isNotUpdate = !d }
func (inst *ObjectData) IsDraw() bool       { return !inst.isNotDraw }
func (inst *ObjectData) SetIsDraw(d bool)   { inst.isNotDraw = !d }

// Set an Instance to Game and initialize it.
func InitObject(inst Object) Object {
	gm.setInstance(inst)
	preInitBehaviours(inst)
	inst.Init()
	postInitBehaviours(inst)
	return inst
}

// Delete an Instance from Game and destroy it.
func DelObject(obj Object) {
	delBehaviours(obj)
	delInstanceDataOnBehaviour(obj)
	gm.delInstance(obj)
}

func preInitBehaviours(obj Object) {
	objReflectVal := reflect.Indirect(reflect.ValueOf(obj))

	for i := 0; i < objReflectVal.NumField(); i++ {
		if !objReflectVal.Field(i).Addr().Type().Implements(reflect.TypeOf((*Behaviour)(nil)).Elem()) {
			continue
		}
		if !objReflectVal.Field(i).Addr().CanInterface() {
			log.Panicf("Behaviour %v in Object %v should be exported\n", objReflectVal.Field(i).Type().Name(), objReflectVal.Type().Name())
		}
		if bhvr, ok := objReflectVal.Field(i).Addr().Interface().(Behaviour); ok {
			gm.setBehaviour(obj, bhvr)
			bhvr.PreInit()
		}
	}
}

func postInitBehaviours(obj Object) {
	objReflectVal := reflect.Indirect(reflect.ValueOf(obj))

	for i := 0; i < objReflectVal.NumField(); i++ {
		if !objReflectVal.Field(i).Addr().Type().Implements(reflect.TypeOf((*Behaviour)(nil)).Elem()) {
			continue
		}
		if !objReflectVal.Field(i).Addr().CanInterface() {
			log.Panicf("Behaviour %v in Object %v should be exported\n", objReflectVal.Field(i).Type().Name(), objReflectVal.Type().Name())
		}
		if bhvr, ok := objReflectVal.Field(i).Addr().Interface().(Behaviour); ok {
			// gm.setBehaviour(obj, bhvr)
			bhvr.PostInit()
		}
	}
}

// func preUpdateBehaviours(obj Object) {
// 	objReflectVal := reflect.Indirect(reflect.ValueOf(obj))

// 	for i := 0; i < objReflectVal.NumField(); i++ {
// 		if !objReflectVal.Field(i).Addr().Type().Implements(reflect.TypeOf((*Behaviour)(nil)).Elem()) {
// 			continue
// 		}
// 		if !objReflectVal.Field(i).Addr().CanInterface() {
// 			log.Panicf("Behaviour %v in Object %v should be exported\n", objReflectVal.Field(i).Type().Name(), objReflectVal.Type().Name())
// 		}
// 		if bhvr, ok := objReflectVal.Field(i).Addr().Interface().(Behaviour); ok {
// 			bhvr.PreUpdate()
// 		}
// 	}
// }

// func postUpdateBehaviours(obj Object) {
// 	objReflectVal := reflect.Indirect(reflect.ValueOf(obj))

// 	for i := 0; i < objReflectVal.NumField(); i++ {
// 		if !objReflectVal.Field(i).Addr().Type().Implements(reflect.TypeOf((*Behaviour)(nil)).Elem()) {
// 			continue
// 		}
// 		if !objReflectVal.Field(i).Addr().CanInterface() {
// 			log.Panicf("Behaviour %v in Object %v should be exported\n", objReflectVal.Field(i).Type().Name(), objReflectVal.Type().Name())
// 		}
// 		if bhvr, ok := objReflectVal.Field(i).Addr().Interface().(Behaviour); ok {
// 			bhvr.PostUpdate()
// 		}
// 	}
// }

func delBehaviours(obj Object) {
	objReflectVal := reflect.Indirect(reflect.ValueOf(obj))

	for i := 0; i < objReflectVal.NumField(); i++ {
		if !objReflectVal.Field(i).Addr().Type().Implements(reflect.TypeOf((*Behaviour)(nil)).Elem()) {
			continue
		}
		if !objReflectVal.Field(i).Addr().CanInterface() {
			log.Panicf("Behaviour %v in Object %v should be exported\n", objReflectVal.Field(i).Type().Name(), objReflectVal.Type().Name())
		}
		if bhvr, ok := objReflectVal.Field(i).Addr().Interface().(Behaviour); ok {
			gm.delBehaviour(obj, bhvr)
		}
	}
}

func delInstanceDataOnBehaviour(obj Object) {
	for _, bhvrsData := range gm.behavioursData.byBhvrType {
		bhvrsData.DelInstance(obj)
	}
}

// func drawBehaviours(obj Object, screen *ebiten.Image) {
// 	objReflectVal := reflect.Indirect(reflect.ValueOf(obj))

// 	for i := 0; i < objReflectVal.NumField(); i++ {
// 		if !objReflectVal.Field(i).Addr().Type().Implements(reflect.TypeOf((*Behaviour)(nil)).Elem()) {
// 			continue
// 		}
// 		if !objReflectVal.Field(i).Addr().CanInterface() {
// 			log.Panicf("Behaviour %v in Object %v should be exported\n", objReflectVal.Field(i).Type().Name(), objReflectVal.Type().Name())
// 		}
// 		if bhvr, ok := objReflectVal.Field(i).Addr().Interface().(Behaviour); ok {
// 			bhvr.Draw(screen)
// 		}
// 	}
// }
