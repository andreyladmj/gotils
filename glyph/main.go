package main

import (
	"github.com/andreyladmj/gotils/glyph/internal"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"log"
)

////go:embed all:_assets
//var gameAssets embed.FS

const sampleText = `The quick brown fox jumps over the lazy dog.`

var (
	mplusNormalFont font.Face
	mplusBigFont    font.Face
)

func init() {
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 72
	mplusNormalFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    24,
		DPI:     dpi,
		Hinting: font.HintingVertical,
	})
	if err != nil {
		log.Fatal(err)
	}
	mplusBigFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    48,
		DPI:     dpi,
		Hinting: font.HintingFull, // Use quantization to save glyph cache images.
	})
	if err != nil {
		log.Fatal(err)
	}

	// Adjust the line height.
	mplusBigFont = text.FaceWithLineHeight(mplusBigFont, 54)
}

func main() {
	game := internal.NewWorkspace(map[string]*font.Face{
		"normal": &mplusNormalFont,
		"big":    &mplusBigFont,
	})
	ebiten.SetWindowTitle("Glyph")
	ebiten.SetWindowSize(internal.WindowWidth, internal.WindowHeight)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
