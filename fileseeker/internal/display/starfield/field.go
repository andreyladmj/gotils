package starfield

import (
	"github.com/andreyladmj/gotils/internal/system"
	"github.com/faiface/pixel/pixelgl"
	"math"
	"math/rand"
)

type StarField struct {
	dir   *system.File
	stars []*Star
	w, h  int
}

func NewStarField(dir *system.File, w, h int) *StarField {
	field := &StarField{dir: dir, w: w, h: h}
	field.parseDir(dir, nil)

	return field
}

func (g *StarField) parseDir(dir *system.File, parent *Star) {
	for _, file := range dir.Files {
		c := NewStar(file.Name, file.Path, file.GetSize(), parent)

		if parent == nil {
			margin := 100
			c.x = float64(getRand(0+margin, g.w-margin))
			c.y = float64(getRand(0+margin, g.h-margin))
		} else {
			alpha := 30.0
			r := 150

			l := float64(len(parent.childs))
			alphaSum := alpha*l + 5*math.Floor(float64(alpha*l)/360)
			rSum := r + r*int(l/(360.0/alpha))

			c.x = parent.x + float64(rSum)*math.Cos(alphaSum*math.Pi/180)
			c.y = parent.y + float64(rSum)*math.Sin(alphaSum*math.Pi/180)
			parent.childs = append(parent.childs, c)
		}

		if file.IsDir {
			g.parseDir(file, c)
		}

		g.stars = append(g.stars, c)
	}
}

func (g *StarField) Update(dt float64, win *pixelgl.Window) error {
	for _, c := range g.stars {
		c.Update(dt, win)
	}

	return nil
}

func (g *StarField) Draw(win *pixelgl.Window) {
	for _, c := range g.stars {
		c.Draw(win)
	}
}

func getRand(min, max int) int {
	return rand.Intn(max+1-min) + min
}
