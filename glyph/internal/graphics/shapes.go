package graphics

import (
	"fmt"
	"github.com/andreyladmj/gotils/glyph/internal/internal"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
	"image/color"
)

// ---------------------------------------------

type Text struct {
	Glyph
	font *font.Face
	text string
}

func NewText(pos internal.Pos, font *font.Face, text string) *Text {
	return &Text{Glyph{pos: pos}, font, text}
}

func (g *Text) Render(screen *ebiten.Image) {
	text.Draw(screen, g.text, *g.font, int(g.pos.X), int(g.pos.Y), color.Black)
}

// ---------------------------------------------

type Box struct {
	Glyph
	backgroundColor color.RGBA
}

func NewBox(pos internal.Pos, w, h float32, backgroundColor color.RGBA) *Box {
	return &Box{Glyph: Glyph{pos: pos, w: w, h: h}, backgroundColor: backgroundColor}
}

func (g *Box) Render(screen *ebiten.Image) {
	clr := color.RGBA{
		R: 0,
		G: 0,
		B: 0,
		A: 255,
	}
	vector.DrawFilledRect(screen, g.pos.X, g.pos.Y, g.w, g.h, g.backgroundColor, false)

	vector.StrokeRect(screen, g.pos.X, g.pos.Y, g.w, g.h, 1, clr, false)

	for i := 0; i < len(g.children); i++ {
		g.children[i].Render(screen)
	}
}

// ---------------------------------------------

type Gradient struct {
	Glyph
	backgroundColor color.RGBA
}

func NewGradient(pos internal.Pos, w, h float32, backgroundColor color.RGBA) *Gradient {
	return &Gradient{Glyph: Glyph{pos: pos, w: w, h: h}, backgroundColor: backgroundColor}
}

func (g *Gradient) Render(screen *ebiten.Image) {
	gradientImg := ebiten.NewImage(int(g.w), int(g.h))
	gradientImg.Fill(color.White)
	op := &ebiten.DrawTrianglesOptions{}
	op.CompositeMode = ebiten.CompositeModeSourceOver
	x := g.pos.X
	y := g.pos.Y
	width := g.w
	height := g.h
	startColor := color.RGBA{255, 0, 0, 255}
	endColor := color.RGBA{0, 0, 255, 255}
	vertices := []ebiten.Vertex{
		{
			DstX:   float32(x),
			DstY:   float32(y),
			SrcX:   0,
			SrcY:   0,
			ColorR: float32(startColor.R) / 255,
			ColorG: float32(startColor.G) / 255,
			ColorB: float32(startColor.B) / 255,
			ColorA: float32(startColor.A) / 255,
		},
		{
			DstX:   float32(x + width),
			DstY:   float32(y),
			SrcX:   float32(width),
			SrcY:   0,
			ColorR: float32(endColor.R) / 255,
			ColorG: float32(endColor.G) / 255,
			ColorB: float32(endColor.B) / 255,
			ColorA: float32(endColor.A) / 255,
		},
		{
			DstX:   float32(x),
			DstY:   float32(y + height),
			SrcX:   0,
			SrcY:   float32(height),
			ColorR: float32(startColor.R) / 255,
			ColorG: float32(startColor.G) / 255,
			ColorB: float32(startColor.B) / 255,
			ColorA: float32(startColor.A) / 255,
		},
		{
			DstX:   float32(x + width),
			DstY:   float32(y + height),
			SrcX:   float32(width),
			SrcY:   float32(height),
			ColorR: float32(endColor.R) / 255,
			ColorG: float32(endColor.G) / 255,
			ColorB: float32(endColor.B) / 255,
			ColorA: float32(endColor.A) / 255,
		},
	}

	// Create indices for the vertices
	indices := []uint16{0, 1, 2, 1, 2, 3}

	screen.DrawTriangles(vertices, indices, gradientImg, op)

	for i := 0; i < len(g.children); i++ {
		g.children[i].Render(screen)
	}
}

// ---------------------------------------------

// ---------------------------------------------

type ClickableGlyph struct {
	Glyph
}

func NewClickableGlyph(pos internal.Pos, w, h float32) *ClickableGlyph {
	return &ClickableGlyph{Glyph{pos: pos, w: w, h: h}}
}

func (g *ClickableGlyph) Update(input *internal.WInput) error {
	if input.IsPressed && inside(&g.Glyph, input.PressedPos) {
		fmt.Println("Clicked!")
	}

	for i := 0; i < len(g.children); i++ {
		if err := g.children[i].Update(input); err != nil {
			return err
		}
	}

	return nil
}

func (g *ClickableGlyph) Render(screen *ebiten.Image) {
	vector.StrokeRect(screen, g.pos.X, g.pos.Y, g.w, g.h, 1, color.Black, false)

	for i := 0; i < len(g.children); i++ {
		g.children[i].Render(screen)
	}
}

// ---------------------------------------------

// ---------------------------------------------

type Command struct {
	Glyph
	font            *font.Face
	text            string
	backgroundColor color.RGBA
}

func NewCommand(pos internal.Pos, w, h float32, font *font.Face, text string, backgroundColor color.RGBA) *Command {
	return &Command{Glyph{pos: pos, w: w, h: h}, font, text, backgroundColor}
}
func (g *Command) Update(input *internal.WInput) error {
	bounds, advance := font.BoundString(*g.font, g.text)
	g.w = float32(advance.Round())
	g.h = float32(bounds.Min.Y.Round())

	for i := 0; i < len(g.children); i++ {
		if err := g.children[i].Update(input); err != nil {
			return err
		}
	}

	return nil
}

func (g *Command) Render(screen *ebiten.Image) {
	//vector.DrawEllipse(screen, 0, 0, 200, 200, g.backgroundColor)
	//vector.
	text.Draw(screen, g.text, *g.font, int(g.pos.X), int(g.pos.Y), color.Black)
	//advance := font.MeasureString(*g.font, g.text)
	//bounds, _ := font.BoundString(*g.font, g.text)
	//fmt.Println(advance.Round())
	//fmt.Println(advance)
	//fmt.Println(float32(bounds.Min.X.Round()), float32(bounds.Min.Y.Round()), float32(bounds.Max.X.Round()), bounds.Max.Y.Round())

	//vector.StrokeRect(screen, g.pos.X+float32(bounds.Min.X.Round()), g.pos.Y+float32(bounds.Min.Y.Round()), float32(bounds.Max.X.Round()), float32(bounds.Max.Y.Round()), 1, color.Black, false)

	vector.StrokeRect(screen, g.pos.X, g.pos.Y, g.w, g.h, 1, color.Black, false)
	//text.Draw(screen, g.text, *g.font, int(g.pos.X), int(g.pos.Y), color.Black)
	//text := NewText(pos internal.Pos, font *font.Face, text string)

	for i := 0; i < len(g.children); i++ {
		g.children[i].Render(screen)
	}
}

// ---------------------------------------------

// ---------------------------------------------

//type Ellipse struct {
//	BoxGlyph
//	backgroundColor color.RGBA
//}
//
//func NewEllipse(pos internal.Pos, w, h float32, backgroundColor color.RGBA) *Ellipse {
//	return &Ellipse{BoxGlyph{Glyph: Glyph{pos: pos}, w: w, h: h}, backgroundColor}
//}
//
//func (g *Ellipse) Render(screen *ebiten.Image) {
//	//vector.DrawEllipse(screen, 0, 0, 200, 200, g.backgroundColor)
//	//vector.
//
//	for i := 0; i < len(g.children); i++ {
//		g.children[i].Render(screen)
//	}
//}
