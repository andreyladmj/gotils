package starfield

import (
	"fmt"
	"github.com/andreyladmj/gotils/internal/display"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/font/basicfont"
	"image"
	"image/png"
	"math"
	"os"
)

var pic pixel.Picture

func init() {
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
	var err error
	pic, err = loadPicture("utils/Shining-Star-PNG-ksndfhbfg.png")
	if err != nil {
		panic(err)
	}
}

type Star struct {
	name     string
	path     string
	byteSize int

	angle              float64
	angleRotationSpeed float64
	x, y               float64

	hovered      bool
	clicked      bool
	parent       *Star
	childs       []*Star
	spriteScaled float64
	sprite       *pixel.Sprite
}

func NewStar(name, path string, byteSize int, parent *Star) *Star {
	return &Star{
		name:     name,
		path:     path,
		byteSize: byteSize,
		parent:   parent,
		sprite:   pixel.NewSprite(pic, pic.Bounds()),

		spriteScaled:       geSpriteScale(byteSize),
		angleRotationSpeed: float64(getRand(3, 7)) / 10,
	}
}

func geSpriteScale(byteSize int) float64 {
	return math.Min(0.02+float64(byteSize)/1024/1024/1024, 0.05)
}

func (s *Star) SetPosition(x, y float64) {
	s.x = x
	s.y = y
}

func (s *Star) GetSize() (w, h float64) {
	return s.sprite.Frame().W() * s.spriteScaled, s.sprite.Frame().H() * s.spriteScaled
}

func (s *Star) Update(dt float64, win *pixelgl.Window) {
	mousePos := win.MousePosition()

	s.angle += s.angleRotationSpeed * dt
	w, h := s.GetSize()

	if mousePos.X >= s.x-w/2 && mousePos.X <= s.x+w/2 && mousePos.Y >= s.y-h/2 && mousePos.Y <= s.y+h/2 {
		s.hovered = true
	} else {
		s.hovered = false
	}

	if s.hovered && win.JustPressed(pixelgl.MouseButtonLeft) {
		s.clicked = true
	} else {
		s.clicked = false
	}
}

func (s *Star) DrawLinesToChilds(win *pixelgl.Window) {
	for _, cs := range s.childs {
		f := pixel.V(s.x, s.y)
		t := pixel.V(cs.x, cs.y)
		imd := imdraw.New(nil)
		imd.Color = pixel.RGBA{173, 0, 0, 1}
		imd.Push(f, t)
		imd.Line(1)
		imd.Draw(win)
	}
}

func (s *Star) DrawLineToParent(win *pixelgl.Window) {
	f := pixel.V(s.x, s.y)
	t := pixel.V(s.parent.x, s.parent.y)
	imd := imdraw.New(nil)
	imd.Color = pixel.RGB(173, 216, 230)
	imd.Push(f, t)
	imd.Line(1)
	imd.Draw(win)
}

func (s *Star) Draw(win *pixelgl.Window) {
	mat := pixel.IM
	mat = mat.Rotated(pixel.ZV, s.angle)
	mat = mat.Moved(pixel.V(s.x, s.y))
	mat = mat.ScaledXY(pixel.V(s.x, s.y), pixel.V(s.spriteScaled, s.spriteScaled))
	s.sprite.Draw(win, mat)
	gbSize := fmt.Sprintf("%.2f", display.Format(float64(s.byteSize), "gb"))

	if s.hovered {
		s.DrawLinesToChilds(win)

		if s.parent != nil {
			s.DrawLineToParent(win)
		}

		basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
		basicTxt := text.New(pixel.V(20, win.Bounds().H()-20), basicAtlas)
		fmt.Fprintln(basicTxt, fmt.Sprintln(s.path, "(", gbSize, "Gb", ")"))
		basicTxt.Draw(win, pixel.IM)
	}
}

func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}
