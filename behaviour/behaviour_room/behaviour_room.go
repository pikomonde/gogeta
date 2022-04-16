package behaviour_room

import (
	"github.com/hajimehoshi/ebiten/v2"
	bhvrCommon "github.com/pikomonde/gogeta/behaviour/behaviour_common"
	"github.com/pikomonde/gogeta/gm"
)

type Room struct {
	Instances []gm.Object
}

func (bhvr *Room) PreInit() {
	bhvr.Instances = make([]gm.Object, 0)
}

func (bhvr *Room) PostInit() {
}

func (bhvr *Room) PreUpdate() {
}

func (bhvr *Room) PostUpdate() {
}

func (bhvr *Room) Draw(screen *ebiten.Image) {
}

// === Behaviour specific method ===

func (bhvr *Room) InitObject(obj gm.Object) gm.Object {
	initObj := gm.InitObject(obj)
	bhvr.Instances = append(bhvr.Instances, initObj)
	return initObj
}

func (bhvr *Room) Restart() {
	for _, junkInst := range bhvr.Instances {
		gm.DelObject(junkInst)
	}
	bhvr.Instances = make([]gm.Object, 0)
	gm.MustGetObjectParent(bhvr).Init()
}

func (bhvr *Room) Resume() {
	for _, inst := range bhvr.Instances {
		bCommon := gm.MustGetBehaviour(inst, bhvr, &bhvrCommon.Common{}).(*bhvrCommon.Common)
		bCommon.IsStopUpdate = false
	}
}

func (bhvr *Room) Pause() {
	for _, inst := range bhvr.Instances {
		bCommon := gm.MustGetBehaviour(inst, bhvr, &bhvrCommon.Common{}).(*bhvrCommon.Common)
		bCommon.IsStopUpdate = true
	}
}

// === Package functions ===
