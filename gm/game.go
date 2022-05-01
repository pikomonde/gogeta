package gm

import (
	"fmt"

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

func Println() string {
	return fmt.Sprintf(`> %2f %2f %d
instances.unused:        %+v
instances.all:           %+v
instances.allTypes:      %+v
instances.byObjInst:     %+v
instances.byObjType:     %+v
instances.byBhvrInst:    %+v
instances.byBhvrType:    %+v
instances.zidxOrdered:   %+v
instances.zidxInstances: %+v
behaviours.unused:       %+v
behaviours.all:          %+v
behaviours.allTypes:     %+v
behaviours.byBhvrInst:   %+v
behaviours.byBhvrType:   %+v
behaviours.byObjInst:    %+v
behavioursData.byBhvrType:  %+v`,
		ebiten.CurrentTPS(), ebiten.CurrentFPS(), len(GetInstIDs()),
		gm.instances.unused,
		gm.instances.all,
		gm.instances.allTypes,
		gm.instances.byObjInst,
		gm.instances.byObjType,
		gm.instances.byBhvrInst,
		gm.instances.byBhvrType,
		gm.instances.zidxOrdered,
		gm.instances.zidxInstances,
		gm.behaviours.unused,
		gm.behaviours.all,
		gm.behaviours.allTypes,
		gm.behaviours.byBhvrInst,
		gm.behaviours.byBhvrType,
		gm.behaviours.byObjInst,
		gm.behavioursData.byBhvrType,
	)
}

// Update all Instances and BehavioursData.
func (g *game) Update() error {
	for _, bhvrInstData := range GetBhvrDatas() {
		bhvrInstData.PreUpdate()
	}

	for _, objInstID := range GetInstIDs() {
		if inst := GetInstByObjInstID(objInstID); inst != nil {
			if inst.IsUpdate() {
				inst.Update()
			}
		}
	}

	for _, bhvrInstData := range GetBhvrDatas() {
		bhvrInstData.PostUpdate()
	}

	// // delete (garbage collected) unused instance in a slice
	// g.instances.all = gogetautil.SliceCutZeros(g.instances.all)
	// for i := range g.instances.byObjType {
	// 	gogetautil.SliceCutZeros(g.instances.byObjType[i])
	// }
	// for i := range g.instances.byBhvrType {
	// 	gogetautil.SliceCutZeros(g.instances.byBhvrType[i])
	// }

	// g.behaviours.all = gogetautil.SliceCutZeros(g.behaviours.all)
	// for i := range g.behaviours.byBhvrType {
	// 	gogetautil.SliceCutZeros(g.behaviours.byBhvrType[i])
	// }
	// for i := range g.behaviours.byObjInst {
	// 	gogetautil.SliceCutZeros(g.behaviours.byObjInst[i])
	// }

	return nil
}

// Draw all Instances and BehavioursData.
func (g *game) Draw(screen *ebiten.Image) {
	for _, zidx := range g.instances.zidxOrdered {
		for _, instID := range g.instances.zidxInstances[zidx] {
			if inst := GetInstByObjInstID(instID); inst != nil {
				if inst.IsDraw() {
					for _, bhvrInstID := range GetBhvrIDsByObjInstID(instID) {
						GetBhvrByBhvrInstID(bhvrInstID).Draw(screen)
					}
					inst.Draw(screen)
				}
			}
		}
	}
	// ebitenutil.DebugPrintAt(screen, Println(), 20, 100)
}

func (g *game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.layoutW, g.layoutH
}

func Init(windowW, windowH, layoutW, layoutH int) error {
	gm.instances.unused = make([]int, 0)
	gm.instances.all = make([]int, 0)
	gm.instances.allTypes = []string{""}
	gm.instances.byObjInst = []Object{nil}
	gm.instances.byObjType = [][]int{nil}
	gm.instances.byBhvrInst = []int{0}
	gm.instances.byBhvrType = [][]int{nil}

	gm.instances.zidxInstances = make(map[int][]int)
	gm.instances.zidxOrdered = make([]int, 0)

	gm.behaviours.unused = make([]int, 0)
	gm.behaviours.all = make([]int, 0)
	gm.behaviours.allTypes = []string{""}
	gm.behaviours.byBhvrInst = []Behaviour{nil}
	gm.behaviours.byObjInst = [][]int{nil}
	gm.behaviours.byBhvrType = [][]int{nil}

	gm.behavioursData.byBhvrType = []BehavioursData{nil}

	gm.layoutW, gm.layoutH = layoutW, layoutH
	ebiten.SetWindowSize(windowW, windowH)
	return nil
}

func Run() error {
	return ebiten.RunGame(&gm)
}
