package gm

import (
	"reflect"

	"github.com/pikomonde/gogeta/gogetautil"
)

const (
	maxObjInstID  = 100000 // TODO: should we enforce it?
	maxBhvrTypeID = 50     // TODO: should we enforce it?
)

// It turns out that map[string] is more expensive than map[int] and map[interface{}]. Theoretically,
// map[int] is O(1), because it has finite domain (bits?), meanwhile, map[string] average is O(N log N),
// because it might apply hashing. So it is preferable to use map[int] instead.
// https://stackoverflow.com/questions/29677670/what-is-the-big-o-performance-of-maps-in-golang
type instances struct {
	unused     []int    // []objInstID, up to 0.8 mb
	all        []int    // []objInstID, up to 0.8 mb
	allTypes   []string // [objTypeID]objType, up to 0.8 mb, undeletable, starts from 1
	byObjInst  []Object // [objInstID]Object, up to 0.8 mb, starts from 1
	byObjType  [][]int  // [objTypeID][]objInstID, up to 0.8 mb, starts from 1
	byBhvrInst []int    // [bhvrInstID]objInstID, up to 40 mb, starts from 1
	byBhvrType [][]int  // [bhvrTypeID][]objInstID, up to 40 mb, starts from 1

	zidxOrdered   []int         // []zidx                 (sorted)
	zidxInstances map[int][]int // map[zidx] []InstanceID (sorted)
}

func GetInstIDs() []int                   { return gm.instances.all }
func GetInstTypes() []string              { return gm.instances.allTypes[1:] }
func GetInstByObjInstID(id int) Object    { return gm.instances.byObjInst[id] }
func GetInstIDsByObjTypeID(id int) []int  { return gm.instances.byObjType[id] }
func GetInstIDByBhvrInstID(id int) int    { return gm.instances.byBhvrInst[id] }
func GetInstIDsByBhvrTypeID(id int) []int { return gm.instances.byBhvrType[id] }
func setInstZidx(inst Object, zidx int) {
	if _, ok := gm.instances.zidxInstances[zidx]; !ok {
		// zidxOrdered
		idx, _ := gogetautil.SliceIntOrderedFindIdx(gm.instances.zidxOrdered, zidx)
		gm.instances.zidxOrdered = gogetautil.SliceInsert(gm.instances.zidxOrdered, idx, zidx)

		// zidxInstance
		gm.instances.zidxInstances[zidx] = make([]int, 0)
	}

	// zidxInstance
	if len(gm.instances.zidxInstances[zidx]) == 0 {
		gm.instances.zidxInstances[zidx] = append(gm.instances.zidxInstances[zidx], inst.getID())
	} else {
		idx2, _ := gogetautil.SliceIntOrderedFindIdx(gm.instances.zidxInstances[zidx], inst.getID())
		gm.instances.zidxInstances[zidx] = gogetautil.SliceInsert(gm.instances.zidxInstances[zidx], idx2, inst.getID())
	}
}
func delInstZidx(inst Object) {
	zidx := inst.Zidx()
	if len(gm.instances.zidxInstances[zidx]) == 1 {
		// zidxOrdered
		idx, _ := gogetautil.SliceIntOrderedFindIdx(gm.instances.zidxOrdered, zidx)
		gm.instances.zidxOrdered = append(gm.instances.zidxOrdered[:idx], gm.instances.zidxOrdered[idx+1:]...)

		// zidxInstances
		delete(gm.instances.zidxInstances, zidx)
		return
	}

	// zidxInstances
	idx, _ := gogetautil.SliceIntOrderedFindIdx(gm.instances.zidxInstances[zidx], inst.getID())
	gm.instances.zidxInstances[zidx] = gogetautil.SliceCut(gm.instances.zidxInstances[zidx], idx)
}
func updateInstZidx(inst Object, zidxNew int) {
	delInstZidx(inst)
	setInstZidx(inst, zidxNew)
	inst.setZidx(zidxNew)
}

type behaviours struct {
	unused     []int       // []bhvrInstID, up to 40 mb
	all        []int       // []bhvrInstID, up to 40 mb
	allTypes   []string    // [bhvrTypeID]bhvrType, up to 0,8 mb, undeletable, starts from 1
	byBhvrInst []Behaviour // [bhvrInstID]Behaviour, up to 40 mb, starts from 1
	byBhvrType [][]int     // [bhvrTypeID][]bhvrInstID, up to 40 mb, starts from 1
	byObjInst  [][]int     // [objInstID][]bhvrInstID, up to 40 mb, starts from 1
}

func GetBhvrIDs() []int                    { return gm.behaviours.all }
func GetBhvrTypes() []string               { return gm.behaviours.allTypes[1:] }
func GetBhvrByBhvrInstID(id int) Behaviour { return gm.behaviours.byBhvrInst[id] }
func GetBhvrIDsByBhvrTypeID(id int) []int  { return gm.behaviours.byBhvrType[id] }
func GetBhvrIDsByObjInstID(id int) []int   { return gm.behaviours.byObjInst[id] }

type behavioursData struct {
	byBhvrType []BehavioursData // [bhvrTypeID]BehavioursData, up to 200 b, undeletable, starts from 1
}

func GetBhvrDatas() []BehavioursData                { return gm.behavioursData.byBhvrType[1:] }
func GetBhvrDataByBhvrTypeID(id int) BehavioursData { return gm.behavioursData.byBhvrType[id] }

// Set an Instance to Game.
func (g *game) setInstance(objInst Object) {
	var objInstID int
	if len(g.instances.unused) > 0 {
		// objInstID
		objInstID = g.instances.unused[0]
		g.instances.unused = g.instances.unused[1:]

		// instances.byObjInst
		g.instances.byObjInst[objInstID] = objInst

		// behaviours.byObjInst part A
		g.behaviours.byObjInst[objInstID] = make([]int, 0)

	} else {
		// objInstID
		objInstID = len(g.instances.byObjInst)

		// instances.byObjInst
		g.instances.byObjInst = append(g.instances.byObjInst, objInst)

		// behaviours.byObjInst part A
		g.behaviours.byObjInst = append(g.behaviours.byObjInst, []int{})

	}
	objInst.setID(objInstID)

	// instances.all
	idx1, _ := gogetautil.SliceIntOrderedFindIdx(g.instances.all, objInstID)
	g.instances.all = gogetautil.SliceInsert(g.instances.all, idx1, objInstID)

	objType := reflect.TypeOf(objInst).String()
	objTypeID := gogetautil.SliceStringFindIdx(g.instances.allTypes, objType)
	if objTypeID < 0 {
		// instances.allTypes
		objTypeID = len(g.instances.allTypes)
		g.instances.allTypes = append(g.instances.allTypes, objType)

		// instances.byObjType
		g.instances.byObjType = append(g.instances.byObjType, []int{objInstID})

	} else {
		// instances.byObjType
		idx2, _ := gogetautil.SliceIntOrderedFindIdx(g.instances.byObjType[objTypeID], objInstID)
		g.instances.byObjType[objTypeID] = gogetautil.SliceInsert(g.instances.byObjType[objTypeID], idx2, objInstID)

	}
	objInst.setType(objType)
	objInst.setTypeID(objTypeID)

	// zidx
	setInstZidx(objInst, 0)
}

// Delete an Instance from Game.
func (g *game) delInstance(objInst Object) {
	objInstID := objInst.getID()
	g.instances.unused = append(g.instances.unused, objInstID)
	objTypeID := objInst.getTypeID()

	// instances.all
	idx1, _ := gogetautil.SliceIntOrderedFindIdx(g.instances.all, objInstID)
	g.instances.all = gogetautil.SliceCut(g.instances.all, idx1)

	// instances.allTypes
	// instances.allTypes is undeletable

	// instances.byObjInst
	g.instances.byObjInst[objInstID] = nil

	// instances.byObjType
	idx2, _ := gogetautil.SliceIntOrderedFindIdx(g.instances.byObjType[objTypeID], objInstID)
	g.instances.byObjType[objTypeID] = gogetautil.SliceCut(g.instances.byObjType[objTypeID], idx2)

	// zidx
	delInstZidx(objInst)
}

// Set a Behaviour to Game.
func (g *game) setBehaviour(objInst Object, bhvrInst Behaviour) {
	objInstID := objInst.getID()
	var bhvrInstID int
	if len(g.behaviours.unused) > 0 {
		// bhvrInstID
		bhvrInstID = g.behaviours.unused[0]
		g.behaviours.unused = g.behaviours.unused[1:]

		// instances.byBhvrInst
		g.instances.byBhvrInst[bhvrInstID] = objInstID

		// behaviours.byBhvrInst
		g.behaviours.byBhvrInst[bhvrInstID] = bhvrInst

	} else {
		// bhvrInstID
		bhvrInstID = len(g.behaviours.byBhvrInst)

		// instances.byBhvrInst
		g.instances.byBhvrInst = append(g.instances.byBhvrInst, objInstID)

		// behaviours.byBhvrInst
		g.behaviours.byBhvrInst = append(g.behaviours.byBhvrInst, bhvrInst)

	}
	bhvrInst.setID(bhvrInstID)

	// behaviours.all
	idx1, _ := gogetautil.SliceIntOrderedFindIdx(g.behaviours.all, bhvrInstID)
	g.behaviours.all = gogetautil.SliceInsert(g.behaviours.all, idx1, bhvrInstID)

	bhvrType := reflect.TypeOf(bhvrInst).String()
	bhvrTypeID := gogetautil.SliceStringFindIdx(g.behaviours.allTypes, bhvrType)
	if bhvrTypeID < 0 {
		// behaviours.allTypes
		bhvrTypeID = len(g.behaviours.allTypes)
		g.behaviours.allTypes = append(g.behaviours.allTypes, bhvrType)

		// instances.byBhvrType
		g.instances.byBhvrType = append(g.instances.byBhvrType, []int{objInstID})

		// behaviours.byBhvrType
		g.behaviours.byBhvrType = append(g.behaviours.byBhvrType, []int{bhvrInstID})

		// behavioursData.byBhvrType
		g.behavioursData.byBhvrType = append(g.behavioursData.byBhvrType, bhvrInst.Data())

		// behavioursData instance
		bhvrInst.Data().setID(bhvrTypeID)
		bhvrInst.Data().setType(bhvrType)
		bhvrInst.Data().setTypeID(bhvrTypeID)

	} else {
		// instances.byBhvrType
		idx2, _ := gogetautil.SliceIntOrderedFindIdx(g.instances.byBhvrType[bhvrTypeID], objInstID)
		g.instances.byBhvrType[bhvrTypeID] = gogetautil.SliceInsert(g.instances.byBhvrType[bhvrTypeID], idx2, objInstID)

		// behaviours.byBhvrType
		idx3, _ := gogetautil.SliceIntOrderedFindIdx(g.behaviours.byBhvrType[bhvrTypeID], bhvrInstID)
		g.behaviours.byBhvrType[bhvrTypeID] = gogetautil.SliceInsert(g.behaviours.byBhvrType[bhvrTypeID], idx3, bhvrInstID)

	}
	bhvrInst.setType(bhvrType)
	bhvrInst.setTypeID(bhvrTypeID)

	// behaviours.byObjInst part B
	g.behaviours.byObjInst[objInstID] = append(g.behaviours.byObjInst[objInstID], bhvrInstID)

}

// Delete a Behaviour from Game.
func (g *game) delBehaviour(objInst Object, bhvrInst Behaviour) {
	objInstID := objInst.getID()
	bhvrInstID := bhvrInst.getID()
	g.behaviours.unused = append(g.behaviours.unused, bhvrInstID)
	bhvrTypeID := bhvrInst.getTypeID()

	// instances.byBhvrInst
	g.instances.byBhvrInst[bhvrInstID] = 0

	// instances.byBhvrType
	idx1, _ := gogetautil.SliceIntOrderedFindIdx(g.instances.byBhvrType[bhvrTypeID], objInstID)
	g.instances.byBhvrType[bhvrTypeID] = gogetautil.SliceCut(g.instances.byBhvrType[bhvrTypeID], idx1)

	// behaviours.all
	idx2, _ := gogetautil.SliceIntOrderedFindIdx(g.behaviours.all, bhvrInstID)
	g.behaviours.all = gogetautil.SliceCut(g.behaviours.all, idx2)

	// behaviours.allTypes
	// behaviours.allTypes is undeletable

	// behaviours.byBhvrInst
	g.behaviours.byBhvrInst[bhvrInstID] = nil

	// behaviours.byBhvrType
	idx3, _ := gogetautil.SliceIntOrderedFindIdx(g.behaviours.byBhvrType[bhvrTypeID], bhvrInstID)
	g.behaviours.byBhvrType[bhvrTypeID] = gogetautil.SliceCut(g.behaviours.byBhvrType[bhvrTypeID], idx3)

	// behaviours.byObjInst
	idx4, _ := gogetautil.SliceIntOrderedFindIdx(g.behaviours.byObjInst[objInstID], bhvrInstID)
	g.behaviours.byObjInst[objInstID] = gogetautil.SliceCut(g.behaviours.byObjInst[objInstID], idx4)

	// behavioursData.byBhvrType
	// behavioursData.byBhvrType is undeletable

}
