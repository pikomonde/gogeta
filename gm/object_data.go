package gm

import (
	"fmt"
	"reflect"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
)

type ObjectData struct {
	object       Object
	ZIdx         float64
	IsStopUpdate bool
}

// Example 1: all object pointer & finding by exact
//            object pointer								map[obj]				map["*objPointer"]ObjectData
//
// Example 2: finding by object type						map[objT::"objType"]	map["*objPointer"]ObjectData
// objects stored ObjectData with mapping by objectType
// and objectName. For example. there are 2 "enemy" objects
// and 1 "player" object, first map indexed by "enemy" and
// "player", meanwhile the second mapping contain the
// address (interface) of "enemy 1", "enemy 2", and
// "player 1".
//
// Example 3: finding by exact behaviour pointer (for
//            parent object case)							map[bhvr::"*bhvPointer"]map["*objPointer"]objectData
//
// Example 4: finding by exact behaviour type				map[bhvr::"bhvType"]	map["*objPointer"]objectData
//
// Example 5: finding by z-index							map[zidx]				map["*objPointer"]objectData
type objects map[string]map[Object]ObjectData

const (
	KeyByObj      = "obj"
	KeyByObjType  = "objT::"
	KeyByBhvr     = "bhvr::"
	KeyByBhvrType = "bhvrT::"
	// KeyByZIdx     = "zidx"
)

// Set an Instance to Game.
func (objs objects) setObject(obj Object) {
	objd := ObjectData{object: obj}
	objType := reflect.TypeOf(obj).String()

	// case: KeyByObj
	key := KeyByObj
	if _, ok := objs[key]; !ok {
		objs[key] = make(map[Object]ObjectData)
	}
	objs[key][obj] = objd

	// case: KeyByObjType
	key = fmt.Sprintf("%s%s", KeyByObjType, objType)
	if _, ok := objs[key]; !ok {
		objs[key] = make(map[Object]ObjectData)
	}
	objs[key][obj] = objd
}

// Delete an Instance from Game.
func (objs objects) delObject(obj Object) {
	objType := reflect.TypeOf(obj).String()

	// case: KeyByObj
	key := KeyByObj
	delete(objs[key], obj)

	// case: KeyByObjType
	key = fmt.Sprintf("%s%s", KeyByObjType, objType)
	delete(objs[key], obj)
	if len(objs[key]) == 0 {
		delete(objs, key)
	}
}

// Update all Instances.
func (objs objects) update() {
	// TODO: behaviours should be updated outside from this loop
	for _, objData := range objs[KeyByObj] {
		if !objData.IsStopUpdate {
			preUpdateBehaviours(objData.object)
			objData.object.Update()
			postUpdateBehaviours(objData.object)
		}
	}
}

// Draw all Instances.
func (objs objects) draw(screen *ebiten.Image) {
	// TODO: consider z-index when draw objects
	// TODO: behaviours should be draw outside from this loop
	objArr := make([]string, 0)
	objDataMap := make(map[string]ObjectData)
	for obj, objData := range objs[KeyByObj] {
		objStr := fmt.Sprintf("%016.10f:%p", objData.ZIdx, obj)
		// fmt.Println(objStr)
		objArr = append(objArr, objStr)
		objDataMap[objStr] = objData
	}

	sort.Strings(objArr)

	for _, str := range objArr {
		objData := objDataMap[str]
		drawBehaviours(objData.object, screen)
		objData.object.Draw(screen)
	}

}

func (objs objects) MustSetObjectData(obj Object, objData ObjectData) {
	if _, ok := objs[KeyByObj]; !ok {
		return
	}

	if _, ok := objs[KeyByObj][obj]; !ok {
		return
	}

	objs[KeyByObj][obj] = ObjectData{
		object:       objs[KeyByObj][obj].object,
		ZIdx:         objData.ZIdx,
		IsStopUpdate: objData.IsStopUpdate,
	}
}
