package behaviour_common

import (
	"image/color"

	"github.com/golang/geo/r1"
	"github.com/golang/geo/r2"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/pikomonde/gogeta/gm"
)

type Common struct {
	Sprite       Sprite
	Position     r2.Point // Position of the instance based on cartesian room
	Speed        r2.Point // Speed of the instance based on cartesian room
	Angle        float64  // Angle of the instance based on sprite anchor
	Scale        r2.Point // Scale of the instance based on sprite anchor
	Zidx         float64  // Depth of the instance
	IsStopUpdate bool     // Toggle instance update or not
	IsDrawMask   bool     // Draw instance's mask
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
	op.GeoM.Translate(-frame.Anchor().X, -frame.Anchor().Y) // step A
	op.GeoM.Scale(bhvr.Scale.X, bhvr.Scale.Y)               // step B
	op.GeoM.Rotate(bhvr.Angle)                              // step C
	op.GeoM.Translate(bhvr.Position.X, bhvr.Position.Y)     // step D
	(*ebiten.Image)(screen).DrawImage(frame.Image(), op)

	if bhvr.IsDrawMask {
		bhvr.drawMask(screen)
	}
}

// === Behaviour specific method ===

func (bhvr *Common) IsInside(pPos r2.Point) bool {
	// TODO: use mask in "Frame" instead (for collision)
	frame := bhvr.Sprite.GetCurrentFrame()
	vectors := frame.Mask().Vectors()

	geoMForP := ebiten.GeoM{}
	geoMForP.Translate(-bhvr.Position.X, -bhvr.Position.Y) // reversed step D
	geoMForP.Rotate(-bhvr.Angle)                           // reversed step C

	geoMForMask := ebiten.GeoM{}
	geoMForMask.Translate(-frame.Anchor().X, -frame.Anchor().Y) // step A
	geoMForMask.Scale(bhvr.Scale.X, bhvr.Scale.Y)               // step B

	switch frame.MaskType() {
	case Sprite_MaskType_Circle:
		// TODO: implement Sprite_MaskType_Circle mask
		return false
	case Sprite_MaskType_Recatangle:
		pX, pY := geoMForP.Apply(pPos.X, pPos.Y)
		maskZeroX, maskZeroY := geoMForMask.Apply(vectors[0].X, vectors[0].Y)
		maskEndX, maskEndY := geoMForMask.Apply(vectors[2].X, vectors[2].Y)
		mask := r2.Rect{
			X: r1.Interval{Lo: maskZeroX, Hi: maskEndX},
			Y: r1.Interval{Lo: maskZeroY, Hi: maskEndY},
		}
		return mask.ContainsPoint(r2.Point{X: pX, Y: pY})
	case Sprite_MaskType_Capsule:
		// TODO: implement Sprite_MaskType_Capsule mask
		return false
	case Sprite_MaskType_ConvexHull:
		// TODO: implement Sprite_MaskType_ConvexHull mask
		return false
	default:
		return false
	}
}

func (bhvr *Common) TrasnformedMask() Mask {
	frame := bhvr.Sprite.GetCurrentFrame()

	geoM := ebiten.GeoM{}
	geoM.Translate(-frame.Anchor().X, -frame.Anchor().Y) // step A
	geoM.Scale(bhvr.Scale.X, bhvr.Scale.Y)               // step B
	geoM.Rotate(bhvr.Angle)                              // step C
	geoM.Translate(bhvr.Position.X, bhvr.Position.Y)     // step D

	return frame.Mask().GeoTransform(geoM)
}

// TODO: geoM.Apply using gpu
func (bhvr *Common) drawMask(screen *ebiten.Image) {
	var path vector.Path
	lineWidth := float64(1)
	frame := bhvr.Sprite.GetCurrentFrame()
	vectors := frame.Mask().Vectors()

	geoM := ebiten.GeoM{}
	geoM.Translate(-frame.Anchor().X, -frame.Anchor().Y) // step A
	geoM.Scale(bhvr.Scale.X, bhvr.Scale.Y)               // step B
	geoM.Rotate(bhvr.Angle)                              // step C
	geoM.Translate(bhvr.Position.X, bhvr.Position.Y)     // step D

	switch frame.MaskType() {
	case Sprite_MaskType_Circle:
		// TODO: implement Sprite_MaskType_Circle mask
	case Sprite_MaskType_Recatangle:
		innerVectors := make([]r2.Point, 0)
		innerVectors = append(innerVectors, r2.Point{X: vectors[0].X + lineWidth, Y: vectors[0].Y + lineWidth})
		innerVectors = append(innerVectors, r2.Point{X: vectors[1].X - lineWidth, Y: vectors[1].Y + lineWidth})
		innerVectors = append(innerVectors, r2.Point{X: vectors[2].X - lineWidth, Y: vectors[2].Y - lineWidth})
		innerVectors = append(innerVectors, r2.Point{X: vectors[3].X + lineWidth, Y: vectors[3].Y - lineWidth})

		for i, vector := range vectors {
			x, y := geoM.Apply(vector.X, vector.Y)
			if i == 0 {
				path.MoveTo(float32(x), float32(y))
				continue
			}
			path.LineTo(float32(x), float32(y))
		}

		for i, vector := range innerVectors {
			x, y := geoM.Apply(vector.X, vector.Y)
			if i == 0 {
				path.MoveTo(float32(x), float32(y))
				continue
			}
			path.LineTo(float32(x), float32(y))
		}
	case Sprite_MaskType_Capsule:
		// TODO: implement Sprite_MaskType_Capsule mask
	case Sprite_MaskType_ConvexHull:
		// TODO: implement Sprite_MaskType_ConvexHull mask
	default:
		return
	}

	vs, is := path.AppendVerticesAndIndicesForFilling(nil, nil)
	emptyImage := ebiten.NewImage(1, 1)
	emptyImage.Fill(color.NRGBA{0xff, 0x00, 0x00, 0xff})
	screen.DrawTriangles(vs, is, emptyImage, &ebiten.DrawTrianglesOptions{
		FillRule: ebiten.EvenOdd,
	})
}

// === Package functions ===

// // GetInstanceByObject get the first instance of an object.
// func GetInstanceByObject(obj gm.Object) (gm.Object, error) {
// 	objDB := gm.GetObjectDB()
// 	objType := reflect.TypeOf(obj).String()

// 	key := fmt.Sprintf("%s%s", gm.KeyByObjType, objType)
// 	if _, ok := objDB[key]; !ok {
// 		return nil, errors.New(gm.ErrInstanceNotFound)
// 	}

// 	for inst, _ := range objDB[key] {
// 		return inst, nil
// 	}
// 	return nil, errors.New(gm.ErrInstanceNotFound)
// }

// // MustGetInstanceByObject get the first instance of an object. Panic if ust found.
// func MustGetInstanceByObject(obj gm.Object) gm.Object {
// 	objDB := gm.GetObjectDB()
// 	objType := reflect.TypeOf(obj).String()

// 	key := fmt.Sprintf("%s%s", gm.KeyByObjType, objType)
// 	if _, ok := objDB[key]; !ok {
// 		log.Panicf("[MustGetInstanceByObject] There is no instance for object %T.", obj)
// 		return nil
// 	}

// 	for inst, _ := range objDB[key] {
// 		return inst
// 	}

// 	log.Panicf("[MustGetInstanceByObject] There is no instance for object %T.", obj)
// 	return nil
// }
