package internal

import (
	"bufio"
	"image"
	"image/color"
	"log"
	"math"
	"os"
	"time"
)

const Width = 21600
const height = 10800
const bathymetryFile = "data/bathymetry.bin"

//const width = 2160
//const height = 1080
//const bathymetryFile = "data/bathymetry_2160x1080.bin"

type Grid struct {
	bathymetry [Width][height]byte
	Longitudes [Width]float64
	Latitudes  [height]float64
	Step       float64
}

func NewGrid() *Grid {
	g := Grid{}
	log.Println(time.Now(), "readBathymetry")
	g.readBathymetry()
	log.Println(time.Now(), "end readBathymetry")

	for i, v := range arange(-180, 180, Width) {
		g.Longitudes[i] = v * math.Pi / 180
	}
	g.Longitudes[Width-1] = 180.0 * math.Pi / 180

	var tmLat [height]float64

	for i, v := range arange(-90, 90, height) {
		tmLat[i] = v * math.Pi / 180
	}
	tmLat[height-1] = 90.0 * math.Pi / 180

	for i := range tmLat {
		g.Latitudes[height-i-1] = tmLat[i]
	}

	g.Step = (g.Longitudes[2] - g.Longitudes[1]) * 5

	return &g
}

func (g *Grid) Heuristic(node1, node2 *Node) float64 {
	return math.Sqrt(math.Pow(float64(node1.LatIdx-node2.LatIdx), 2) + math.Pow(float64(node1.LonIdx-node2.LonIdx), 2))
}

func (g *Grid) Haversine(node1, node2 *Node) float64 {
	return Haversine(g.Latitudes[node1.LatIdx], g.Longitudes[node1.LonIdx], g.Latitudes[node2.LatIdx], g.Longitudes[node2.LonIdx])
}

func (g *Grid) IsTraversable(lonIdx, latIdx int) bool {
	lonIdx = normalizeLon(lonIdx)
	return g.bathymetry[lonIdx][latIdx] == 1
}

func (g *Grid) LatLon(node *Node) (float64, float64) {
	return g.Latitudes[node.LatIdx], g.Longitudes[node.LonIdx]
}

func (g *Grid) GetImage(wm *Map) *image.RGBA {
	seaColor := color.RGBA{0, 0, 255, 100}
	seaVisitedColor := color.RGBA{100, 100, 255, 20}
	seaCurrentColor := color.RGBA{255, 255, 255, 70}
	seaInNeibersColor := color.RGBA{0, 255, 0, 120}
	obstacleColor := color.RGBA{144, 12, 63, 255}

	sX := 0
	eX := Width
	sY := 0
	eY := height

	img := image.NewRGBA(image.Rect(sX, sY, eX, eY))

	log.Println(time.Now(), "GetImage")

	for x := sX; x < eX; x++ {
		for y := sY; y < eY; y++ {
			if g.IsTraversable(x, y) {
				if wm.ExistNode(y, x) {

					n := wm.GetNode(y, x)
					if n.InNeibers {
						img.Set(x, y, seaInNeibersColor)
					}
					if n.Visited {
						img.Set(x, y, seaVisitedColor)
					}
					if n.Current {
						img.Set(x, y, seaCurrentColor)
					}

				} else {
					img.Set(x, y, seaColor)
				}
			} else {
				img.Set(x, y, obstacleColor)
			}
		}
	}

	log.Println(time.Now(), "end GetImage")

	return img
}

func (g *Grid) readBathymetry() {
	file, err := os.Open(bathymetryFile)

	if err != nil {
		log.Fatalf("open file error: %v", err)
	}
	defer file.Close()

	stats, err := file.Stat()
	if err != nil {
		log.Fatalf("file stat error: %v", err)
	}

	var size int64 = stats.Size()
	bytes := make([]byte, size)

	bufr := bufio.NewReader(file)
	_, err = bufr.Read(bytes)
	//
	if err != nil {
		log.Fatalf("read buffer error: %v", err)
	}
	//
	for i := 0; i < Width; i++ {
		for j := 0; j < height; j++ {
			g.bathymetry[i][j] = bytes[i*height+(height-j-1)]
			//g.bathymetry[i][j] = bytes[i*height+j]
		}
	}
}

func arange(start, stop float64, count int) []float64 {
	step := (stop - start) / float64(count-1)
	N := count
	rnge := make([]float64, N, N)
	i := 0
	for x := start; x <= stop; x += step {
		rnge[i] = x
		i += 1
	}
	return rnge
}

// Haversine Calculate the great circle arc distance in radians between two points on the Earth.
func Haversine(lat1, lon1, lat2, lon2 float64) float64 {
	diffLat := lat2 - lat1
	diffLon := lon2 - lon1

	a := math.Pow(math.Sin(diffLat/2), 2) + math.Cos(lat1)*math.Cos(lat2)*math.Pow(math.Sin(diffLon/2), 2)

	return 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
}

func (g *Grid) BisecLat(x float64) int {
	lo := 0
	hi := len(g.Latitudes)
	for lo < hi {
		mid := int((lo + hi) / 2)
		if x > g.Latitudes[mid] {
			hi = mid
		} else {
			lo = mid + 1
		}
	}

	if lo > 0 {
		return lo - 1
	}

	return 0
}

func (g *Grid) BisecLon(x float64) int {
	lo := 0
	hi := len(g.Longitudes)
	for lo < hi {
		mid := int((lo + hi) / 2)
		if x < g.Longitudes[mid] {
			hi = mid
		} else {
			lo = mid + 1
		}
	}

	if lo > 0 {
		return lo - 1
	}

	return 0
}
