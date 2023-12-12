package graphics

import "github.com/andreyladmj/gotils/glyph/internal/internal"

type Draggable interface {
	IsDraggable() bool
	IsDragging() bool
}

type DraggableGlyph struct {
	Glyph
	initPos     internal.Pos
	isDraggable bool
	isDragging  bool
}

func NewDraggableGlyph() *DraggableGlyph {
	return &DraggableGlyph{}
}

func inside(g *Glyph, pos internal.Pos) bool {
	return pos.X >= g.pos.X && pos.X <= g.pos.X+g.w && pos.Y >= g.pos.Y && pos.Y <= g.pos.Y+g.h
}

func (g *DraggableGlyph) IsDraggable() bool {
	return g.isDraggable
}

func (g *DraggableGlyph) IsDragging() bool {
	return g.isDragging
}

func (g *DraggableGlyph) Update(input *internal.WInput) error {

	if !g.IsDragging() && g.isDraggable && inside(&g.Glyph, input.PressedPos) {
		g.isDragging = true
	}

	if g.IsDragging() {
		if input.IsPressed {
			diff := input.PressedPos.Subtract(input.Pos)
			g.pos = g.initPos.Add(diff)
		}

		if input.IsJustReleased {
			g.initPos = g.pos
			g.isDragging = false
		}
	}

	for i := 0; i < len(g.children); i++ {
		if err := g.children[i].Update(input); err != nil {
			return err
		}
	}

	return nil
}
