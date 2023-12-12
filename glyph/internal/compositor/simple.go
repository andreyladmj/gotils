package compositor

import (
	"github.com/andreyladmj/gotils/glyph/internal/graphics"
	"github.com/hajimehoshi/ebiten/v2"
)

type Compositor interface {
	compose(g graphics.Glyph)
}
type SimpleCompositor struct {
	screen *ebiten.Image
}

func NewSimpleCompositor(screen *ebiten.Image) *SimpleCompositor {
	return &SimpleCompositor{screen: screen}
}

func (c *SimpleCompositor) compose(g graphics.Glyph) {
	g.Render(c.screen)
}
