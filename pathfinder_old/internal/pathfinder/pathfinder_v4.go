package pathfinder

import (
	"fmt"
	"math"
	"pathfinder/internal"
	"sync"
)

type Pathfinderv4 struct {
	candidateWeights map[*internal.Node]float64
	lineOfSight      *internal.LineOfSightChecking

	startNode *internal.Node
	endNode   *internal.Node
	heap      *internal.NodesMinHeap
	found     bool

	Pathfinder
}

func NewPathfinder4(grid *internal.Grid) *Pathfinderv4 {
	return &Pathfinderv4{
		candidateWeights: map[*internal.Node]float64{},
		lineOfSight:      nil,
		startNode:        nil,
		endNode:          nil,
		heap:             nil,
		Pathfinder:       Pathfinder{grid: grid, WorldMap: &internal.Map{Grid: grid}},
	}
}

func (pf *Pathfinderv4) searchNodes(wg *sync.WaitGroup, currentNode *internal.Node) {
	heap := internal.NewHeap(25000)
	heap.Clear()
	heap.Insert(currentNode)
	stop := false
	candidateWeights := map[*internal.Node]float64{}
	defer wg.Done()
	c := 0

	for !heap.Empty() {
		currentNode = heap.GetMin()

		if currentNode.Visited {
			continue
		}

		currentNode.Visited = true
		parent := currentNode.Parent
		neighbours := make([]*internal.Node, 0, 16)
		pf.WorldMap.GetNeighboursAdoptedV2(&neighbours, currentNode)

		for _, neighbour := range neighbours {

			if pf.found {
				return
			}

			c++
			if c > 1000000 {
				return
			}

			if neighbour.Parent != nil && neighbour.Parent == parent {
			} else if hasCandidate(candidateWeights, neighbour) {
				candidateWeight := candidateWeights[neighbour]

				if candidateWeight != internal.INFINITY {
					newWeight, updateParent := pf.checkAbilityToMove(parent, currentNode, neighbour)

					if candidateWeight > newWeight {
						pf.updateVertex(currentNode, neighbour, newWeight, updateParent, heap, candidateWeights)
					}
				}
			} else {
				if pf.isVisible(currentNode, neighbour) {
					newWeight, updateParent := pf.checkAbilityToMove(parent, currentNode, neighbour)
					pf.updateVertex(currentNode, neighbour, newWeight, updateParent, heap, candidateWeights)
				}
			}

			if neighbour.LatIdx == pf.endNode.LatIdx && neighbour.LonIdx == pf.endNode.LonIdx {
				fmt.Println("FINISHED")
				pf.endNode = currentNode
				pf.found = true
				stop = true
				break
			}
		}

		if stop {
			break
		}
	}
}

func (pf *Pathfinderv4) find(startIndex, endIndex Index) *internal.Node {
	pf.startNode = pf.WorldMap.GetNode(startIndex.LatIdx, startIndex.LonIdx)
	pf.endNode = pf.WorldMap.GetNode(endIndex.LatIdx, endIndex.LonIdx)

	pf.startNode.Parent = pf.startNode
	pf.found = false

	fmt.Println("startNode", pf.startNode.LatIdx, pf.startNode.LonIdx)
	fmt.Println("endNode", pf.endNode.LatIdx, pf.endNode.LonIdx)

	pf.startNode.GScore = pf.grid.Haversine(pf.startNode, pf.endNode)
	pf.startNode.FScore = 0
	var wg sync.WaitGroup
	c := 0

	for _, nodeParent := range pf.WorldMap.GetNeighbours(pf.startNode) {
		for _, node := range pf.WorldMap.GetNeighbours(nodeParent) {
			wg.Add(1)
			c++
			node.Parent = pf.startNode
			go pf.searchNodes(&wg, node)
		}
	}

	fmt.Println("Start", c, "threads")

	wg.Wait()

	return pf.endNode
}

func (pf *Pathfinderv4) updateVertex(currentNode, neighbour *internal.Node, weight float64, updateParent bool, heap *internal.NodesMinHeap, candidateWeights map[*internal.Node]float64) {
	if hasCandidate(candidateWeights, neighbour) {
		removeCandidate(candidateWeights, neighbour)
	}

	updateWeight(candidateWeights, neighbour, weight)

	if updateParent {
		neighbour.Parent = currentNode
	} else {
		neighbour.Parent = currentNode.Parent
	}

	neighbour.FScore = weight
	neighbour.GScore = weight + pf.grid.Haversine(neighbour, pf.endNode)

	heap.Insert(neighbour)
}

func (pf *Pathfinderv4) checkAbilityToMove(parent, currentNode, neighbour *internal.Node) (float64, bool) {
	var newWeight float64
	var updateParent bool

	if pf.isVisible(parent, neighbour) {
		newWeight = parent.FScore + pf.grid.Haversine(parent, neighbour)
		updateParent = false
	} else {
		newWeight = currentNode.FScore + pf.grid.Haversine(currentNode, neighbour)
		updateParent = true
	}

	return newWeight, updateParent
}

func (pf *Pathfinderv4) hasCandidate(node *internal.Node) bool {
	_, ok := pf.candidateWeights[node]
	return ok
}

func (pf *Pathfinderv4) updateWeight(node *internal.Node, weight float64) {
	pf.candidateWeights[node] = weight
}

func (pf *Pathfinderv4) removeCandidate(node *internal.Node) {
	delete(pf.candidateWeights, node)
}

func (pf *Pathfinderv4) isVisible(currentNode, neighbour *internal.Node) bool {
	lat1, lon1 := pf.grid.LatLon(currentNode)
	lat2, lon2 := pf.grid.LatLon(neighbour)
	return fillGreatCircle(lat1, lon1, lat2, lon2, pf.grid)
}

func fillGreatCircle(lat1, lon1, lat2, lon2 float64, grid *internal.Grid) bool {
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

	getGCLatLon := func(f float64) (float64, float64) {
		a := math.Sin((1 - f) * d)
		b := math.Sin(f * d)
		x := a*cos_lat1_cos_lon1 + b*cos_lat2_cos_lon2
		y := a*cos_lat1_sin_lon1 + b*cos_lat2_sin_lon2
		z := a*sin_lat1 + b*sin_lat2
		lat := math.Atan2(z, math.Pow(math.Pow(x, 2)+math.Pow(y, 2), 0.5))
		lon := math.Atan2(y, x)
		return lat, lon
	}

	//N := int(1 / step)

	//if N < 50 {
	for f := step; f <= 1+1e-9; f += step {
		lat, lon := getGCLatLon(f)
		latIdx := grid.BisecLat(lat)
		lonIdx := grid.BisecLon(lon)

		if !grid.IsTraversable(lonIdx, latIdx) {
			return false
		}

		if f > step {
			prevLat, prevLon := getGCLatLon(f - step)
			prevLatIdx := grid.BisecLat(prevLat)
			prevLonIdx := grid.BisecLon(prevLon)

			if !internal.TraversableSegmentOnPlane(lonIdx, latIdx, prevLonIdx, prevLatIdx, grid) {
				return false
			}
		}
	}

	return true
	//}

	//var wg sync.WaitGroup
	//nonTraversablePoints := 0
	//
	//for f:=step;f<=1+1e-9;f+=step {
	//	wg.Add(1)
	//	go func(wg *sync.WaitGroup, f float64) {
	//		defer wg.Done()
	//
	//		lat, lon := getGCLatLon(f)
	//		latIdx := grid.BisecLat(lat)
	//		lonIdx := grid.BisecLon(lon)
	//
	//		if !grid.IsTraversable(lonIdx, latIdx) {
	//			nonTraversablePoints++
	//		}
	//
	//		if nonTraversablePoints > 0 {
	//			return
	//		}
	//
	//		if f > step {
	//			prevLat, prevLon := getGCLatLon(f-step)
	//			prevLatIdx := grid.BisecLat(prevLat)
	//			prevLonIdx := grid.BisecLon(prevLon)
	//
	//			if !traversableSegmentOnPlane(lonIdx, latIdx, prevLonIdx, prevLatIdx, grid) {
	//				nonTraversablePoints++
	//			}
	//		}
	//	}(&wg, f)
	//}
	//
	//wg.Wait()
	//
	//return nonTraversablePoints == 0
}

func hasCandidate(candidateWeights map[*internal.Node]float64, node *internal.Node) bool {
	_, ok := candidateWeights[node]
	return ok
}

func updateWeight(candidateWeights map[*internal.Node]float64, node *internal.Node, weight float64) {
	candidateWeights[node] = weight
}

func removeCandidate(candidateWeights map[*internal.Node]float64, node *internal.Node) {
	delete(candidateWeights, node)
}
