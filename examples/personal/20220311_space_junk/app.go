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
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	bhvrCommon "github.com/pikomonde/gogeta/behaviour/behaviour_common"
	bhvrRoom "github.com/pikomonde/gogeta/behaviour/behaviour_room"
	"github.com/pikomonde/gogeta/gm"
	"github.com/pikomonde/gogeta/gogetautil"
)

const (
	WindowWidth  = 600
	WindowHeight = 600
	CanvasWidth  = WindowWidth
	CanvasHeight = WindowHeight
	// CanvasWidth  = 186
	// CanvasHeight = 300
	// CanvasWidth  = 93  // safe zone area is 67.5 pixel (based on 20:9), which is pixel 13-80
	// CanvasHeight = 150 // safe zone area is 124 pixel (based on 4:3), which is pixel 13-137
)

// TODO:
// ui should get from actual width & height
// room's action should on canvas's safe zone

func main() {
	// Initialize objects
	gm.Init(WindowWidth, WindowHeight, CanvasWidth, CanvasHeight)
	gm.InitObject(&room01{})
	// img := ebiten.NewImage(32, 8)
	// img.Fill(color.Black)
	// (*objJunkAnimation["junk_energy"])[0].SetImage(img).SetAnchorToggle(bhvrCommon.Sprite_FrameAnchor_ToggleMiddleLeft).SetMaskFill()

	// Run game
	if err := gm.Run(); err != nil {
		log.Fatal("error game run: ", err)
	}
}

type room01 struct {
	BhvrRoom  bhvrRoom.Room
	Tick      uint64
	SpawnRate uint64
	JunkSpeed float64
	Bread     int64
	Energy    int64
	IsGameEnd bool
}

func (obj *room01) Init() {
	rand.Seed(time.Now().UnixNano())
	obj.Tick = 0
	obj.SpawnRate = 60
	obj.JunkSpeed = 1.4
	obj.Bread = 0
	obj.Energy = 12
}

func (obj *room01) Update() {
	// ebiten.SetWindowTitle(fmt.Sprintf("%2f %2f", ebiten.CurrentTPS(), ebiten.CurrentFPS()))

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
		if junkJunkTypeRand <= 20 {
			junkJunkType = junkTypeBread
		} else {
			junkJunkType = junkTypeEnergy
		}
		junkDistance := float64(CanvasWidth) * 0.8 // (CanvasWidth / 2) * 2^(0.5)
		junkAngle := 2 * math.Pi * rand.Float64()
		junkDirAngle := 0.25 * math.Pi * (rand.Float64() - 0.5)

		// create junk instance
		obj.BhvrRoom.InitObject(&objJunk{BhvrCommon: bhvrCommon.Common{
			Position: r2.Point{
				X: float64(CanvasWidth)/2 + junkDistance*math.Sin(junkAngle),
				Y: float64(CanvasHeight)/2 + junkDistance*math.Cos(junkAngle),
			},
			Speed: r2.Point{
				X: -math.Sin(junkAngle + junkDirAngle),
				Y: -math.Cos(junkAngle + junkDirAngle),
			}.Normalize().Mul(obj.JunkSpeed),
			Scale: r2.Point{X: 2, Y: 2},
			Zidx:  50,
			// IsDrawMask: true,
		},
			JunkType: junkJunkType,
		})
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

func (obj *room01) Draw(screen *ebiten.Image) {
	// Color the background
	screen.Fill(color.NRGBA{0x25, 0x20, 0x20, 0xff})

	// // TODO: add to util
	// objs := gm.GetObjectDB()
	// strArr := make([]string, 0)
	// for k, _ := range objs {
	// 	strArr = append(strArr, k)
	// }
	// sort.Strings(strArr)
	// lenSum := 0
	// for _, v := range strArr {
	// 	strArr2 := make([]string, 0)
	// 	strArr2Map := make(map[string]gm.Object)
	// 	for k, _ := range objs[v] {
	// 		strArr2 = append(strArr2, fmt.Sprintf("%p", k))
	// 		strArr2Map[fmt.Sprintf("%p", k)] = k
	// 	}
	// 	sort.Strings(strArr2)
	// 	for _, v2 := range strArr2 {
	// 		lenSum++
	// 		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("objects[%s]: %d %+v", v, len(objs[v]), objs[v][strArr2Map[v2]]), 8, 40+(lenSum*12))
	// 	}
	// }

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Energy: %d", obj.Energy), 8, 8)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Bread: %d", obj.Bread), 8, 24)

	if obj.IsGameEnd {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("THE END: YOU GOT %d BREAD", obj.Bread), CanvasWidth/2-75, CanvasHeight/2-4)
		ebitenutil.DebugPrintAt(screen, "CLICK ANYWHERE TO CONTINUE", CanvasWidth/2-78, CanvasHeight/2+12)
	}

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

type objJunk struct {
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

	// TODO: destroy out of room

	// slow down
	if obj.Tick > 200 {
		// gm.DelObject(obj)
		if obj.BhvrCommon.Speed.Norm() > 0.1 {
			obj.BhvrCommon.Speed = obj.BhvrCommon.Speed.Mul(0.97)
		}
	}

	// when object is clicked
	touchIDs := inpututil.AppendJustPressedTouchIDs([]ebiten.TouchID{})
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || (len(touchIDs) > 0) {
		mouseX, mouseY := ebiten.CursorPosition()
		if len(touchIDs) > 0 {
			mouseX, mouseY = ebiten.TouchPosition(touchIDs[0])
		}
		if obj.BhvrCommon.IsInside(r2.Point{float64(mouseX), float64(mouseY)}) {

			// increment collected junk
			objRoom01 := bhvrCommon.MustGetInstanceByObject(&room01{}).(*room01)
			if obj.JunkType == junkTypeBread {
				objRoom01.Bread++
			} else if obj.JunkType == junkTypeEnergy {
				objRoom01.Energy += 2
			}

			// increase difficulty
			objRoom01.JunkSpeed += 0.02

			// delete junk
			gm.DelObject(obj)
		}
	}
}

func (obj *objJunk) Draw(screen *ebiten.Image) {
	// Color the background
	// screen.Fill(color.NRGBA{0x00, 0x40, 0x80, 0xff})
}
