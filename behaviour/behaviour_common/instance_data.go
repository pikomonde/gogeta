package behaviour_common

import (
	"github.com/pikomonde/gogeta/gm"
)

var Data CommonsData

func init() {
	Data.instancesData = make(map[gm.Object]*InstanceData)
}

// === instances data ===

type CommonsData struct {
	gm.Instancer
	instancesData map[gm.Object]*InstanceData
}

func (data *CommonsData) Behaviour() gm.Behaviour                 { return &Common{} }
func (data *CommonsData) ByInstance(indt gm.Object) *InstanceData { return data.instancesData[indt] }
func (data *CommonsData) DelInstance(indt gm.Object)              { delete(data.instancesData, indt) }

func (data *CommonsData) PreUpdate() {
	for _, bhvrInstID := range gm.GetBhvrIDsByBhvrTypeID(gm.TypeID(&Common{})) {
		instID := gm.GetInstIDByBhvrInstID(bhvrInstID)
		inst := gm.GetInstByObjInstID(instID)
		if !inst.IsUpdate() {
			continue
		}
		bhvrCommonInst := gm.GetBhvrByBhvrInstID(bhvrInstID).(*Common)
		// TODO: using delta time instead per tick for stability.
		// ebiten.MaxTPS(). The downside is we can't reproducing step by step.
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

func (data *CommonsData) PostUpdate() {

}

// it is a very expensive function
// desktop ebiten                    : 68.900
// desktop without                   : 54.000
// desktop with (string & sprintf)   : 12.000
// desktop with (int)                : 24.000
// desktop with (int) & map[int]     : 36.000

// === instance data ===

type InstanceData struct{}
