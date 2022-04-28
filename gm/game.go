package gm

import (
	"github.com/hajimehoshi/ebiten/v2"
)

var gm game

type game struct {
	instances        instances      // all instances in the game, type "Object" is "object interface pointer"
	behaviours       behaviours     // all behaviours in the game, type "Behaviours" is "behaviour interface pointer"
	behavioursData   behavioursData //
	layoutW, layoutH int
}

// TODO: restore this function for debugging
// func GetObjectDB() objects {
// 	return gm.objects
// }

// Update all Instances and BehavioursData.
func (g *game) Update() error {
	for bhvrType, bhvrInstData := range g.behavioursData.byBhvrType {
		if _, ok := gm.behavioursData.byBhvrType[bhvrType]; ok {
			bhvrInstData.PreUpdate()
		}
	}

	for _, inst := range g.instances.byObjInst {
		if _, ok := gm.instances.byObjInst[inst.ID()]; ok {
			if inst.IsUpdate() {
				inst.Update()
			}
		}
	}

	for bhvrType, bhvrInstData := range g.behavioursData.byBhvrType {
		if _, ok := gm.behavioursData.byBhvrType[bhvrType]; ok {
			bhvrInstData.PostUpdate()
		}
	}
	return nil
}

// Draw all Instances and BehavioursData.
func (g *game) Draw(screen *ebiten.Image) {
	// Draw object
	// objArr := make([]string, 0)
	// objDataMap := make(map[string]ObjectData)
	// for obj, objData := range objs[KeybyObj] {
	// 	objStr := fmt.Sprintf("%016.10f:%p", objData.ZIdx, obj)
	// 	// fmt.Println(objStr)
	// 	objArr = append(objArr, objStr)
	// 	objDataMap[objStr] = objData
	// }

	// sort.Strings(objArr)

	// for _, str := range objArr {
	// 	objData := objDataMap[str]
	// 	drawBehaviours(objData.object, screen)
	// 	objData.object.Draw(screen)
	// }

	// for bhvrType, bhvrInstData := range g.behavioursData.byBhvrType {
	// 	if _, ok := gm.behavioursData.byBhvrType[bhvrType]; ok {
	// 		bhvrInstData.Draw(screen)
	// 	}
	// }

	// for _, inst := range g.instances.byDrawOrder {
	// 	if _, ok := gm.instances.byObjInst[inst.ID()]; ok {
	// 		inst.Draw(screen)
	// 	}
	// }

	// ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%+v", g.instances.zidxOrdered), 100, 100)
	for _, zidx := range g.instances.zidxOrdered {
		for _, instID := range g.instances.zidxInstances[zidx] {
			if inst, ok := gm.instances.byObjInst[instID]; ok {
				if inst.IsDraw() {
					for _, bhvr := range GetBehavioursByObjInst()[inst] {
						bhvr.Draw(screen)
					}
					inst.Draw(screen)
				}
			}
		}
	}
}

func (g *game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.layoutW, g.layoutH
}

func Init(windowW, windowH, layoutW, layoutH int) error {
	gm.instances.byObjInst = make(map[int]Object)
	gm.instances.byObjType = make(map[string]map[Object]Object)
	gm.instances.byBhvrInst = make(map[Behaviour]Object)
	gm.instances.byBhvrType = make(map[int]map[Object]Object)

	gm.instances.zidxInstances = make(map[int][]int)
	gm.instances.zidxOrdered = make([]int, 0)

	gm.behaviours.byObjInst = make(map[Object]map[int]Behaviour)
	gm.behaviours.byBhvrInst = make(map[Behaviour]Behaviour)
	gm.behaviours.byBhvrType = make(map[int]map[Behaviour]Behaviour)

	gm.behavioursData.byBhvrType = make(map[int]BehaviourInstancesData, 0)

	gm.layoutW, gm.layoutH = layoutW, layoutH
	ebiten.SetWindowSize(windowW, windowH)
	return nil
}

func Run() error {
	return ebiten.RunGame(&gm)
}
