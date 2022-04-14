package gogetautil

import (
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func MustNewEbitenImageFromFile(path string) *ebiten.Image {
	img, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		panic(err)
	}
	return img
}

// func MustLoadEbitenImagesFromDirectory(path string) map[string]*ebiten.Image {
// 	img, _, err := ebitenutil.NewImageFromFile(path)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return img
// }
