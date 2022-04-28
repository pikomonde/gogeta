package behaviour_common

import (
	"github.com/pikomonde/gogeta/gm"
)

var Data BehaviourInstancesData

func init() {
	Data.instancesData = make(map[gm.Object]*InstanceData)
}

// === instances data ===

type BehaviourInstancesData struct {
	gm.BehaviourInstancesDataData
	instancesData map[gm.Object]*InstanceData
}

func (data *BehaviourInstancesData) TypeString() string { return "Common" }

func (data *BehaviourInstancesData) ByInstance(instance gm.Object) *InstanceData {
	return data.instancesData[instance]
}

func (data *BehaviourInstancesData) DelInstance(instance gm.Object) {
	delete(data.instancesData, instance)
}

func (data *BehaviourInstancesData) PreUpdate() {
	for _, inst := range gm.GetInstancesByBhvrType()[Common{}.Type()] {
		if !inst.IsUpdate() {
			continue
		}

		bhvrCommonRaw, ok := gm.GetBehavioursByObjInst()[inst][Common{}.Type()]
		if !ok {
			continue
		}

		bhvrCommonInst := bhvrCommonRaw.(*Common)
		// TODO: using delta time instead per tick for stability
		// ebiten.MaxTPS()
		bhvrCommonInst.Position.X += bhvrCommonInst.Speed.X
		bhvrCommonInst.Position.Y += bhvrCommonInst.Speed.Y
	}
	// // TODO: move this to other behaviour
	// bhvr.Angle += 0.01
	// // inpututil.IsKeyJustPressed
	// if ebiten.IsKeyPressed(ebiten.KeyA) {
	// 	bhvr.Position.X -= 3.5
	// }
	// if ebiten.IsKeyPressed(ebiten.KeyD) {
	// 	bhvr.Position.X += 3.5
	// }
	// if ebiten.IsKeyPressed(ebiten.KeyW) {
	// 	bhvr.Position.Y -= 3.5
	// }
	// if ebiten.IsKeyPressed(ebiten.KeyS) {
	// 	bhvr.Position.Y += 3.5
	// }
}

func (data *BehaviourInstancesData) PostUpdate() {

}

// it is a very expensive function
// desktop ebiten                    : 68.900
// desktop without                   : 54.000
// desktop with (string & sprintf)   : 12.000
// desktop with (int)                : 24.000
// desktop with (int) & map[int]     : 36.000

// === instance data ===

type InstanceData struct{}
