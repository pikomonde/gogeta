package gm

import (
	"errors"
	"fmt"
	"log"
	"reflect"

	"github.com/hajimehoshi/ebiten/v2"
)

// TODO: Consider different interface between Behaviour and Object
type Behaviour interface {
	Init()
	Update()
	Draw(*ebiten.Image)
}

// type behaviourData struct {
// }

// type behaviours map[string]Behaviour

// // Set a Behaviour to Object.
// func (bhvrs behaviours) set(bhvr Behaviour) {
// 	bhvrType := reflect.TypeOf(bhvr).String()
// 	bhvrs[bhvrType] = bhvr
// }

// // Get a Behaviour from Object.
// func (bhvrs behaviours) get(bhvr Behaviour) Behaviour {
// 	bhvrType := reflect.TypeOf(bhvr).String()
// 	return bhvrs[bhvrType]
// }

// // Update all Instances.
// func (bhvrs behaviours) Update() {
// 	for _, bhvr := range bhvrs {
// 		bhvr.Update()
// 	}
// }

// // Draw all Instances.
// func (bhvrs behaviours) Draw(screen Screen) {
// 	for _, bhvr := range bhvrs {
// 		bhvr.Draw(screen)
// 	}
// }

func (objs objects) setBehaviour(obj Object, bhvr Behaviour) {
	// no need to set behaviour to object's behaviour, because behaviour
	// already set programatically
	objd := objectData{object: obj}
	bhvrType := reflect.TypeOf(bhvr).String()

	// case: keyByBhvr (set behaviour to parentObject-behvaiour relation on
	// top level game)
	key := fmt.Sprintf("%s%p", keyByBhvr, bhvr)
	if _, ok := gm.objects[key]; !ok {
		gm.objects[key] = make(map[Object]objectData)
	}
	gm.objects[key][obj] = objd

	// case: keyByBhvrType
	key = fmt.Sprintf("%s%s", keyByBhvrType, bhvrType)
	if _, ok := gm.objects[key]; !ok {
		gm.objects[key] = make(map[Object]objectData)
	}
	gm.objects[key][obj] = objd
}

func (objs objects) delBehaviour(obj Object, bhvr Behaviour) {
	bhvrType := reflect.TypeOf(bhvr).String()

	// case: keyByBhvr (set behaviour to parentObject-behvaiour relation on
	// top level game)
	key := fmt.Sprintf("%s%p", keyByBhvr, bhvr)
	delete(objs[key], obj)
	delete(objs, key)

	// case: keyByBhvrType
	key = fmt.Sprintf("%s%s", keyByBhvrType, bhvrType)
	delete(objs[key], obj)
	if len(objs[key]) == 0 {
		delete(objs, key)
	}
}

func (objs objects) getParentObjectByBehaviour(bhvr Behaviour) (Object, error) {
	key := fmt.Sprintf("%s%p", keyByBhvr, bhvr)
	mapObj, ok := gm.objects[key]
	if !ok {
		return nil, errors.New(ErrParentObjectNotFound)
	}
	for _, objData := range mapObj {
		return objData.object, nil
	}
	// there should be no case of this, because if there is no member of
	// mapObj, gm.objects[key] should not be exist
	return nil, errors.New(ErrParentObjectNotFound)
}

// // Get Object's Behaviour by type.
// func GetBehaviour(obj Object, bhvrType Behaviour) (Behaviour, error) {
// 	if bhvr := gm.objects.getObjectData(obj).behaviours.get(bhvrType); bhvr != nil {
// 		return bhvr, nil
// 	}
// 	return nil, errors.New(ErrBehaviourNotFound)
// }

// // Get Object's Behaviour by type. Must return, panic if not found.
// func MustGetBehaviour(obj Object, bhvrType Behaviour) Behaviour {
// 	if bhvr := gm.objects.getObjectData(obj).behaviours.get(bhvrType); bhvr != nil {
// 		return bhvr
// 	}
// 	log.Fatalf("[GetBehaviour] Behaviour %T is not found in Object %T", bhvrType, obj)
// 	return nil
// }

// // Get relative's Behaviour by type.
// func GetBehaviourRel(bhvrThis Behaviour, bhvrType Behaviour) (Behaviour, error) {
// 	obj := GetObject(bhvrThis)
// 	if bhvr := gm.objects.getObjectData(obj).behaviours.get(bhvrType); bhvr != nil {
// 		return bhvr, nil
// 	}
// 	return nil, errors.New(ErrBehaviourNotFound)
// }

// Get relative's Behaviour by type. Must return, panic if not found.
func MustGetBehaviourRel(bhvrThis Behaviour, bhvrType Behaviour) Behaviour {
	obj, _ := gm.objects.getParentObjectByBehaviour(bhvrThis)
	objReflectVal := reflect.Indirect(reflect.ValueOf(obj))

	for i := 0; i < objReflectVal.NumField(); i++ {
		field := objReflectVal.Field(i).Addr().Interface()
		if bhvr, ok := field.(Behaviour); ok && (reflect.TypeOf(bhvr) == reflect.TypeOf(bhvrType)) {
			return bhvr
		}
	}
	log.Fatalf("[GetBehaviourRel] Behaviour %T is not found in Object %T. It is required by Behaviour %T.", bhvrType, obj, bhvrThis)
	return nil
}

func initBehaviours(obj Object) {
	objReflectVal := reflect.Indirect(reflect.ValueOf(obj))

	for i := 0; i < objReflectVal.NumField(); i++ {
		if !objReflectVal.Field(i).Addr().Type().Implements(reflect.TypeOf((*Behaviour)(nil)).Elem()) {
			continue
		}
		if !objReflectVal.Field(i).Addr().CanInterface() {
			log.Panicf("Behaviour %v in Object %v should be exported\n", objReflectVal.Field(i).Type().Name(), objReflectVal.Type().Name())
		}
		if bhvr, ok := objReflectVal.Field(i).Addr().Interface().(Behaviour); ok {
			gm.objects.setBehaviour(obj, bhvr)
			bhvr.Init()
		}
	}
}

func updateBehaviours(obj Object) {
	objReflectVal := reflect.Indirect(reflect.ValueOf(obj))

	for i := 0; i < objReflectVal.NumField(); i++ {
		if !objReflectVal.Field(i).Addr().Type().Implements(reflect.TypeOf((*Behaviour)(nil)).Elem()) {
			continue
		}
		if !objReflectVal.Field(i).Addr().CanInterface() {
			log.Panicf("Behaviour %v in Object %v should be exported\n", objReflectVal.Field(i).Type().Name(), objReflectVal.Type().Name())
		}
		if bhvr, ok := objReflectVal.Field(i).Addr().Interface().(Behaviour); ok {
			bhvr.Update()
		}
	}
}

func drawBehaviours(obj Object, screen *ebiten.Image) {
	objReflectVal := reflect.Indirect(reflect.ValueOf(obj))

	for i := 0; i < objReflectVal.NumField(); i++ {
		if !objReflectVal.Field(i).Addr().Type().Implements(reflect.TypeOf((*Behaviour)(nil)).Elem()) {
			continue
		}
		if !objReflectVal.Field(i).Addr().CanInterface() {
			log.Panicf("Behaviour %v in Object %v should be exported\n", objReflectVal.Field(i).Type().Name(), objReflectVal.Type().Name())
		}
		if bhvr, ok := objReflectVal.Field(i).Addr().Interface().(Behaviour); ok {
			bhvr.Draw(screen)
		}
	}
}
