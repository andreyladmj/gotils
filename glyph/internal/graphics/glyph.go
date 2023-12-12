package graphics

import (
	"github.com/andreyladmj/gotils/glyph/internal/internal"
	"github.com/hajimehoshi/ebiten/v2"
)

type Drawable interface {
	Render(screen *ebiten.Image)
	Update(input *internal.WInput) error
	SetParent(glyph Drawable)
}

type Glyph struct {
	pos      internal.Pos
	children []Drawable
	parent   Drawable
	w, h     float32
}

func NewGlyph() *Glyph {
	return &Glyph{}
}

func (g *Glyph) SetPosition(pos internal.Pos) error {
	g.pos = pos
	return nil
}

func (g *Glyph) GetPosition() internal.Pos {
	return g.pos
}

func (g *Glyph) Update(input *internal.WInput) error {

	for i := 0; i < len(g.children); i++ {
		if err := g.children[i].Update(input); err != nil {
			return err
		}
	}

	return nil
}

func (g *Glyph) Render(screen *ebiten.Image) {
	for i := 0; i < len(g.children); i++ {
		g.children[i].Render(screen)
	}
}

func (g *Glyph) AddChildren(glh ...Drawable) {
	for i := 0; i < len(glh); i++ {
		glh[i].SetParent(g)
		g.children = append(g.children, glh[i])
	}
}

func (g *Glyph) GetParent() Drawable {
	return g.parent
}

func (g *Glyph) SetParent(glyph Drawable) {
	g.parent = glyph
}

func (g *Glyph) GetChildren() []Drawable {
	return g.children
}
