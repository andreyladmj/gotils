package pathfinder

import (
	"fmt"
	"pathfinder/internal"
	"sync"
)

type Pathfinderv6 struct {
	candidateWeights map[*internal.Node]float64
	lineOfSight      *internal.LineOfSightChecking

	startNode *internal.Node
	endNode   *internal.Node

	startNodeTmp *internal.Node
	endNodeTmp   *internal.Node

	heap  *internal.NodesMinHeap
	found bool

	Pathfinder
}

func NewPathfinderv6(grid *internal.Grid) *Pathfinderv6 {
	return &Pathfinderv6{
		candidateWeights: map[*internal.Node]float64{},
		lineOfSight:      nil,
		startNode:        nil,
		endNode:          nil,
		heap:             nil,
		Pathfinder:       Pathfinder{grid: grid, WorldMap: &internal.Map{Grid: grid}},
	}
}

func (pf *Pathfinderv6) searchNodes(wg *sync.WaitGroup, currentNode, searchNode *internal.Node, isStart bool) {
	heap := internal.NewHeap(10000)
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

		if isStart {
			pf.startNodeTmp = currentNode
		} else {
			pf.endNodeTmp = currentNode
		}

		currentNode.Visited = true
		parent := currentNode.Parent
		neighbours := make([]*internal.Node, 0, 8)
		pf.WorldMap.GetNeighboursAdopted(&neighbours, currentNode)

		c++
		if c > 500000 {
			return
		}

		if currentNode.LatIdx == searchNode.LatIdx && currentNode.LonIdx == searchNode.LonIdx {
			fmt.Println("FINISHED")
			//pf.endNode = currentNode
			pf.found = true
			stop = true
			break
		}

		for _, neighbour := range neighbours {

			if pf.found {
				return
			}

			if neighbour.Parent != nil && neighbour.Parent == parent {
			} else if hasCandidate(candidateWeights, neighbour) {
				candidateWeight := candidateWeights[neighbour]

				if candidateWeight != internal.INFINITY {
					newWeight, updateParent := pf.checkAbilityToMove(parent, currentNode, neighbour)

					if candidateWeight > newWeight {
						pf.updateVertex(currentNode, neighbour, newWeight, updateParent, heap, candidateWeights, searchNode)
					}
				}
			} else {
				if pf.isVisible(currentNode, neighbour) {
					newWeight, updateParent := pf.checkAbilityToMove(parent, currentNode, neighbour)
					pf.updateVertex(currentNode, neighbour, newWeight, updateParent, heap, candidateWeights, searchNode)
				}
			}

			if neighbour.LatIdx == searchNode.LatIdx && neighbour.LonIdx == searchNode.LonIdx {
				fmt.Println("FINISHED")
				//pf.endNode = currentNode
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

func (pf *Pathfinderv6) find(startIndex, endIndex Index) *internal.Node {
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

	pf.endNode.Parent = pf.endNode
	pf.startNode.Parent = pf.startNode

	pf.endNodeTmp = pf.endNode
	pf.startNodeTmp = pf.startNode

	wg.Add(2)
	go pf.searchNodes(&wg, pf.startNode, pf.endNodeTmp, true)
	go pf.searchNodes(&wg, pf.endNode, pf.startNodeTmp, false)
	//go pf.searchNodes(&wg, pf.startNode, pf.endNodeTmp, true)
	//go pf.searchNodes(&wg, pf.endNode, pf.startNodeTmp, false)
	//go pf.searchNodes(&wg, pf.startNode, pf.endNodeTmp, true)
	//go pf.searchNodes(&wg, pf.endNode, pf.startNodeTmp, false)
	//go pf.searchNodes(&wg, pf.startNode, pf.endNodeTmp, true)
	//go pf.searchNodes(&wg, pf.endNode, pf.startNodeTmp, false)
	//go pf.searchNodes(&wg, pf.startNode, pf.endNodeTmp, true)
	//go pf.searchNodes(&wg, pf.endNode, pf.startNodeTmp, false)

	//for _, nodeParent := range pf.WorldMap.GetNeighbours(pf.startNode) {
	//	for _, node := range pf.WorldMap.GetNeighbours(nodeParent) {
	//		wg.Add(1)
	//		c++
	//		node.Parent = pf.startNode
	//		go pf.searchNodes(&wg, node, pf.endNodeTmp, true)
	//	}
	//	break
	//}
	//
	//for _, nodeParent := range pf.WorldMap.GetNeighbours(pf.endNode) {
	//	for _, node := range pf.WorldMap.GetNeighbours(nodeParent) {
	//		wg.Add(1)
	//		c++
	//		node.Parent = pf.endNode
	//		go pf.searchNodes(&wg, node, pf.startNodeTmp, true)
	//	}
	//	break
	//}

	fmt.Println("Start", c, "threads")

	wg.Wait()

	return pf.endNode
}

func (pf *Pathfinderv6) updateVertex(currentNode, neighbour *internal.Node, weight float64, updateParent bool, heap *internal.NodesMinHeap, candidateWeights map[*internal.Node]float64, searchNode *internal.Node) {
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
	neighbour.GScore = weight + pf.grid.Haversine(neighbour, searchNode)

	heap.Insert(neighbour)
}

func (pf *Pathfinderv6) checkAbilityToMove(parent, currentNode, neighbour *internal.Node) (float64, bool) {
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

func (pf *Pathfinderv6) hasCandidate(node *internal.Node) bool {
	_, ok := pf.candidateWeights[node]
	return ok
}

func (pf *Pathfinderv6) updateWeight(node *internal.Node, weight float64) {
	pf.candidateWeights[node] = weight
}

func (pf *Pathfinderv6) removeCandidate(node *internal.Node) {
	delete(pf.candidateWeights, node)
}

func (pf *Pathfinderv6) isVisible(currentNode, neighbour *internal.Node) bool {

	//return pf.lineOfSight.lineOfSight(currentNode, neighbour, true, true)

	lat1, lon1 := pf.grid.LatLon(currentNode)
	lat2, lon2 := pf.grid.LatLon(neighbour)
	return fillGreatCircle(lat1, lon1, lat2, lon2, pf.grid)
}
