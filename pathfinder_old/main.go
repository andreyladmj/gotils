package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"pathfinder/internal"
	"pathfinder/internal/pathfinder"
	"time"
)

func main() {

	ports := map[string][2]float64{}

	ports["canada"] = [2]float64{48.43254 + 12, -89.21899}
	ports["black sea"] = [2]float64{43.60967, 51.22023 - 15}
	ports["japan"] = [2]float64{35.92849, 140.68824}
	ports["new zeland"] = [2]float64{-19.53236, 169.26580}
	ports["africa"] = [2]float64{-17.29381, 38.13973}
	ports["africa 2"] = [2]float64{21.46388, 39.17029}
	ports["brasil"] = [2]float64{-0.61084, -80.42418 - 0.1}
	ports["turkish"] = [2]float64{35.32798 - 0.1, 33.99319}
	ports["belarus"] = [2]float64{57.04010 + 0.1, 24.07700}

	checks := []struct {
		name string
		p1   [2]float64
		p2   [2]float64
	}{
		{
			name: "canada to black sea",
			p1:   ports["canada"],
			p2:   ports["black sea"],
		},
	}

	for _, check := range checks {
		grid := internal.NewGrid()
		pf := pathfinder.NewPathfinderv8(grid)

		idx1 := pathfinder.Index{LatIdx: grid.BisecLat(pathfinder.Deg2rad(check.p1[0])), LonIdx: grid.BisecLon(pathfinder.Deg2rad(check.p1[1]))}
		idx2 := pathfinder.Index{LatIdx: grid.BisecLat(pathfinder.Deg2rad(check.p2[0])), LonIdx: grid.BisecLon(pathfinder.Deg2rad(check.p2[1]))}
		log.Println("\n\nStart", check.name)
		log.Println("START POINT INDEX", idx1.LatIdx, idx1.LonIdx)
		log.Println("END POINT INDEX", idx2.LatIdx, idx2.LonIdx)
		log.Println("START POINT", grid.Latitudes[idx1.LatIdx], grid.Longitudes[idx1.LonIdx])
		log.Println("END POINT", grid.Latitudes[idx2.LatIdx], grid.Longitudes[idx2.LonIdx])

		start := time.Now()
		s1, s2 := pf.Find(idx1, idx2)
		elapsed := time.Since(start)
		log.Println("time", elapsed)

		img := grid.GetImage(pf.WorldMap)
		c := 0

		node := s1.CurrentNode

		for node != nil && node.Parent != nil && node != node.Parent {
			drawPoint(img, node, false)
			pathfinder.DrawGreatCircle(node, node.Parent, grid, img)
			node = node.Parent
			c++
			if c > 50 {
				break
			}
		}

		node = s2.CurrentNode

		for node != nil && node.Parent != nil && node != node.Parent {
			drawPoint(img, node, false)
			pathfinder.DrawGreatCircle(node, node.Parent, grid, img)
			node = node.Parent
			c++
			if c > 50 {
				break
			}
		}

		drawPoint(img, pf.StartNode, true)
		drawPoint(img, pf.EndNode, true)

		f, err := os.Create(fmt.Sprintf("images/1 Path %s, %s.png", check.name, elapsed))
		if err != nil {
			panic(err)
		}
		defer f.Close()
		png.Encode(f, img)
	}

}

func drawPoint(img *image.RGBA, node *internal.Node, mainPoint bool) {
	dx := 6

	if internal.Width == 2160 {
		dx = 3
	}

	c := color.RGBA{0, 255, 0, 255}

	if mainPoint {
		c = color.RGBA{255, 0, 255, 255}
	}

	img.Set(node.LonIdx, node.LatIdx, c)
	for i := -dx; i < dx; i++ {
		for j := -dx; j < dx; j++ {
			img.Set(node.LonIdx+i, node.LatIdx+j, c)
		}
	}
}
