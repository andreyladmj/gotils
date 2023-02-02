package pathfinder

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"pathfinder/internal"
)

type Pathfinderv3 struct {
	candidateWeights map[*internal.Node]float64
	lineOfSight      *internal.LineOfSightChecking

	startNode *internal.Node
	endNode   *internal.Node
	heap      *internal.NodesMinHeap

	Pathfinder
}

func NewPathfinderv3(grid *internal.Grid) *Pathfinderv3 {
	return &Pathfinderv3{
		candidateWeights: map[*internal.Node]float64{},
		lineOfSight:      nil,
		startNode:        nil,
		endNode:          nil,
		heap:             nil,
		Pathfinder:       Pathfinder{grid: grid, WorldMap: &internal.Map{Grid: grid}},
	}
}

func (pf *Pathfinderv3) find(startIndex, endIndex Index) *internal.Node {
	pf.heap = internal.NewHeap(1000)

	//pf.lineOfSight = NewLineOfSight(pf.grid)

	pf.startNode = pf.WorldMap.GetNode(startIndex.LatIdx, startIndex.LonIdx)
	pf.endNode = pf.WorldMap.GetNode(endIndex.LatIdx, endIndex.LonIdx)

	pf.startNode.Parent = pf.startNode

	fmt.Println("startNode", pf.startNode.LatIdx, pf.startNode.LonIdx)
	fmt.Println("endNode", pf.endNode.LatIdx, pf.endNode.LonIdx)

	pf.startNode.GScore = pf.grid.Haversine(pf.startNode, pf.endNode)
	pf.startNode.FScore = 0
	pf.heap.Clear()
	pf.heap.Insert(pf.startNode)

	var currentNode *internal.Node
	stop := false
	c := 0

	for !pf.heap.Empty() {
		currentNode = pf.heap.GetMin()

		if currentNode == nil {
			log.Fatalf("Path wasnt found")
		}

		if currentNode.Visited {
			continue
		}

		//fmt.Println("currentNode", currentNode.LatIdx, currentNode.LonIdx, "heap size", heap.heapSize, "visited", currentNode.Visited)
		currentNode.Visited = true
		parent := currentNode.Parent

		//dl := currentNode.Size
		//if pf.endNode.LatIdx > currentNode.LatIdx - dl && pf.endNode.LatIdx < currentNode.LatIdx + dl &&
		//	pf.endNode.LonIdx > currentNode.LonIdx - dl && pf.endNode.LonIdx < currentNode.LonIdx + dl {
		//	//pf.endNode = currentNode
		//	stop = true
		//	break
		//}

		neighbours := make([]*internal.Node, 0, 8)
		pf.WorldMap.GetNeighboursAdopted(&neighbours, currentNode)

		for _, neighbour := range neighbours {
			//for _, neighbour := range pf.WorldMap.GetNeighbours(currentNode) {
			c++

			//if c == 10000 {
			//	return currentNode
			//}

			if neighbour.Parent != nil && neighbour.Parent == parent {

			} else if pf.hasCandidate(neighbour) {
				candidateWeight := pf.candidateWeights[neighbour]

				if candidateWeight != internal.INFINITY {
					newWeight, updateParent := pf.checkAbilityToMove(parent, currentNode, neighbour)

					if candidateWeight > newWeight {
						pf.updateVertex(currentNode, neighbour, newWeight, updateParent)
					}
				}
			} else {
				if pf.isVisible(currentNode, neighbour) {
					newWeight, updateParent := pf.checkAbilityToMove(parent, currentNode, neighbour)
					pf.updateVertex(currentNode, neighbour, newWeight, updateParent)
				}
			}
			//
			//dl := neighbour.Size
			//if pf.endNode.LatIdx > neighbour.LatIdx - dl && pf.endNode.LatIdx < neighbour.LatIdx + dl &&
			//	pf.endNode.LonIdx > neighbour.LonIdx - dl && pf.endNode.LonIdx < neighbour.LonIdx + dl {
			//	//pf.endNode = currentNode
			//	stop = true
			//	break
			//}

			if neighbour.LatIdx == pf.endNode.LatIdx && neighbour.LonIdx == pf.endNode.LonIdx {
				//pf.endNode = neighbour
				pf.endNode = currentNode
				stop = true
				break
			}
		}

		if stop {
			break
		}

	}

	return pf.endNode
}

func (pf *Pathfinderv3) updateVertex(currentNode, neighbour *internal.Node, weight float64, updateParent bool) {
	if pf.hasCandidate(neighbour) {
		pf.removeCandidate(neighbour)
	}

	pf.updateWeight(neighbour, weight)

	if updateParent {
		neighbour.Parent = currentNode
	} else {
		neighbour.Parent = currentNode.Parent
	}

	neighbour.FScore = weight
	neighbour.GScore = weight + pf.grid.Haversine(neighbour, pf.endNode)

	pf.heap.Insert(neighbour)
}

func (pf *Pathfinderv3) checkAbilityToMove(parent, currentNode, neighbour *internal.Node) (float64, bool) {
	var newWeight float64
	var updateParent bool

	if pf.isVisible(parent, neighbour) {
		newWeight = parent.FScore + 0.5*pf.grid.Haversine(parent, neighbour)
		updateParent = false
	} else {
		newWeight = currentNode.FScore + 0.5*pf.grid.Haversine(currentNode, neighbour)
		updateParent = true
	}

	return newWeight, updateParent
}

func (pf *Pathfinderv3) hasCandidate(node *internal.Node) bool {
	_, ok := pf.candidateWeights[node]
	return ok
}

func (pf *Pathfinderv3) updateWeight(node *internal.Node, weight float64) {
	pf.candidateWeights[node] = weight
}

func (pf *Pathfinderv3) removeCandidate(node *internal.Node) {
	delete(pf.candidateWeights, node)
}

func (pf *Pathfinderv3) isVisible(currentNode, neighbour *internal.Node) bool {
	lat1, lon1 := pf.grid.LatLon(currentNode)
	lat2, lon2 := pf.grid.LatLon(neighbour)

	//return fill_great_circle(deg2rad(lat1), deg2rad(lon1), deg2rad(lat2), deg2rad(lon2), pf.grid)
	return fill_great_circle(lat1, lon1, lat2, lon2, pf.grid)
}

func Deg2rad(degrees float64) float64 {
	return degrees * (math.Pi / 180.0)
}

func rad2deg(radians float64) float64 {
	return radians * (180.0 / math.Pi)
}

func fill_great_circle(lat1, lon1, lat2, lon2 float64, grid *internal.Grid) bool {
	// d between lats 0.0029102294150895602
	// gc step 0.07267292292361319
	// AVG_EARTH_RADIUS_NM = 3440.0707986216
	step := grid.Step

	cos_lat1 := math.Cos(lat1)
	sin_lat1 := math.Sin(lat1)
	cos_lon1 := math.Cos(lon1)
	sin_lon1 := math.Sin(lon1)

	cos_lat2 := math.Cos(lat2)
	sin_lat2 := math.Sin(lat2)
	cos_lon2 := math.Cos(lon2)
	sin_lon2 := math.Sin(lon2)

	d := math.Acos(sin_lat1*sin_lat2 + cos_lat1*cos_lat2*math.Cos(lon1-lon2))
	sind := math.Sin(d)
	if sind == 0 {
		d += 1e-9
	}

	cos_lat1_cos_lon1 := cos_lat1 * cos_lon1
	cos_lat1_sin_lon1 := cos_lat1 * sin_lon1
	cos_lat2_cos_lon2 := cos_lat2 * cos_lon2
	cos_lat2_sin_lon2 := cos_lat2 * sin_lon2

	step /= d
	//n := int(1 / step)
	//range_ = np.arange(step, 1 + 1e-9, step)

	nonTraversablePoints := 0
	//var wg sync.WaitGroup

	getGCLatLon := func(f float64) (float64, float64) {
		a := math.Sin((1 - f) * d)
		b := math.Sin(f * d)
		x := a*cos_lat1_cos_lon1 + b*cos_lat2_cos_lon2
		y := a*cos_lat1_sin_lon1 + b*cos_lat2_sin_lon2
		z := a*sin_lat1 + b*sin_lat2
		lat := math.Atan2(z, math.Pow(math.Pow(x, 2)+math.Pow(y, 2), 0.5))
		lon := math.Atan2(y, x)
		//return rad2deg(lat), rad2deg(lon)
		return lat, lon
	}

	for f := step; f <= 1+1e-9; f += step {
		//wg.Add(1)
		//go func(wg *sync.WaitGroup, f float64) {
		//	defer wg.Done()

		//a := math.Sin((1 - f) * d)
		//b := math.Sin(f * d)
		//x := a * cos_lat1_cos_lon1 + b * cos_lat2_cos_lon2
		//y := a * cos_lat1_sin_lon1 + b * cos_lat2_sin_lon2
		//z := a * sin_lat1 + b * sin_lat2
		//lat_ := math.Atan2(z, math.Pow(math.Pow(x, 2) + math.Pow(y, 2), 0.5))
		//lon_ := math.Atan2(y, x)
		//lat, lon = rad2deg(lat), rad2deg(lon)
		//
		//fmt.Println(lat_, lon_)

		lat, lon := getGCLatLon(f)
		latIdx := grid.BisecLat(lat)
		lonIdx := grid.BisecLon(lon)

		if !grid.IsTraversable(lonIdx, latIdx) {
			return false
			//nonTraversablePoints++
		}

		if f > step {
			prevLat, prevLon := getGCLatLon(f - step)
			//prevLat, prevLon = rad2deg(prevLat), rad2deg(prevLon)
			prevLatIdx := grid.BisecLat(prevLat)
			prevLonIdx := grid.BisecLon(prevLon)

			if !internal.TraversableSegmentOnPlane(lonIdx, latIdx, prevLonIdx, prevLatIdx, grid) {
				return false
				//nonTraversablePoints++
			}
		}

		//}(&wg, f)
	}

	//wg.Wait()

	return nonTraversablePoints == 0
}

func DrawGreatCircle(node1, node2 *internal.Node, grid *internal.Grid, img *image.RGBA) {
	// d between lats 0.0029102294150895602
	// gc step 0.07267292292361319
	// AVG_EARTH_RADIUS_NM = 3440.0707986216
	step := grid.Step

	//lat1 := deg2rad(grid.Latitudes[node1.LatIdx])
	//lon1 := deg2rad(grid.Longitudes[node1.LonIdx])
	//lat2 := deg2rad(grid.Latitudes[node2.LatIdx])
	//lon2 := deg2rad(grid.Longitudes[node2.LonIdx])
	lat1 := grid.Latitudes[node1.LatIdx]
	lon1 := grid.Longitudes[node1.LonIdx]
	lat2 := grid.Latitudes[node2.LatIdx]
	lon2 := grid.Longitudes[node2.LonIdx]

	cos_lat1 := math.Cos(lat1)
	sin_lat1 := math.Sin(lat1)
	cos_lon1 := math.Cos(lon1)
	sin_lon1 := math.Sin(lon1)

	cos_lat2 := math.Cos(lat2)
	sin_lat2 := math.Sin(lat2)
	cos_lon2 := math.Cos(lon2)
	sin_lon2 := math.Sin(lon2)

	d := math.Acos(sin_lat1*sin_lat2 + cos_lat1*cos_lat2*math.Cos(lon1-lon2))
	sind := math.Sin(d)
	if sind == 0 {
		d += 1e-9
	}

	cos_lat1_cos_lon1 := cos_lat1 * cos_lon1
	cos_lat1_sin_lon1 := cos_lat1 * sin_lon1
	cos_lat2_cos_lon2 := cos_lat2 * cos_lon2
	cos_lat2_sin_lon2 := cos_lat2 * sin_lon2

	step /= d
	//n := int(1 / step)
	//range_ = np.arange(step, 1 + 1e-9, step)

	c := color.RGBA{255, 255, 0, 255}
	dx := 2

	if internal.Width == 2160 {
		dx = 1
	}

	//wg.Add(n)
	for f := step; f <= 1+1e-9; f += step {

		a := math.Sin((1 - f) * d)
		b := math.Sin(f * d)
		x := a*cos_lat1_cos_lon1 + b*cos_lat2_cos_lon2
		y := a*cos_lat1_sin_lon1 + b*cos_lat2_sin_lon2
		z := a*sin_lat1 + b*sin_lat2
		lat := math.Atan2(z, math.Pow(math.Pow(x, 2)+math.Pow(y, 2), 0.5))
		lon := math.Atan2(y, x)

		//latIdx := grid.BisecLat(rad2deg(lat))
		//lonIdx := grid.BisecLon(rad2deg(lon))
		latIdx := grid.BisecLat(lat)
		lonIdx := grid.BisecLon(lon)

		for i := -dx; i < dx; i++ {
			for j := -dx; j < dx; j++ {
				img.Set(lonIdx+i, latIdx+j, c)
			}
		}
	}
}
