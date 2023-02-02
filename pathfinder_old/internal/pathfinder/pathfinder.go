package pathfinder

import (
	"fmt"
	"log"
	"pathfinder/internal"
)

type Pathfinder struct {
	grid     *internal.Grid
	WorldMap *internal.Map
}

type Index struct {
	LatIdx, LonIdx int
}

func NewPathfinder(grid *internal.Grid) *Pathfinder {
	return &Pathfinder{grid: grid, WorldMap: &internal.Map{Grid: grid}}
}

func (pf *Pathfinder) find(startIndex, endIndex Index) *internal.Node {
	heap := internal.NewHeap(25000)

	startNode := pf.WorldMap.GetNode(startIndex.LatIdx, startIndex.LonIdx)
	endNode := pf.WorldMap.GetNode(endIndex.LatIdx, endIndex.LonIdx)

	fmt.Println("startNode", startNode.LatIdx, startNode.LonIdx)
	fmt.Println("endNode", endNode.LatIdx, endNode.LonIdx)

	startNode.GScore = pf.grid.Haversine(startNode, endNode)
	startNode.FScore = 0
	heap.Clear()
	heap.Insert(startNode)

	var currentNode *internal.Node
	stop := false
	c := 0

	for !heap.Empty() {
		currentNode = heap.GetMin()

		if currentNode == nil {
			log.Fatalf("Path wasnt found")
		}

		if currentNode.Visited {
			continue
		}

		//fmt.Println("currentNode", currentNode.LatIdx, currentNode.LonIdx, "heap size", heap.heapSize, "visited", currentNode.Visited)
		currentNode.Visited = true

		for _, neighbour := range pf.WorldMap.GetNeighbours(currentNode) {
			c++
			fPossibleLowerGoal := currentNode.FScore + pf.grid.Haversine(currentNode, neighbour)

			if c > 10000000 {
				log.Fatalf("Pathfinding is working too long")
			}

			if fPossibleLowerGoal < neighbour.FScore {
				neighbour.Parent = currentNode
				neighbour.FScore = fPossibleLowerGoal

				neighbour.GScore = neighbour.FScore + 5*pf.grid.Haversine(neighbour, endNode)

				if neighbour.LatIdx == endNode.LatIdx && neighbour.LonIdx == endNode.LonIdx {
					endNode = neighbour
					stop = true
					break
				}
			}

			if !neighbour.Visited {
				heap.Insert(neighbour)
			}
		}

		if stop {
			break
		}

	}

	return endNode
}
