package behaviour_room

import (
	"log"
	"reflect"

	"github.com/golang/geo/r1"
	"github.com/golang/geo/r2"
	"github.com/hajimehoshi/ebiten/v2"
	bhvrCommon "github.com/pikomonde/gogeta/behaviour/behaviour_common"
	"github.com/pikomonde/gogeta/gm"
)

type Room struct {
	gm.Instancer
	Position  r2.Point
	Size      r2.Point
	instances map[gm.Object]gm.Object
}

func (bhvr *Room) Data() gm.BehavioursData { return &Data }

func (bhvr *Room) PreInit() {
	bhvr.instances = make(map[gm.Object]gm.Object)
}

func (bhvr *Room) PostInit() {
	if (bhvr.Size.X <= 0) || (bhvr.Size.Y <= 0) {
		log.Panicf("Room size cannot be empty")
	}
}

func (bhvr *Room) Draw(screen *ebiten.Image) {

}

func (bhvr *Room) Instances() map[gm.Object]gm.Object { return bhvr.instances }

// TODO: what to do with room's window? is it as view camera? do we still need game's canvas width/height?
// is there any limitation on "active" room (because we need to attach camera to a room)? should we attach
// a camera to a room? how camera works?
func (bhvr *Room) window() r2.Rect {
	return r2.Rect{
		X: r1.Interval{Lo: bhvr.Position.X, Hi: bhvr.Position.X + bhvr.Size.X},
		Y: r1.Interval{Lo: bhvr.Position.Y, Hi: bhvr.Position.Y + bhvr.Size.Y},
	}
}

// === Behaviour specific method ===

// InitObject register current instance to this room. It is also register to gm
// through gm.InitObject().
func (bhvr *Room) InitObject(instance gm.Object, data InstanceData) gm.Object {
	if _, ok := Data.instancesData[gm.ID(instance)]; ok {
		instance := gm.MustGetObjectParent(bhvr)
		instanceName := reflect.TypeOf(instance).String()
		log.Panicf("This instance is already assigned to room %s", instanceName)
	}
	initObj := gm.InitObject(instance)
	bhvr.instances[initObj] = initObj
	Data.instancesData[gm.ID(instance)] = &InstanceData{
		roomInstance:        gm.MustGetObjectParent(bhvr),
		room:                bhvr,
		IsDeleteWhenOutside: data.IsDeleteWhenOutside,
	}
	return initObj
}

func (bhvr *Room) Restart() {
	for _, inst := range bhvr.instances {
		gm.DelObject(inst)
	}
	bhvr.instances = make(map[gm.Object]gm.Object)
	gm.MustGetObjectParent(bhvr).Init()
}

func (bhvr *Room) Resume() {
	for _, inst := range bhvr.instances {
		inst.SetIsUpdate(true)
	}
}

func (bhvr *Room) Pause() {
	for _, inst := range bhvr.instances {
		inst.SetIsUpdate(false)
	}
}

// === Package functions ===

func IsOutside(instance gm.Object) bool {
	instCommon := gm.MustGetBehaviour(instance, gm.TypeID(&Room{}), gm.TypeID(&bhvrCommon.Common{})).(*bhvrCommon.Common)
	roomWindow := Data.instancesData[gm.ID(instance)].room.window()
	maskOuterRect := instCommon.TrasnformedMask().OuterRectangle()
	return !roomWindow.ContainsPoint(maskOuterRect.Lo()) && !roomWindow.ContainsPoint(maskOuterRect.Hi())
}
