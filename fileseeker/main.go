package main

import (
	"flag"
	"fmt"
	"github.com/andreyladmj/gotils/env"
	"github.com/andreyladmj/gotils/internal"
	"github.com/andreyladmj/gotils/internal/display"
	"github.com/andreyladmj/gotils/internal/display/starfield"
	"github.com/andreyladmj/gotils/internal/scanners"
	"github.com/andreyladmj/gotils/internal/system"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	log "github.com/sirupsen/logrus"
	"golang.org/x/image/colornames"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	_, err := env.ReadOS()
	if err != nil {
		log.Fatal("config read error:", err)
	}

	fmt.Println("basepath: ", os.Args[1])

	displayer := flag.String("display", "cmd", "starfield or cmd")
	minSize := flag.String("min-size", "1Kb", "1Kb, 1Mb, 1Gb, default: 1Kb")
	flag.Parse()

	filters := &internal.Filters{}
	filters.FileSizeFilter(parseMinSize(*minSize))
	scanner := scanners.NewDirsRecursiveScanner(0)
	scanner.SetFilter(filters)
	scanner.Scan(os.Args[1])

	switch *displayer {
	case "starfield":
		pixelgl.Run(func() {
			run(scanner.GetDir())
		})
	default:
		p := display.NewCommandLine()
		p.Render(scanner.GetDir())
	}
}

func run(dir *system.File) {

	cfg := pixelgl.WindowConfig{
		Title:  "Starfield!",
		Bounds: pixel.R(0, 0, 1920, 1080),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	sf := starfield.NewStarField(dir, int(win.Bounds().W()), int(win.Bounds().H()))
	win.SetSmooth(true)
	last := time.Now()

	for !win.Closed() {
		win.Clear(colornames.Black)

		dt := time.Since(last).Seconds()
		last = time.Now()

		sf.Update(dt, win)
		sf.Draw(win)
		win.Update()
	}
}

func parseMinSize(minSize string) (int, string) {
	if len(minSize) < 3 {
		log.Fatal("wrong size format")
	}

	units := minSize[len(minSize)-2:]
	size, err := strconv.Atoi(minSize[:len(minSize)-2])

	if err != nil {
		log.Fatal(err)
	}

	return size, units
}
