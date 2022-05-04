package gm

type LayoutType int

const (
	// If you develop game for mobile it's better to use this default size (if portrait)
	DefaultCanvasW int = 372 // safe zone area is 270 pixel (based on 20:9), which is pixel 51-321
	DefaultCanvasH int = 600 // safe zone area is 496 pixel (based on 4:3), which is pixel 52-548
	// DefaultCanvasW = 93  // safe zone area is 67.5 pixel (based on 20:9), which is pixel 13-80
	// DefaultCanvasH = 150 // safe zone area is 124 pixel (based on 4:3), which is pixel 13-137
)

const (
	// LayoutType_Canvas means that the size of game screen will follow the game's canvas. See the
	// implementation in game's Layout function.
	LayoutType_SnapCanvas LayoutType = iota

	// LayoutType_Outside means that the size of game screen will follow the game's device/window. See
	// the implementation in game's Layout function.
	LayoutType_SnapOutside

	// LayoutType_Custom means that the size of game screen will follow user's function.
	LayoutType_Custom
)

type layout struct {
	canvasW, canvasH int
	screenW, screenH int
	layoutType       LayoutType
	layoutCustomFunc func(outsideWidth, outsideHeight int) (screenWidth, screenHeight int)
}

func SetLayoutType(layoutType LayoutType) {
	gm.layout.layoutType = layoutType
}

// SetCustomLayoutFunction set user's custom Layout function. Setting this not automatically use this
// function. You should set the layout type to LayoutType_Custom to use this.
func SetCustomLayoutFunction(fn func(outsideWidth, outsideHeight int) (screenWidth, screenHeight int)) {
	gm.layout.layoutCustomFunc = fn
}

func SetCanvasSize(w, h int) {
	gm.layout.canvasW, gm.layout.canvasH = w, h
}

func GetCanvasSize() (w, h int) {
	return gm.layout.canvasW, gm.layout.canvasH
}

func GetScreenSize() (w, h int) {
	return gm.layout.screenW, gm.layout.screenH
}
