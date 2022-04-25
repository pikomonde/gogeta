package behaviour_room

import "github.com/pikomonde/gogeta/gm"

var (
	instancesData map[gm.Object]*InstanceData
)

func init() {
	instancesData = make(map[gm.Object]*InstanceData)
}

// === instances data ===

type InstanceData struct {
	room                *Room
	IsDeleteWhenOutside bool
}

func (data *InstanceData) Room() *Room { return data.room }

func Data(instance gm.Object) *InstanceData {
	return instancesData[instance]
}
