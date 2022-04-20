package txt

import (
	"image"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// === CharSet ===

type CharSet uint8

const (
	CharSet_0020_007F CharSet = iota
)

func (cs CharSet) Min() rune {
	switch cs {
	case CharSet_0020_007F:
		return 0x0020
	}
	return 0x0000
}

func (cs CharSet) Max() rune {
	switch cs {
	case CharSet_0020_007F:
		return 0x007F
	}
	return 0x0000
}

// === Allignment ===

type Allignment uint8

const (
	Allignment_TopLeft Allignment = iota
	Allignment_TopCenter
	Allignment_TopRight
	Allignment_MiddleLeft
	Allignment_MiddleCenter
	Allignment_MiddleRight
	Allignment_BottomLeft
	Allignment_BottomCenter
	Allignment_BottomRight
)

type Font struct {
	Size       uint64 // Size is based on the height of character
	LineHeight uint64 // Space between lines
	Allignment Allignment
	cw, ch     int
	chars      map[rune]*ebiten.Image
}

func (f Font) Draw(screen *ebiten.Image, str string, px, py int) {
	op := &ebiten.DrawImageOptions{}
	scale := float64(f.Size) / float64(f.ch)
	lines := strings.Split(str, "\n")
	yMax := float64(len(lines)*f.ch) * scale
	x := float64(0)
	y := float64(0)
	for _, line := range lines {
		xMax := float64(len(line)*f.cw) * scale
		for _, c := range line {
			chr, ok := f.chars[c]
			if !ok {
				chr, ok = f.chars[' ']
				if !ok {
					chr = ebiten.NewImage(1, 1)
				}
			}

			op.GeoM.Reset()
			op.GeoM.Scale(scale, scale)
			op.GeoM.Translate(x+float64(px), y+float64(py))
			switch f.Allignment {
			case Allignment_TopLeft:
				op.GeoM.Translate(0, 0)
			case Allignment_TopCenter:
				op.GeoM.Translate(-xMax/2, 0)
			case Allignment_TopRight:
				op.GeoM.Translate(-xMax, 0)
			case Allignment_MiddleLeft:
				op.GeoM.Translate(0, -yMax/2)
			case Allignment_MiddleCenter:
				op.GeoM.Translate(-xMax/2, -yMax/2)
			case Allignment_MiddleRight:
				op.GeoM.Translate(-xMax, -yMax/2)
			case Allignment_BottomLeft:
				op.GeoM.Translate(0, -yMax)
			case Allignment_BottomCenter:
				op.GeoM.Translate(-xMax/2, -yMax)
			case Allignment_BottomRight:
				op.GeoM.Translate(-xMax, -yMax)
			}
			screen.DrawImage(chr, op)
			x += float64(f.cw) * scale
		}
		x = 0
		y += (float64(f.ch) + float64(f.LineHeight)) * scale
	}
}

func MustNewFontFromFile(path string, cw, ch int, charSet CharSet, opt Font) Font {
	charMap, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		panic(err)
	}

	w, _ := charMap.Size()
	chars := make(map[rune]*ebiten.Image)
	for c := charSet.Min(); c <= charSet.Max(); c++ {
		n := w / cw
		cx := (int(c-charSet.Min()) % n) * cw
		cy := (int(c-charSet.Min()) / n) * ch
		chr := charMap.SubImage(image.Rect(cx, cy, cx+cw, cy+ch)).(*ebiten.Image)
		chars[c] = chr
	}

	return Font{chars: chars, cw: cw, ch: ch}
}
