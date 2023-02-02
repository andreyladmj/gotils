package pathfinder

import (
	"fmt"
	"log"
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

	log.Println("startNode", pf.startNode.LatIdx, pf.startNode.LonIdx)
	log.Println("endNode", pf.endNode.LatIdx, pf.endNode.LonIdx)

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

	log.Println("Start", c, "threads")

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
