package graphics

import (
	"bufio"
	"fmt"
	"github.com/andreyladmj/gotils/glyph/internal/internal"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
	"image/color"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

type Words struct {
	Glyph
	font     *font.Face
	index    int
	text     string
	words    []string
	lines    []string
	time     time.Time
	duration int
}

func NewWords(pos internal.Pos, w, h float32, file string, font *font.Face) *Words {
	g := &Words{Glyph: Glyph{pos: pos, w: w, h: h}, font: font, time: time.Now(), duration: 60}
	err := g.ReadFile(file)
	rand.Seed(time.Now().Unix())
	g.text = g.lines[rand.Int()%len(g.lines)]

	if err != nil {
		log.Fatal(err)
		return nil
	}

	return g
}

func (g *Words) Update(input *internal.WInput) error {
	if time.Now().Sub(g.time).Seconds() > float64(g.duration) {
		g.index++
		g.text = g.lines[rand.Int()%len(g.lines)]
		g.time = time.Now()
	}

	for i := 0; i < len(g.children); i++ {
		if err := g.children[i].Update(input); err != nil {
			return err
		}
	}

	return nil
}

func (g *Words) Render(screen *ebiten.Image) {
	text.Draw(screen, g.text, *g.font, int(g.pos.X)+20, int(g.pos.Y)+50, color.White)
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

func (g *Words) ReadFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	start := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if strings.Contains(line, "[START]") {
			start = true
			continue
		}

		if start && line != "" {
			g.lines = append(g.lines, line)
		}

		if strings.Contains(line, "[END]") {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println(time.Now().Sub(g.time).Seconds())

	return nil
}
