package pathfinder

import (
	"fmt"
	"pathfinder/internal"
	"sync"
)

type Pathfinderv7 struct {
	candidateWeights map[*internal.Node]float64
	lineOfSight      *internal.LineOfSightChecking

	startNode *internal.Node
	endNode   *internal.Node

	g1StartNode *internal.Node
	g1EndNode   *internal.Node

	g2StartNode *internal.Node
	g2EndNode   *internal.Node

	g3StartNode *internal.Node
	g3EndNode   *internal.Node

	g4StartNode *internal.Node
	g4EndNode   *internal.Node

	heap  *internal.NodesMinHeap
	found bool

	Pathfinder
}

func NewPathfinderv7(grid *internal.Grid) *Pathfinderv7 {
	return &Pathfinderv7{
		candidateWeights: map[*internal.Node]float64{},
		lineOfSight:      nil,
		startNode:        nil,
		endNode:          nil,
		heap:             nil,
		Pathfinder:       Pathfinder{grid: grid, WorldMap: &internal.Map{Grid: grid}},
	}
}

func (pf *Pathfinderv7) searchNodes(wg *sync.WaitGroup, currentNode, searchNode *internal.Node, isStart bool, gn int) {
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
			if gn == 1 {
				pf.g1StartNode = currentNode
			}
			if gn == 2 {
				pf.g2StartNode = currentNode
			}
			if gn == 3 {
				pf.g3StartNode = currentNode
			}
			if gn == 4 {
				pf.g4StartNode = currentNode
			}
		} else {
			if gn == 1 {
				pf.g1EndNode = currentNode
			}
			if gn == 2 {
				pf.g2EndNode = currentNode
			}
			if gn == 3 {
				pf.g3EndNode = currentNode
			}
			if gn == 4 {
				pf.g4EndNode = currentNode
			}
		}

		currentNode.Visited = true
		parent := currentNode.Parent
		neighbours := make([]*internal.Node, 0, 8)
		//pf.WorldMap.GetNeighboursAdoptedV2(&neighbours, currentNode)
		//pf.WorldMap.GetMinNeighbour(&neighbours, currentNode)
		//pf.WorldMap.GetNeighboursAdopted(&neighbours, currentNode)
		pf.WorldMap.GetNeighboursAdoptedFarAway(&neighbours, currentNode)

		c++
		if c > 20000 {
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

func (pf *Pathfinderv7) find(startIndex, endIndex Index) *internal.Node {
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

	pf.g1EndNode = pf.endNode
	pf.g2EndNode = pf.endNode
	pf.g3EndNode = pf.endNode
	pf.g4EndNode = pf.endNode
	pf.g1StartNode = pf.startNode
	pf.g2StartNode = pf.startNode
	pf.g3StartNode = pf.startNode
	pf.g4StartNode = pf.startNode

	wg.Add(8)
	//go pf.searchNodes(&wg, pf.startNode, pf.g1EndNode, true, 1)
	//go pf.searchNodes(&wg, pf.endNode, pf.g1StartNode, false, 1)
	//
	//go pf.searchNodes(&wg, pf.startNode, pf.g2EndNode, true, 2)
	//go pf.searchNodes(&wg, pf.endNode, pf.g2StartNode, false, 2)
	//
	//go pf.searchNodes(&wg, pf.startNode, pf.g3EndNode, true,3 )
	//go pf.searchNodes(&wg, pf.endNode, pf.g3StartNode, false, 3)
	//
	//go pf.searchNodes(&wg, pf.startNode, pf.g4EndNode, true, 4)
	//go pf.searchNodes(&wg, pf.endNode, pf.g4StartNode, false, 4)

	c = 0
	for _, nodeParent := range pf.WorldMap.GetNeighbours(pf.startNode) {
		for _, node := range pf.WorldMap.GetNeighbours(nodeParent) {
			//wg.Add(1)
			c++
			node.Parent = pf.startNode

			if c == 1 {
				go pf.searchNodes(&wg, node, pf.g1EndNode, true, 1)
			}
			if c == 2 {
				go pf.searchNodes(&wg, node, pf.g2EndNode, true, 2)
			}
			if c == 3 {
				go pf.searchNodes(&wg, node, pf.g3EndNode, true, 3)
			}
			if c == 4 {
				go pf.searchNodes(&wg, node, pf.g4EndNode, true, 4)
			}

			//go pf.searchNodes(&wg, node, pf.endNodeTmp, true)
		}
		break
	}

	c = 0

	for _, nodeParent := range pf.WorldMap.GetNeighbours(pf.endNode) {
		for _, node := range pf.WorldMap.GetNeighbours(nodeParent) {
			//wg.Add(1)
			c++
			node.Parent = pf.endNode

			if c == 1 {
				go pf.searchNodes(&wg, node, pf.g1StartNode, false, 1)
			}
			if c == 2 {
				go pf.searchNodes(&wg, node, pf.g2StartNode, false, 2)
			}
			if c == 3 {
				go pf.searchNodes(&wg, node, pf.g3StartNode, false, 3)
			}
			if c == 4 {
				go pf.searchNodes(&wg, node, pf.g4StartNode, false, 4)
			}
		}
		break
	}

	fmt.Println("Start", c, "threads")

	wg.Wait()

	return pf.endNode
}

func (pf *Pathfinderv7) updateVertex(currentNode, neighbour *internal.Node, weight float64, updateParent bool, heap *internal.NodesMinHeap, candidateWeights map[*internal.Node]float64, searchNode *internal.Node) {
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

func (pf *Pathfinderv7) checkAbilityToMove(parent, currentNode, neighbour *internal.Node) (float64, bool) {
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

func (pf *Pathfinderv7) hasCandidate(node *internal.Node) bool {
	_, ok := pf.candidateWeights[node]
	return ok
}

func (pf *Pathfinderv7) updateWeight(node *internal.Node, weight float64) {
	pf.candidateWeights[node] = weight
}

func (pf *Pathfinderv7) removeCandidate(node *internal.Node) {
	delete(pf.candidateWeights, node)
}

func (pf *Pathfinderv7) isVisible(currentNode, neighbour *internal.Node) bool {

	//return pf.lineOfSight.lineOfSight(currentNode, neighbour, true, true)

	lat1, lon1 := pf.grid.LatLon(currentNode)
	lat2, lon2 := pf.grid.LatLon(neighbour)
	return fillGreatCircle(lat1, lon1, lat2, lon2, pf.grid)
}
