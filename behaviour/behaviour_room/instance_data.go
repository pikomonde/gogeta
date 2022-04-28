package behaviour_room

import (
	bhvrCommon "github.com/pikomonde/gogeta/behaviour/behaviour_common"
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

func (data *BehaviourInstancesData) TypeString() string { return "Room" }

func (data *BehaviourInstancesData) ByInstance(instance gm.Object) *InstanceData {
	return data.instancesData[instance]
}

func (data *BehaviourInstancesData) DelInstance(instance gm.Object) {
	bhvrRoomInst := data.ByInstance(instance).Room()
	delete(bhvrRoomInst.instances, instance)
	delete(data.instancesData, instance)
}

func (data *BehaviourInstancesData) PreUpdate() {
	// TODO: should filter with room active?
	for instance, instanceData := range Data.instancesData {
		if instanceData.IsDeleteWhenOutside {
			instCommon := gm.MustGetBehaviour(instance, Room{}.Type(), bhvrCommon.Common{}.Type()).(*bhvrCommon.Common)
			maskOuterRect := instCommon.TrasnformedMask().OuterRectangle()
			bhvrRoomInst := instanceData.Room()
			if !bhvrRoomInst.window().ContainsPoint(maskOuterRect.Lo()) && !bhvrRoomInst.window().ContainsPoint(maskOuterRect.Hi()) {
				gm.DelObject(instance)
			}
		}
	}
}

func (data *BehaviourInstancesData) PostUpdate() {

}

// === instance data ===

// TODO: I don't feel convenince with this instance data model, because instance data is the data of
// instance that not have Room's behaviour.
//
// The option that I think of is to put this
// data in Common as new field, or even in game as map[string]. I don't think the later will make
// the development easier, because there is no autocomplete that prevents typo, and we might need check
// isExist. Meanwhile the former will add bunch of other Behaviour instance data to Common that might not be used
//
// The other option that I can think is with new behaviour called RoomInstance, then we can let this
// Room behaviour interact wirh RoomInstance behaviour. But I think this will be a little bit tedious, since
// I think every instance can be child of the Room.
//
// If this case is true (every instance can be child of the Room), then it is a good idea
// to put this on the game package (as map[string]?), but again this is will prevent the ide for autocomplete.
// I think it's better to put in this instance data model for now.
type InstanceData struct {
	room                *Room
	IsDeleteWhenOutside bool
}

func (data *InstanceData) Room() *Room { return data.room }
