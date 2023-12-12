package internal

import (
	"github.com/andreyladmj/gotils/glyph/internal/graphics"
	"github.com/andreyladmj/gotils/glyph/internal/internal"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"golang.org/x/image/font"
	"image/color"
)

const (
	WindowWidth  = 64 * 28
	WindowHeight = 64 * 16
)

type workspace struct {
	fonts map[string]*font.Face
	gl    *graphics.Glyph
	wi    *internal.WInput
}

func NewWorkspace(fonts map[string]*font.Face) *workspace {
	g := &workspace{
		fonts: fonts,
		wi:    &internal.WInput{},
		gl:    graphics.NewGlyph(),
	}

	clr := color.RGBA{
		R: 255,
		G: 0,
		B: 0,
		A: 200,
	}
	clr2 := clr
	clr2.G = 255

	//b1 := graphics.NewBox(internal.NewPos(30, 30), 500, 500, clr)
	//b2 := graphics.NewBox(internal.NewPos(300, 30), 100, 100, clr2)
	//b1.AddChildren(b2)
	//b1.AddChildren(graphics.NewGradient(internal.NewPos(300, 300), 500, 500, clr))
	//t1 := graphics.NewText(internal.NewPos(100, 400), g.fonts["normal"], "Hello, World!")
	//c1 := graphics.NewConnection(b1, b2)
	//c2 := graphics.NewConnection(b2, t1)
	//g.gl.AddChildren(b1, b2, c1, c2)
	//g.gl.AddChildren(b1)
	//g.gl.AddChildren(graphics.NewClickableGlyph(internal.NewPos(900, 400), 100, 100))
	//g.gl.AddChildren(graphics.NewCommand(internal.NewPos(900, 300), 100, 100, g.fonts["normal"], "Hello, \n World 123!", clr))
	g.gl.AddChildren(graphics.NewWords(internal.NewPos(0, 0), 800, 1000, "D:\\words.txt", g.fonts["normal"]))
	g.gl.AddChildren(graphics.NewWords(internal.NewPos(800, 0), 800, 1000, "D:\\esp.txt", g.fonts["normal"]))

	return g
}

func (g *workspace) Update() error {
	g.wi.UpdatePosition(ebiten.CursorPosition())
	g.wi.IsJustPressed = false
	g.wi.IsJustReleased = false

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if !g.wi.IsPressed {
			g.wi.UpdatePressedPosition(ebiten.CursorPosition())
			g.wi.IsJustPressed = true
		}

		g.wi.IsPressed = true
	}

	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		if g.wi.IsPressed {
			g.wi.IsJustReleased = true
		}
		g.wi.IsPressed = false
	}

	if err := g.gl.Update(g.wi); err != nil {
		return err
	}

	return nil
}

//// MouseStrokeSource is a StrokeSource implementation of mouse.
//type MouseStrokeSource struct{}
//
//func (m *MouseStrokeSource) Position() (int, int) {
//	return ebiten.CursorPosition()
//}
//
//func (m *MouseStrokeSource) IsJustReleased() bool {
//	return inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft)
//}
//
//// TouchStrokeSource is a StrokeSource implementation of touch.
//type TouchStrokeSource struct {
//	ID ebiten.TouchID
//}
//
//func (t *TouchStrokeSource) Position() (int, int) {
//	return ebiten.TouchPosition(t.ID)
//}
//
//func (t *TouchStrokeSource) IsJustReleased() bool {
//	return inpututil.IsTouchJustReleased(t.ID)
//}

func (g *workspace) Draw(screen *ebiten.Image) {
	// Локальные сокращения, чтобы уменьшить код по ширине
	// (Формат статьи накладывает свои ограничения)
	const w = WindowWidth
	const h = WindowHeight

	ebitenutil.DrawRect(screen, 0, 0, w, h, color.RGBA{
		R: 17,
		G: 38,
		B: 46,
		A: 255,
	})

	gridColor64 := &color.RGBA{A: 50}
	gridColor32 := &color.RGBA{A: 20}
	for y := 0.0; y < h; y += 32 {
		ebitenutil.DrawLine(screen, 0, y, w, y, gridColor32)
	}
	for y := 0.0; y < h; y += 64 {
		ebitenutil.DrawLine(screen, 0, y, w, y, gridColor64)
	}
	for x := 0.0; x < w; x += 32 {
		ebitenutil.DrawLine(screen, x, 0, x, h, gridColor32)
	}
	for x := 0.0; x < w; x += 64 {
		ebitenutil.DrawLine(screen, x, 0, x, h, gridColor64)
	}

	//gradientImg := ebiten.NewImage(500, 500)
	//gradientImg.Fill(color.White)
	//img1 := ebiten.NewImage(500, 500)
	//img1.Fill(color.White)
	//op := &ebiten.DrawTrianglesOptions{}
	//op.CompositeMode = ebiten.CompositeModeSourceOver
	////drawGradientRect(gradientImg, op, 50, 50, 500, 500, color.RGBA{255, 0, 0, 255}, color.RGBA{0, 0, 255, 255})
	//x := 50
	//y := 50
	//width := 500
	//height := 500
	//startColor := color.RGBA{255, 0, 0, 255}
	//endColor := color.RGBA{0, 0, 255, 255}
	//vertices := []ebiten.Vertex{
	//	{
	//		DstX:   float32(x),
	//		DstY:   float32(y),
	//		SrcX:   0,
	//		SrcY:   0,
	//		ColorR: float32(startColor.R) / 255,
	//		ColorG: float32(startColor.G) / 255,
	//		ColorB: float32(startColor.B) / 255,
	//		ColorA: float32(startColor.A) / 255,
	//	},
	//	{
	//		DstX:   float32(x + width),
	//		DstY:   float32(y),
	//		SrcX:   float32(width),
	//		SrcY:   0,
	//		ColorR: float32(endColor.R) / 255,
	//		ColorG: float32(endColor.G) / 255,
	//		ColorB: float32(endColor.B) / 255,
	//		ColorA: float32(endColor.A) / 255,
	//	},
	//	{
	//		DstX:   float32(x),
	//		DstY:   float32(y + height),
	//		SrcX:   0,
	//		SrcY:   float32(height),
	//		ColorR: float32(startColor.R) / 255,
	//		ColorG: float32(startColor.G) / 255,
	//		ColorB: float32(startColor.B) / 255,
	//		ColorA: float32(startColor.A) / 255,
	//	},
	//	{
	//		DstX:   float32(x + width),
	//		DstY:   float32(y + height),
	//		SrcX:   float32(width),
	//		SrcY:   float32(height),
	//		ColorR: float32(endColor.R) / 255,
	//		ColorG: float32(endColor.G) / 255,
	//		ColorB: float32(endColor.B) / 255,
	//		ColorA: float32(endColor.A) / 255,
	//	},
	//}
	//
	//// Create indices for the vertices
	//indices := []uint16{0, 1, 2, 1, 2, 3}
	//
	//// Draw the gradient rectangle using the vertices and indices
	////gradientImg.DrawTriangles(vertices, indices, img1, op)
	////screen.DrawImage(gradientImg, nil)
	//screen.DrawTriangles(vertices, indices, gradientImg, op)
	//screen.DrawImage(gradientImg, nil)

	g.gl.Render(screen)
}

func (g *workspace) Layout(outsideWidth, outsideHeight int) (int, int) {
	//fmt.Println("outsideWidth: ", outsideWidth, "outsideHeight: ", outsideHeight)
	return WindowWidth, WindowHeight
}

func drawGradientRect(gradientImg *ebiten.Image, op *ebiten.DrawTrianglesOptions, x, y, width, height int, startColor, endColor color.RGBA) {
	// Create a new image with the gradient

	// Create vertices for the gradient rectangle
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

	// Draw the gradient rectangle using the vertices and indices
	gradientImg.DrawTriangles(vertices, indices, nil, op)

}
