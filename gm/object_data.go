package gm

import (
	"reflect"
	"sort"
)

// It turns out that map[string] is more expensive than map[int] and map[interface{}]. Theoretically,
// map[int] is O(1), because it has finite domain (bits?), meanwhile, map[string] average is O(N log N),
// because it might apply hashing. So it is preferable to use map[int] instead.
// https://stackoverflow.com/questions/29677670/what-is-the-big-o-performance-of-maps-in-golang
type instances struct {
	lastObjInst int
	byObjInst   map[int]Object
	byObjType   map[string]map[Object]Object // TODO: find away to change this map[string] to map[int]
	byBhvrInst  map[Behaviour]Object
	byBhvrType  map[int]map[Object]Object

	zidxOrdered   []int         // []zidx                 (sorted)
	zidxInstances map[int][]int // map[zidx] []InstanceID (sorted)
}

func GetInstancesByObjInst() map[int]Object               { return gm.instances.byObjInst }
func GetInstancesByObjType() map[string]map[Object]Object { return gm.instances.byObjType }
func GetInstancesByBhvrInst() map[Behaviour]Object        { return gm.instances.byBhvrInst }
func GetInstancesByBhvrType() map[int]map[Object]Object   { return gm.instances.byBhvrType }
func SetInstancesZidx(inst Object, zidx int) {
	if _, ok := gm.instances.zidxInstances[zidx]; !ok {
		// zidxOrdered
		gm.instances.zidxOrdered = append(gm.instances.zidxOrdered, zidx)
		// TODO: remove this sort, use findIntOnSlice instead
		sort.Ints(gm.instances.zidxOrdered)

		// zidxInstance
		gm.instances.zidxInstances[zidx] = make([]int, 0)
	}

	// zidxInstance
	if len(gm.instances.zidxInstances[zidx]) == 0 {
		gm.instances.zidxInstances[zidx] = append(gm.instances.zidxInstances[zidx], inst.ID())
	} else {
		idx, _ := findIntOnSlice(gm.instances.zidxInstances[zidx], inst.ID())
		gm.instances.zidxInstances[zidx] = append(gm.instances.zidxInstances[zidx][:idx], append([]int{inst.ID()}, gm.instances.zidxInstances[zidx][idx:]...)...)
	}
}
func DelInstancesZidx(inst Object) {
	zidx := inst.Zidx()
	if len(gm.instances.zidxInstances[zidx]) == 1 {
		// zidxOrdered
		idx, _ := findIntOnSlice(gm.instances.zidxOrdered, zidx)
		gm.instances.zidxOrdered = append(gm.instances.zidxOrdered[:idx], gm.instances.zidxOrdered[idx+1:]...)

		// zidxInstances
		delete(gm.instances.zidxInstances, zidx)
		return
	}

	// zidxInstances
	idx, _ := findIntOnSlice(gm.instances.zidxInstances[zidx], inst.ID())
	gm.instances.zidxInstances[zidx] = append(gm.instances.zidxInstances[zidx][:idx], gm.instances.zidxInstances[zidx][idx+1:]...)
}
func UpdateInstancesZidx(inst Object, zidxNew int) {
	DelInstancesZidx(inst)
	SetInstancesZidx(inst, zidxNew)
	inst.setZidx(zidxNew)
}

// findIntOnSlice returns the index of match the searched value, if not found it will returns the index of
// the value closest and greater to the searched value. The input array MUST be SORTED.
// TODO: create unit test
// WARNING: watch this! this function is not properly unit tested, it might be returns wrong result
func findIntOnSlice(arr []int, value int) (int, bool) {
	start := 0
	end := len(arr) - 1

	for start <= end {
		mid := (start + end) / 2

		if value == arr[mid] {
			return mid, true
		} else if value < arr[mid] {
			end = mid - 1
		} else {
			start = mid + 1
		}
	}
	return end + 1, false
}

type behaviours struct {
	byBhvrInst map[Behaviour]Behaviour
	byObjInst  map[Object]map[int]Behaviour
	byBhvrType map[int]map[Behaviour]Behaviour
}

func GetBehavioursByBhvrInst() map[Behaviour]Behaviour         { return gm.behaviours.byBhvrInst }
func GetBehavioursByObjInst() map[Object]map[int]Behaviour     { return gm.behaviours.byObjInst }
func GetBehavioursByBhvrType() map[int]map[Behaviour]Behaviour { return gm.behaviours.byBhvrType }

type behavioursData struct {
	byBhvrType map[int]BehaviourInstancesData
}

func GetBehavioursDataByBhvrType() map[int]BehaviourInstancesData {
	return gm.behavioursData.byBhvrType
}

// Set an Instance to Game.
func (g *game) setInstance(inst Object) {
	// objd := ObjectData{object: obj}
	instType := reflect.TypeOf(inst).String()

	// case: instances.byObjInst
	g.instances.byObjInst[g.instances.lastObjInst] = inst
	inst.setID(g.instances.lastObjInst)
	g.instances.lastObjInst++

	// case: instances.byObjType
	if g.instances.byObjType[instType] == nil {
		g.instances.byObjType[instType] = make(map[Object]Object)
	}
	g.instances.byObjType[instType][inst] = inst

	// case: zidx
	SetInstancesZidx(inst, 0)
}

// Delete an Instance from Game.
func (g *game) delInstance(inst Object) {
	instType := reflect.TypeOf(inst).String()

	// case: instances.byObjInst
	delete(g.instances.byObjInst, inst.ID())

	// case: instances.byObjType
	delete(g.instances.byObjType[instType], inst)
	if len(g.instances.byObjType[instType]) == 0 {
		delete(g.instances.byObjType, instType)
	}

	// case: zidx
	DelInstancesZidx(inst)
}

// Set a Behaviour to Game.
func (g *game) setBehaviour(inst Object, bhvrInst Behaviour) {
	// bhvrType := reflect.TypeOf(bhvrInst).String()

	// case: behavioursData.byBhvrType
	bhvrType := bhvrInst.Type()
	// TODO: why if we use this conditional `(bhvrType == 0) || !ok`, it drops fps by almost 50%
	if _, ok := g.behavioursData.byBhvrType[bhvrType]; !ok {
		bhvrType = len(g.behavioursData.byBhvrType) + 1
		bhvrInst.Data().setID(bhvrType)
		g.behavioursData.byBhvrType[bhvrType] = bhvrInst.Data()
	}

	// case: instances.byBhvrInst (set behaviour to parentObject-behvaiour relation on
	// top level game)
	g.instances.byBhvrInst[bhvrInst] = inst

	// case: instances.byBhvrType
	if g.instances.byBhvrType[bhvrType] == nil {
		g.instances.byBhvrType[bhvrType] = make(map[Object]Object)
	}
	g.instances.byBhvrType[bhvrType][inst] = inst

	// case: behaviours.byBhvrInst
	g.behaviours.byBhvrInst[bhvrInst] = bhvrInst

	// case: behaviours.byObjInst
	if g.behaviours.byObjInst[inst] == nil {
		g.behaviours.byObjInst[inst] = make(map[int]Behaviour)
	}
	g.behaviours.byObjInst[inst][bhvrType] = bhvrInst

	// case: behaviours.byBhvrType
	if g.behaviours.byBhvrType[bhvrType] == nil {
		g.behaviours.byBhvrType[bhvrType] = make(map[Behaviour]Behaviour)
	}
	g.behaviours.byBhvrType[bhvrType][bhvrInst] = bhvrInst

}

// Delete a Behaviour from Game.
func (g *game) delBehaviour(inst Object, bhvrInst Behaviour) {
	// bhvrType := reflect.TypeOf(bhvrInst).String()
	bhvrType := bhvrInst.Type()

	// case: instances.byBhvrInst (set behaviour to parentObject-behvaiour relation on
	// top level game)
	delete(g.instances.byBhvrInst, bhvrInst)

	// case: instances.byBhvrType
	delete(g.instances.byBhvrType[bhvrType], inst)
	if len(g.instances.byBhvrType[bhvrType]) == 0 {
		delete(g.instances.byBhvrType, bhvrType)
		// case: behavioursData.byBhvrType
		// delete(bhvrsData.byBhvrType, bhvrType)
	}

	// case: behaviours.byBhvrInst
	delete(g.behaviours.byBhvrInst, bhvrInst)

	// case: behaviours.byObjInst
	delete(g.behaviours.byObjInst[inst], bhvrType)
	if len(g.behaviours.byObjInst[inst]) == 0 {
		delete(g.behaviours.byObjInst, inst)
	}

	// case: behaviours.byBhvrType
	delete(g.behaviours.byBhvrType[bhvrType], bhvrInst)
	if len(g.behaviours.byBhvrType[bhvrType]) == 0 {
		delete(g.behaviours.byBhvrType, bhvrType)
	}

}
