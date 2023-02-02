package pathfinder

import (
	"fmt"
	"log"
	"pathfinder/internal"
)

type Pathfinderv2 struct {
	candidateWeights map[*internal.Node]float64
	lineOfSight      *internal.LineOfSightChecking

	startNode *internal.Node
	endNode   *internal.Node
	heap      *internal.NodesMinHeap

	Pathfinder
}

func NewPathfinder2(grid *internal.Grid) *Pathfinderv2 {
	return &Pathfinderv2{
		candidateWeights: map[*internal.Node]float64{},
		lineOfSight:      nil,
		startNode:        nil,
		endNode:          nil,
		heap:             nil,
		Pathfinder:       Pathfinder{grid: grid, WorldMap: &internal.Map{Grid: grid}},
	}
}

func (pf *Pathfinderv2) find(startIndex, endIndex Index) *internal.Node {
	pf.heap = internal.NewHeap(25000)

	pf.lineOfSight = internal.NewLineOfSight(pf.grid)

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

		for _, neighbour := range pf.WorldMap.GetNeighbours(currentNode) {

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

func (pf *Pathfinderv2) isVisible(currentNode, neighbour *internal.Node) bool {
	return pf.lineOfSight.LineOfSight(currentNode, neighbour, true, true)
}

func (pf *Pathfinderv2) updateVertex(currentNode, neighbour *internal.Node, weight float64, updateParent bool) {
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

func (pf *Pathfinderv2) checkAbilityToMove(parent, currentNode, neighbour *internal.Node) (float64, bool) {
	var newWeight float64
	var updateParent bool

	if pf.lineOfSight.LineOfSight(parent, neighbour, true, true) {
		newWeight = parent.FScore + 0.5*pf.grid.Haversine(parent, neighbour)
		updateParent = false
	} else {
		newWeight = currentNode.FScore + 0.5*pf.grid.Haversine(currentNode, neighbour)
		updateParent = true
	}

	return newWeight, updateParent
}

func (pf *Pathfinderv2) hasCandidate(node *internal.Node) bool {
	_, ok := pf.candidateWeights[node]
	return ok
}

func (pf *Pathfinderv2) updateWeight(node *internal.Node, weight float64) {
	pf.candidateWeights[node] = weight
}

func (pf *Pathfinderv2) removeCandidate(node *internal.Node) {
	delete(pf.candidateWeights, node)
}
