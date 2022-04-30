package gm

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

// TODO: Consider different interface between Behaviour and Object
type Behaviour interface {
	Instance
	Data() BehavioursData
	PreInit()
	PostInit()
	// PreUpdate()
	// PostUpdate()
	Draw(*ebiten.Image)
}

type BehavioursData interface {
	Instance
	Behaviour() Behaviour
	// ByInstance(Object) interface{}
	DelInstance(Object)
	PreUpdate()
	PostUpdate()
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
	return MustGetBehaviour(MustGetObjectParent(bhvrThis), bhvrThis.getTypeID(), bhvrType)
}

// Get Behaviour by type. Must return, panic if not found.
func MustGetBehaviour(instThis Object, bhvrThisTypeID int, bhvrTypeID int) Behaviour {
	for _, bhvrID := range GetBhvrIDsByObjInstID(instThis.getID()) {
		if GetBhvrByBhvrInstID(bhvrID).getTypeID() == bhvrTypeID {
			return GetBhvrByBhvrInstID(bhvrID)
		}
	}
	log.Panicf("[MustGetBehaviour] Behaviour %T is not found in Object %T. It is required by Behaviour %T.",
		GetBhvrDataByBhvrTypeID(bhvrTypeID), instThis, GetBhvrDataByBhvrTypeID(bhvrThisTypeID))
	return nil
}

// Get behaviour's parent. Must return, panic if not found.
func MustGetObjectParent(bhvrThis Behaviour) Object {
	return GetInstByObjInstID(GetInstIDByBhvrInstID(bhvrThis.getID()))
}
