package main

import (
	"fmt"
	"image/color"
	_ "image/png"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/golang/geo/r2"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	bhvrCommon "github.com/pikomonde/gogeta/behaviour/behaviour_common"
	bhvrRoom "github.com/pikomonde/gogeta/behaviour/behaviour_room"
	"github.com/pikomonde/gogeta/gm"
	"github.com/pikomonde/gogeta/gogetautil"
	"github.com/pikomonde/gogeta/gogetautil/txt"
)

const (
	CanvasWidth  = 600
	CanvasHeight = 600
)

var fontBocil57 = txt.MustNewFontFromFile("asset/sprite/font_bocil_57_0020_007F.png", 5, 7, txt.CharSet_0020_007F, txt.Font{Size: 24})

func main() {
	// Initialize objects
	gm.SetLayoutType(gm.LayoutType_SnapOutside)
	gm.InitObject(&roomMain{})
	ebiten.SetWindowSize(CanvasWidth, CanvasHeight)
	// ebiten.SetWindowResizable(true)

	// Run game
	if err := gm.Run(); err != nil {
		log.Fatal("error game run: ", err)
	}
}

type ui struct {
	gm.Objecter
}

func (obj *ui) Init()   {}
func (obj *ui) Update() {}
func (obj *ui) Draw(screen *ebiten.Image) {
	scrW, scrH := gm.GetScreenSize()
	fontSizeMul := float64(scrW) / float64(CanvasHeight)
	instRoomMain := gm.MustGetObjectParent(bhvrRoom.Data.ByInstance(gm.ID(obj)).Room()).(*roomMain)

	fontBocil57.LineHeight = 3
	fontBocil57.Size = uint64(24 * fontSizeMul)
	fontBocil57.Allignment = txt.Allignment_TopLeft
	fontBocil57.Draw(screen, fmt.Sprintf("ENERGY: %d\nBREAD: %d", instRoomMain.Energy, instRoomMain.Bread), 16, 16)

	if instRoomMain.IsGameEnd {
		fontBocil57.Size = uint64(20 * fontSizeMul)
		fontBocil57.Allignment = txt.Allignment_MiddleCenter
		fontBocil57.Draw(screen,
			fmt.Sprintf("THE END: YOU GOT %d BREAD\nCLICK ANYWHERE TO CONTINUE", instRoomMain.Bread),
			scrW/2, scrH/2)
	}

	fontBocil57.Size = uint64(12 * fontSizeMul)
	fontBocil57.Allignment = txt.Allignment_BottomLeft
	fontBocil57.Draw(screen,
		fmt.Sprintf("%2f %2f %d", ebiten.CurrentTPS(), ebiten.CurrentFPS(), len(gm.GetInstIDs())),
		16, scrH-16)
}

type roomMain struct {
	gm.Objecter
	BhvrRoom  bhvrRoom.Room
	Tick      uint64
	SpawnRate uint64
	JunkSpeed float64
	Bread     int64
	Energy    int64
	IsGameEnd bool
	TouchIDs  []ebiten.TouchID
}

func (obj *roomMain) Init() {
	rand.Seed(time.Now().UnixNano())
	obj.BhvrRoom.Size = r2.Point{X: float64(CanvasWidth), Y: float64(CanvasHeight)}
	obj.Tick = 0
	obj.SpawnRate = 60
	obj.JunkSpeed = 1.4
	obj.Bread = 0
	obj.Energy = 12

	// Zidx: 95
	obj.BhvrRoom.InitObject(&ui{}, bhvrRoom.InstanceData{}).SetZidx(95)
}

func (obj *roomMain) Update() {
	scrW, scrH := gm.GetScreenSize()
	obj.BhvrRoom.Size = r2.Point{X: float64(scrW), Y: float64(scrH)}

	// move this to other object called control/controller
	obj.TouchIDs = inpututil.AppendJustPressedTouchIDs([]ebiten.TouchID{})

	// is ended
	if obj.IsGameEnd {
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || (len(inpututil.AppendJustPressedTouchIDs([]ebiten.TouchID{})) > 0) {
			obj.IsGameEnd = false
			obj.BhvrRoom.Restart()
		}
		return
	}

	// every (obj.SpawnRate / TPS) seconds
	if obj.Tick%obj.SpawnRate == 0 {

		// initialize junk instance value
		junkJunkTypeRand := rand.Intn(100)
		var junkJunkType junkType
		if junkJunkTypeRand <= 35 {
			junkJunkType = junkTypeBread
		} else {
			junkJunkType = junkTypeEnergy
		}
		junkDistanceWidth := float64(scrW) * 0.8  // (scrW / 2) * 2^(0.5)
		junkDistanceHeight := float64(scrH) * 0.8 // (scrH / 2) * 2^(0.5)
		junkAngle := 2 * math.Pi * rand.Float64()
		junkDirAngle := 0.25 * math.Pi * (rand.Float64() - 0.5)

		// if ebiten.IsKeyPressed(ebiten.KeyRight) {
		// 	for i := 0; i < 500; i++ {
		// 		obj.BhvrRoom.InitObject(&objJunk{}, bhvrRoom.InstanceData{})
		// 	}
		// }
		// if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		// 	i := 0
		// 	for _, instID := range gm.GetInstIDsByObjTypeID(gm.TypeID(&instJunk)) {
		// 		gm.DelObject(gm.GetInstByObjInstID(instID))
		// 		i++
		// 		if i >= 500 {
		// 			break
		// 		}
		// 	}
		// }

		// create junk instance
		obj.BhvrRoom.InitObject(
			&objJunk{
				BhvrCommon: bhvrCommon.Common{
					Position: r2.Point{
						X: junkDistanceWidth * math.Sin(junkAngle),
						Y: junkDistanceHeight * math.Cos(junkAngle),
					}.Add(r2.Point{float64(scrW) / 2, float64(scrH) / 2}),
					Speed: r2.Point{
						X: -math.Sin(junkAngle + junkDirAngle),
						Y: -math.Cos(junkAngle + junkDirAngle),
					}.Normalize().Mul(obj.JunkSpeed),
					Scale: r2.Point{X: 2, Y: 2},
					// IsDrawMask: true,
				},
				JunkType: junkJunkType,
			},
			bhvrRoom.InstanceData{},
		).SetZidx(50)
		obj.Energy--
	}

	// end the game
	if obj.Energy <= 0 {
		// pause the game
		obj.BhvrRoom.Pause()
		obj.IsGameEnd = true
	}

	obj.Tick++
}

func (obj *roomMain) Draw(screen *ebiten.Image) {
	// Color the background
	screen.Fill(color.NRGBA{0x25, 0x20, 0x20, 0xff})

	// // TODO: add to util
	// objsObj := gm.GetBehavioursByObjInst()
	// objs := make(map[string]map[int]gm.Behaviour)
	// for k, v := range objsObj {
	// 	objs[fmt.Sprintf("%p", k)] = v
	// }
	// strArr := make([]string, 0)
	// for k, _ := range objs {
	// 	strArr = append(strArr, k)
	// }
	// sort.Strings(strArr)
	// lenSum := 0
	// for _, v := range strArr {
	// 	strArr2 := make([]int, 0)
	// 	for k, _ := range objs[v] {
	// 		strArr2 = append(strArr2, k)
	// 	}
	// 	sort.Ints(strArr2)
	// 	for _, v2 := range strArr2 {
	// 		lenSum++
	// 		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("objects[%s]: %d %+v", v, len(objs[v]), objs[v][v2]), 8, 40+(lenSum*12))
	// 	}
	// }

}

/*
  objJunk
*/

var objJunkAnimation bhvrCommon.Animations = bhvrCommon.Animations{
	"junk_ball": &bhvrCommon.Animation{
		(&bhvrCommon.Frame{}).SetImage(gogetautil.MustNewEbitenImageFromFile("asset/sprite/spr_red_ball_16x16.png")).SetAnchorToggle(bhvrCommon.Sprite_FrameAnchor_ToggleMiddleCenter).SetMaskFill(),
	},
	"junk_energy": &bhvrCommon.Animation{
		(&bhvrCommon.Frame{}).SetImage(gogetautil.MustNewEbitenImageFromFile("asset/sprite/spr_battery_16x16.png")).SetAnchorToggle(bhvrCommon.Sprite_FrameAnchor_ToggleMiddleCenter).SetMaskFill(),
	},
	"junk_bread": &bhvrCommon.Animation{
		(&bhvrCommon.Frame{}).SetImage(gogetautil.MustNewEbitenImageFromFile("asset/sprite/spr_bread_slice_16x16.png")).SetAnchorToggle(bhvrCommon.Sprite_FrameAnchor_ToggleMiddleCenter).SetMaskFill(),
		(&bhvrCommon.Frame{}).SetImage(gogetautil.MustNewEbitenImageFromFile("asset/sprite/spr_croisant_16x16.png")).SetAnchorToggle(bhvrCommon.Sprite_FrameAnchor_ToggleMiddleCenter).SetMaskFill(),
	},
	"ui_energy": &bhvrCommon.Animation{
		(&bhvrCommon.Frame{}).SetImage(gogetautil.MustNewEbitenImageFromFile("asset/sprite/spr_electricity_16x16.png")).SetAnchorToggle(bhvrCommon.Sprite_FrameAnchor_ToggleMiddleCenter).SetMaskFill(),
	},
	"ui_bread": &bhvrCommon.Animation{
		(&bhvrCommon.Frame{}).SetImage(gogetautil.MustNewEbitenImageFromFile("asset/sprite/spr_bread_slice_16x16.png")).SetAnchorToggle(bhvrCommon.Sprite_FrameAnchor_ToggleMiddleCenter).SetMaskFill(),
	},
}

type junkType uint8

const (
	junkTypeEnergy junkType = iota + 1
	junkTypeBread
)

var instJunk objJunk

type objJunk struct {
	gm.Objecter
	BhvrCommon bhvrCommon.Common
	Tick       uint64
	JunkType   junkType
}

func (obj *objJunk) Init() {
	obj.BhvrCommon.Sprite.Animations = objJunkAnimation
	switch obj.JunkType {
	case junkTypeEnergy:
		obj.BhvrCommon.Sprite.CurrentAnimation = "junk_energy"
	case junkTypeBread:
		obj.BhvrCommon.Sprite.CurrentAnimation = "junk_bread"
	default:
		obj.BhvrCommon.Sprite.CurrentAnimation = "junk_ball"
	}
	obj.BhvrCommon.Sprite.CurrentFrame = rand.Intn(len(*objJunkAnimation[obj.BhvrCommon.Sprite.CurrentAnimation]))

	// frame := &bhvrCommon.Frame{Image: ebiten.NewImage(int(32), int(32))}
	// frame.Image.Fill(color.NRGBA{0x00, 0x80, 0x00, 0xff})
	// frame.SetAnchorToggle(bhvrCommon.Sprite_FrameAnchor_ToggleMiddleCenter)
	// obj.BhvrCommon.Sprite.InsertFrame("flat", frame)
	// obj.BhvrCommon.Sprite.CurrentAnimation = "flat"
}

func (obj *objJunk) Update() {
	obj.Tick++
	obj.BhvrCommon.Angle += 0.01

	// activate out-of-room destroy instance, once inside the room
	// TODO: this is expensive
	if !bhvrRoom.IsOutside(obj) {
		bhvrRoom.Data.ByInstance(gm.ID(obj)).IsDeleteWhenOutside = true
	}

	// slow down
	if obj.Tick > 200 {
		if obj.BhvrCommon.Speed.Norm() > 0.3 {
			obj.BhvrCommon.Speed = obj.BhvrCommon.Speed.Mul(0.97)
		}
	}

	// TODO: this function is actually expensive, because ther are 10000 instance calling this function,
	// and inside this function, it is querying in a map of 10000 instance
	// _ = bhvrRoom.Data.ByInstance(obj.ID())
	// when object is clicked
	touchIDs := bhvrRoom.Data.ByInstance(gm.ID(obj)).Parent().(*roomMain).TouchIDs
	// touchIDs := []ebiten.TouchID{}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) ||
		len(touchIDs) > 0 {
		mouseX, mouseY := ebiten.CursorPosition()
		if len(touchIDs) > 0 {
			mouseX, mouseY = ebiten.TouchPosition(touchIDs[0])
		}
		// TODO: This is expensive
		if obj.BhvrCommon.IsInside(r2.Point{X: float64(mouseX), Y: float64(mouseY)}) {
			// increment collected junk
			instRoomMain := gm.MustGetObjectParent(bhvrRoom.Data.ByInstance(gm.ID(obj)).Room()).(*roomMain)
			if obj.JunkType == junkTypeBread {
				instRoomMain.Bread++
			} else if obj.JunkType == junkTypeEnergy {
				instRoomMain.Energy += 2
			}

			// increase difficulty
			instRoomMain.JunkSpeed += 0.02

			// delete junk
			gm.DelObject(obj)
		}
	}
}

func (obj *objJunk) Draw(screen *ebiten.Image) {
}
