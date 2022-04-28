package gm

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

// TODO: Consider different interface between Behaviour and Object
type Behaviour interface {
	IDerInterface
	Type() int
	Data() BehaviourInstancesData
	PreInit()
	PostInit()
	// PreUpdate()
	// PostUpdate()
	Draw(*ebiten.Image)
}

type BehaviourData struct {
	IDer
}

type BehaviourInstancesData interface {
	IDerInterface
	TypeString() string
	// ByInstance(Object) interface{}
	DelInstance(Object)
	PreUpdate()
	PostUpdate()
	// Draw(*ebiten.Image)
}

type BehaviourInstancesDataData struct {
	IDer
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
// 	log.Panicf("[GetBehaviour] Behaviour %T is not found in Object %T", bhvrType, obj)
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
func MustGetBehaviourRel(bhvrThis Behaviour, bhvrType int) Behaviour {
	return MustGetBehaviour(MustGetObjectParent(bhvrThis), bhvrThis.Type(), bhvrType)
}

// Get Behaviour by type. Must return, panic if not found.
func MustGetBehaviour(instThis Object, bhvrThisType int, bhvrType int) Behaviour {
	bhvrByInst, ok := gm.behaviours.byObjInst[instThis]
	if !ok {
		log.Panicf("[MustGetBehaviour] Behaviour %s is not found in Object %T. It is required by Behaviour %s.", GetBehavioursDataByBhvrType()[bhvrType].TypeString(), instThis, GetBehavioursDataByBhvrType()[bhvrThisType].TypeString())
	}

	bhvrInst, ok := bhvrByInst[bhvrType]
	if !ok {
		log.Panicf("[MustGetBehaviour] Behaviour %s is not found in Object %T. It is required by Behaviour %s.", GetBehavioursDataByBhvrType()[bhvrType].TypeString(), instThis, GetBehavioursDataByBhvrType()[bhvrThisType].TypeString())
	}

	return bhvrInst
}

// Get behaviour's parent. Must return, panic if not found.
func MustGetObjectParent(bhvrThis Behaviour) Object {
	inst, ok := gm.instances.byBhvrInst[bhvrThis]
	if !ok {
		log.Panicf("[MustGetObjectParent] Behaviour %T is not have a parent.", bhvrThis)
	}
	return inst
}

// Get behaviour's parent. Must return, panic if not found.
func MustGetBehavioursByType(bhvrType int) map[Object]Object {
	insts, ok := gm.instances.byBhvrType[bhvrType]
	if !ok {
		log.Panicf("[MustGetObjectParent] Behaviour %s is not exist.", GetBehavioursDataByBhvrType()[bhvrType].TypeString())
	}
	return insts
}
