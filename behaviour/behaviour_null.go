package behaviour

import "github.com/hajimehoshi/ebiten/v2"

type Null struct {
}

func (bhvr *Null) Init() {
}

func (bhvr *Null) Update() {
}

func (bhvr *Null) Draw(screen *ebiten.Image) {
}
