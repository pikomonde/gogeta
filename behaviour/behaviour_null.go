package behaviour

import "github.com/hajimehoshi/ebiten/v2"

type Null struct {
}

func (bhvr *Null) PreInit() {
}

func (bhvr *Null) PostInit() {
}

func (bhvr *Null) PreUpdate() {
}

func (bhvr *Null) PostUpdate() {
}

func (bhvr *Null) Draw(screen *ebiten.Image) {
}
