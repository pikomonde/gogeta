package gm

import (
	"reflect"

	"github.com/pikomonde/gogeta/gogetautil"
)

// === Instancer ===

type Instance interface {
	getID() int
	setID(int)
	getType() string
	setType(string)
	getTypeID() int
	setTypeID(int)
}

type Instancer struct {
	id     int
	typ    string
	typeID int
}

func (ins *Instancer) getID() int       { return ins.id }
func (ins *Instancer) setID(d int)      { ins.id = d }
func (ins *Instancer) getType() string  { return ins.typ }
func (ins *Instancer) setType(d string) { ins.typ = d }
func (ins *Instancer) getTypeID() int   { return ins.typeID }
func (ins *Instancer) setTypeID(d int)  { ins.typeID = d }

// TODO: better alternative than ID(ins) or TypeID(ins) directly from instance using struct embedding?
// using reflect?

// ID returns the instance's ID. It is not safe for use on search query. If it is never been set on
// game_data before, it will return 0.
func ID(ins Instance) int {
	// TODO: unit test here
	// general
	if ins.getID() != 0 {
		return ins.getID()
	}

	// never been set before
	if _, ok := ins.(BehavioursData); ok {
		// id and typeID is same for BehavioursData, get logic from there
		return TypeID(ins)
	}

	return 0
}

// Type returns the instance's Type. It is a guarantee that it will not return empty string, but doesn't
// mean that this already set on game_data. If it is an instance of Object, it is better to use it once,
// and cache it on some variable, because it is VERY EXPENSIVE (using reflect). If it is a Behaviour or
// BehavioursData, you don't need to cache it, because it will autaomatically cached on BehavioursData.
func Type(ins Instance) string {
	// TODO: unit test here
	// general
	if ins.getType() != "" {
		return ins.getType()
	}

	// never been set before
	if _, ok := ins.(Object); ok {
		// TODO: log warn here, because reflect can be expensive, especially because the result is not
		// cached automatically.
		ins.setType(reflect.TypeOf(ins).String())
		return ins.getType()
	} else if bhvrInst, ok := ins.(Behaviour); ok {
		if bhvrInst.Data().getType() != "" {
			ins.setType(bhvrInst.Data().getType())
			return ins.getType()
		}
		ins.setType(reflect.TypeOf(ins).String())
		bhvrInst.Data().setType(ins.getType())
		return ins.getType()
	} else if bhvrDataInst, ok := ins.(BehavioursData); ok {
		ins.setType(reflect.TypeOf(bhvrDataInst.Behaviour()).String())
		return ins.getType()
	}

	return ""
}

// TypeID returns the instance's TypeID. It is a guarantee that it will not return 0. It will
// automatically set the Type, TypeID, and BehavioursData to the game_data. If it is an instance of
// Object, it is better to use it once, and cache it on some variable, because it is VERY EXPENSIVE
// (using reflect). If it is a Behaviour or BehavioursData, you don't need to cache it, because it will
// automatically cached on BehavioursData.
func TypeID(ins Instance) int {
	// TODO: unit test here
	// general
	if ins.getTypeID() != 0 {
		return ins.getTypeID()
	}

	// never been set before
	if _, ok := ins.(Object); ok {
		typeID := gogetautil.SliceStringFindIdx(gm.instances.allTypes, Type(ins))
		if typeID < 0 {
			// instances.allTypes
			typeID = len(gm.instances.allTypes)
			gm.instances.allTypes = append(gm.instances.allTypes, Type(ins))

			// instances.byObjType
			gm.instances.byObjType = append(gm.instances.byObjType, []int{})

		}
		return typeID
	} else if bhvrInst, ok := ins.(Behaviour); ok {
		typeID := gogetautil.SliceStringFindIdx(gm.behaviours.allTypes, Type(ins))
		if typeID < 0 {
			// behaviours.allTypes
			typeID = len(gm.behaviours.allTypes)
			gm.behaviours.allTypes = append(gm.behaviours.allTypes, Type(ins))

			// instances.byBhvrType
			gm.instances.byBhvrType = append(gm.instances.byBhvrType, []int{})

			// behaviours.byBhvrType
			gm.behaviours.byBhvrType = append(gm.behaviours.byBhvrType, []int{})

			// behavioursData.byBhvrType
			gm.behavioursData.byBhvrType = append(gm.behavioursData.byBhvrType, bhvrInst.Data())

			ins.setType(Type(ins))
			ins.setTypeID(typeID)
			bhvrInst.Data().setID(typeID)
			bhvrInst.Data().setType(Type(ins))
			bhvrInst.Data().setTypeID(typeID)
		}
		return typeID
	} else if bhvrsDataInst, ok := ins.(BehavioursData); ok {
		typeID := gogetautil.SliceStringFindIdx(gm.behaviours.allTypes, Type(ins))
		if typeID < 0 {
			// behaviours.allTypes
			typeID = len(gm.behaviours.allTypes)
			gm.behaviours.allTypes = append(gm.behaviours.allTypes, Type(ins))

			// instances.byBhvrType
			gm.instances.byBhvrType = append(gm.instances.byBhvrType, []int{})

			// behaviours.byBhvrType
			gm.behaviours.byBhvrType = append(gm.behaviours.byBhvrType, []int{})

			// behavioursData.byBhvrType
			gm.behavioursData.byBhvrType = append(gm.behavioursData.byBhvrType, bhvrsDataInst)

			ins.setID(typeID)
			ins.setType(Type(ins))
			ins.setTypeID(typeID)
		}
		return typeID
	}

	return 0
}
