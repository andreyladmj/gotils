package pathfinder

import (
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

	pf.startNode = pf.WorldMap.GetNode(startIndex.LatIdx, startIndex.LonIdx)
	pf.endNode = pf.WorldMap.GetNode(endIndex.LatIdx, endIndex.LonIdx)

	pf.startNode.Parent = pf.startNode

	log.Println("startNode", pf.startNode.LatIdx, pf.startNode.LonIdx)
	log.Println("endNode", pf.endNode.LatIdx, pf.endNode.LonIdx)

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

		currentNode.Visited = true
		parent := currentNode.Parent

		neighbours := make([]*internal.Node, 0, 8)
		pf.WorldMap.GetNeighboursAdopted(&neighbours, currentNode)

		for _, neighbour := range neighbours {
			c++

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

			if neighbour.LatIdx == pf.endNode.LatIdx && neighbour.LonIdx == pf.endNode.LonIdx {
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
	return fillGreatCircle(lat1, lon1, lat2, lon2, pf.grid)
}

func Deg2rad(degrees float64) float64 {
	return degrees * (math.Pi / 180.0)
}

func rad2deg(radians float64) float64 {
	return radians * (180.0 / math.Pi)
}
