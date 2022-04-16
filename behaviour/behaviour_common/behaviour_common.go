package behaviour_common

import (
	"errors"
	"fmt"
	"log"
	"reflect"

	"github.com/golang/geo/r1"
	"github.com/golang/geo/r2"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/pikomonde/gogeta/gm"
)

type Common struct {
	Sprite       Sprite
	Position     r2.Point // Position of the object based on cartesian room
	Speed        r2.Point // Speed of the object based on cartesian room
	Angle        float64  // Angle of the object based on sprite anchor
	Zidx         float64  // Depth of the object
	IsStopUpdate bool     // Toggle object update or not
}

func (bhvr *Common) PreInit() {
	bhvr.Sprite.PreInit()
}

func (bhvr *Common) PostInit() {
	bhvr.Sprite.PostInit()
}

func (bhvr *Common) PreUpdate() {

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
	bhvr.Position.X += bhvr.Speed.X
	bhvr.Position.Y += bhvr.Speed.Y
}

func (bhvr *Common) PostUpdate() {
	objectDB := gm.GetObjectDB()
	obj, err := gm.GetObjectParent(bhvr)
	if err == nil {
		objectDB.MustSetObjectData(obj, gm.ObjectData{
			ZIdx:         bhvr.Zidx,
			IsStopUpdate: bhvr.IsStopUpdate,
		})
	}
}

func (bhvr *Common) Draw(screen *ebiten.Image) {
	// 	bhvrCommon := gm.MustGetBehaviourRel(bhvr, &Common{}).(*Common)
	frame := bhvr.Sprite.GetCurrentFrame()
	if frame == nil {
		return
	}
	// ebitenutil.DebugPrintAt(screen, "frame", 8, 8)
	// ebitenutil.DebugPrintAt(screen, fmt.Sprintf("animation: %s", bhvr.Sprite.CurrentAnimation), 8, 16)
	// ebitenutil.DebugPrintAt(screen, fmt.Sprintf("frame: %d", bhvr.Sprite.CurrentFrame), 8, 24)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-frame.Anchor.X, -frame.Anchor.Y)
	op.GeoM.Rotate(bhvr.Angle)
	op.GeoM.Translate(bhvr.Position.X, bhvr.Position.Y)
	// op.GeoM.Translate(60, 60)
	(*ebiten.Image)(screen).DrawImage(frame.Image, op)
}

// === Behaviour specific method ===

func (bhvr *Common) IsInside(p r2.Point) bool {
	// TODO: use mask in "Frame" instead (for collision)
	frame := bhvr.Sprite.GetCurrentFrame()
	w, h := frame.Image.Size()
	maskRectZero := bhvr.Position.Sub(frame.Anchor)
	maskRect := r2.Rect{
		X: r1.Interval{Lo: maskRectZero.X, Hi: maskRectZero.X + float64(w)},
		Y: r1.Interval{Lo: maskRectZero.Y, Hi: maskRectZero.Y + float64(h)},
	}
	return maskRect.ContainsPoint(p)
}

// === Package functions ===

// GetInstanceByObject get the first instance of an object.
func GetInstanceByObject(obj gm.Object) (gm.Object, error) {
	objDB := gm.GetObjectDB()
	objType := reflect.TypeOf(obj).String()

	key := fmt.Sprintf("%s%s", gm.KeyByObjType, objType)
	if _, ok := objDB[key]; !ok {
		return nil, errors.New(gm.ErrInstanceNotFound)
	}

	for inst, _ := range objDB[key] {
		return inst, nil
	}
	return nil, errors.New(gm.ErrInstanceNotFound)
}

// MustGetInstanceByObject get the first instance of an object. Panic if ust found.
func MustGetInstanceByObject(obj gm.Object) gm.Object {
	objDB := gm.GetObjectDB()
	objType := reflect.TypeOf(obj).String()

	key := fmt.Sprintf("%s%s", gm.KeyByObjType, objType)
	if _, ok := objDB[key]; !ok {
		log.Fatalf("[MustGetInstanceByObject] There is no instance for object %T.", obj)
		return nil
	}

	for inst, _ := range objDB[key] {
		return inst
	}

	log.Fatalf("[MustGetInstanceByObject] There is no instance for object %T.", obj)
	return nil
}
