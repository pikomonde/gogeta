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

// === Behaviour data ===

var (
	instancesData map[gm.Object]*InstanceData
)

type InstanceData struct {
	room                *Room
	IsDeleteWhenOutside bool
}

func (data *InstanceData) Room() *Room { return data.room }

func init() {
	instancesData = make(map[gm.Object]*InstanceData)
}

// === Behaviour ===

type Room struct {
	Position  r2.Point
	Size      r2.Point
	instances []gm.Object
}

func (bhvr *Room) Instances() []gm.Object { return bhvr.instances }

// TODO: what to do with room's window? is it as view camera? do we still need game's canvas width/height?
// is there any limitation on "active" room (because we need to attach camera to a room)? should we attach
// a camera to a room? how camera works?
func (bhvr *Room) window() r2.Rect {
	return r2.Rect{
		X: r1.Interval{Lo: bhvr.Position.X, Hi: bhvr.Position.X + bhvr.Size.X},
		Y: r1.Interval{Lo: bhvr.Position.Y, Hi: bhvr.Position.Y + bhvr.Size.Y},
	}
}

func (bhvr *Room) PreInit() {
	bhvr.instances = make([]gm.Object, 0)
}

func (bhvr *Room) PostInit() {
	if (bhvr.Size.X <= 0) || (bhvr.Size.Y <= 0) {
		log.Panicf("Room size cannot be empty")
	}
}

func (bhvr *Room) PreUpdate() {

	// TODO: maybe move this to behaviour's package update (if we want to add behaviour's package update)
	for _, instance := range bhvr.instances {
		if instancesData[instance].IsDeleteWhenOutside {
			instCommon := gm.MustGetBehaviour(instance, bhvr, &bhvrCommon.Common{}).(*bhvrCommon.Common)
			maskOuterRect := instCommon.TrasnformedMask().OuterRectangle()
			if !bhvr.window().ContainsPoint(maskOuterRect.Lo()) && !bhvr.window().ContainsPoint(maskOuterRect.Hi()) {
				gm.DelObject(instance)
			}
		}
	}
}

func (bhvr *Room) PostUpdate() {
}

func (bhvr *Room) Draw(screen *ebiten.Image) {
}

// === Behaviour specific method ===

// InitObject register current instance to this room. It is also register to gm
// through gm.InitObject().
func (bhvr *Room) InitObject(instance gm.Object, data InstanceData) gm.Object {
	if _, ok := instancesData[instance]; ok {
		instance := gm.MustGetObjectParent(bhvr)
		instanceName := reflect.TypeOf(instance).String()
		log.Panicf("This instance is already assigned to room %s", instanceName)
	}
	initObj := gm.InitObject(instance)
	bhvr.instances = append(bhvr.instances, initObj)
	instancesData[instance] = &InstanceData{
		room:                bhvr,
		IsDeleteWhenOutside: data.IsDeleteWhenOutside,
	}
	return initObj
}

func (bhvr *Room) Restart() {
	for _, inst := range bhvr.instances {
		gm.DelObject(inst)
	}
	bhvr.instances = make([]gm.Object, 0)
	gm.MustGetObjectParent(bhvr).Init()
}

func (bhvr *Room) Resume() {
	for _, inst := range bhvr.instances {
		bCommon := gm.MustGetBehaviour(inst, bhvr, &bhvrCommon.Common{}).(*bhvrCommon.Common)
		bCommon.IsStopUpdate = false
	}
}

func (bhvr *Room) Pause() {
	for _, inst := range bhvr.instances {
		bCommon := gm.MustGetBehaviour(inst, bhvr, &bhvrCommon.Common{}).(*bhvrCommon.Common)
		bCommon.IsStopUpdate = true
	}
}

// === Package functions ===
func IsOutside(instance gm.Object) bool {
	instCommon := gm.MustGetBehaviour(instance, &Room{}, &bhvrCommon.Common{}).(*bhvrCommon.Common)
	roomWindow := instancesData[instance].room.window()
	maskOuterRect := instCommon.TrasnformedMask().OuterRectangle()
	return !roomWindow.ContainsPoint(maskOuterRect.Lo()) && !roomWindow.ContainsPoint(maskOuterRect.Hi())
}

func Data(instance gm.Object) *InstanceData {
	return instancesData[instance]
}
