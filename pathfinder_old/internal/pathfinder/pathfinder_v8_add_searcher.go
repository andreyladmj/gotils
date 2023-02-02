package pathfinder

import (
	"fmt"
	"pathfinder/internal"
	"sync"
)

type Pathfinderv8 struct {
	lineOfSight *internal.LineOfSightChecking

	StartNode *internal.Node
	EndNode   *internal.Node

	found bool
	Pathfinder
}

type searcher struct {
	candidateWeights map[*internal.Node]float64

	startNode   *internal.Node
	searchNode  *internal.Node
	CurrentNode *internal.Node

	anotherSearchers []*searcher

	heap  *internal.NodesMinHeap
	grid  *internal.Grid
	found bool
}

func NewPathfinderv8(grid *internal.Grid) *Pathfinderv8 {
	return &Pathfinderv8{
		lineOfSight: nil,
		StartNode:   nil,
		EndNode:     nil,
		Pathfinder:  Pathfinder{grid: grid, WorldMap: &internal.Map{Grid: grid}},
	}
}
func NewSearcher(grid *internal.Grid, searchNode *internal.Node) *searcher {
	return &searcher{
		candidateWeights: nil,
		searchNode:       searchNode,
		heap:             nil,
		grid:             grid,
	}
}

func (s *searcher) searchNodes(wg *sync.WaitGroup, currentNode *internal.Node, pf *Pathfinderv8) {
	s.heap = internal.NewHeap(10000)
	s.heap.Clear()
	s.heap.Insert(currentNode)
	stop := false
	s.candidateWeights = map[*internal.Node]float64{}
	defer wg.Done()
	c := 0

	for !s.heap.Empty() {
		currentNode = s.heap.GetMin()
		s.CurrentNode = currentNode

		if currentNode.Visited {
			continue
		}

		currentNode.Visited = true
		parent := currentNode.Parent
		neighbours := make([]*internal.Node, 0, 8)
		//pf.WorldMap.GetNeighboursAdoptedFarAwayIncreasing(&neighbours, currentNode)
		pf.WorldMap.GetNeighboursAdoptedFarAwayIncreasing4ForSmallMap(&neighbours, currentNode)

		c++
		if c > 100000 {
			fmt.Println("EXCEED LENGTH")
			return
		}

		for _, neighbour := range neighbours {

			if pf.found {
				fmt.Println("Found")
				return
			}

			if neighbour.Parent != nil && neighbour.Parent == parent {

			} else if hasCandidate(s.candidateWeights, neighbour) {
				candidateWeight := s.candidateWeights[neighbour]

				if candidateWeight != internal.INFINITY {
					newWeight, updateParent := s.checkAbilityToMove(parent, currentNode, neighbour)

					if candidateWeight > newWeight {
						s.updateVertex(currentNode, neighbour, newWeight, updateParent)
					}
				}
			} else {
				if s.isVisible(currentNode, neighbour) {
					newWeight, updateParent := s.checkAbilityToMove(parent, currentNode, neighbour)
					s.updateVertex(currentNode, neighbour, newWeight, updateParent)
				}
			}

			if s.isVisible(neighbour, s.searchNode) {
				s.searchNode.Parent = neighbour
				fmt.Println("Point searchNode", s.searchNode.LatIdx, s.searchNode.LonIdx)
				fmt.Println("Point neighbour", neighbour.LatIdx, neighbour.LonIdx)
				fmt.Println("Point currentNode", currentNode.LatIdx, currentNode.LonIdx)
				fmt.Println("FINISHED")
				s.found = true
				pf.found = true
				stop = true
				return
			}

			for _, s2 := range s.anotherSearchers {
				node := s2.CurrentNode
				for node != nil && node.Parent != node {
					if s.isVisible(neighbour, s2.CurrentNode) {
						fmt.Println("found another Searchers", s2.CurrentNode.LatIdx, s2.CurrentNode.LonIdx)
						fmt.Println("Point neighbour", neighbour.LatIdx, neighbour.LonIdx)
						fmt.Println("Point currentNode", currentNode.LatIdx, currentNode.LonIdx)
						fmt.Println("FINISHED")
						s.found = true
						pf.found = true
						stop = true
						return
					}

					node = node.Parent
				}
			}
		}

		if stop {
			break
		}
	}
}

func (s *searcher) addSearchers(searchers ...*searcher) {
	for _, searcher := range searchers {
		s.anotherSearchers = append(s.anotherSearchers, searcher)
	}
}

func (pf *Pathfinderv8) Find(startIndex, endIndex Index) (*searcher, *searcher) {
	pf.StartNode = pf.WorldMap.GetNode(startIndex.LatIdx, startIndex.LonIdx)
	pf.EndNode = pf.WorldMap.GetNode(endIndex.LatIdx, endIndex.LonIdx)
	pf.found = false

	fmt.Println("startNode", pf.StartNode.LatIdx, pf.StartNode.LonIdx)
	fmt.Println("endNode", pf.EndNode.LatIdx, pf.EndNode.LonIdx)

	pf.StartNode.GScore = pf.grid.Haversine(pf.StartNode, pf.EndNode)
	pf.StartNode.FScore = 0
	var wg sync.WaitGroup

	pf.EndNode.Parent = pf.EndNode
	pf.StartNode.Parent = pf.StartNode

	s1 := NewSearcher(pf.grid, pf.EndNode)
	s2 := NewSearcher(pf.grid, pf.StartNode)

	s3 := NewSearcher(pf.grid, pf.EndNode)
	s4 := NewSearcher(pf.grid, pf.StartNode)
	s5 := NewSearcher(pf.grid, pf.EndNode)
	s6 := NewSearcher(pf.grid, pf.StartNode)

	s1.addSearchers(s2, s4, s6)
	s2.addSearchers(s1, s3, s5)

	s3.addSearchers(s2, s4, s6)
	s4.addSearchers(s1, s3, s5)

	s5.addSearchers(s2, s4, s6)
	s6.addSearchers(s1, s3, s5)

	wg.Add(6)

	go s1.searchNodes(&wg, pf.StartNode, pf)
	go s2.searchNodes(&wg, pf.EndNode, pf)
	go s3.searchNodes(&wg, pf.StartNode, pf)
	go s4.searchNodes(&wg, pf.EndNode, pf)
	go s5.searchNodes(&wg, pf.StartNode, pf)
	go s6.searchNodes(&wg, pf.EndNode, pf)

	fmt.Println("Start")

	wg.Wait()

	return s1, s2
}

func (s *searcher) updateVertex(currentNode, neighbour *internal.Node, weight float64, updateParent bool) {
	if s.hasCandidate(neighbour) {
		s.removeCandidate(neighbour)
	}

	s.updateWeight(neighbour, weight)

	if updateParent {
		neighbour.Parent = currentNode
	} else {
		neighbour.Parent = currentNode.Parent
	}

	neighbour.FScore = weight
	neighbour.GScore = weight + s.grid.Haversine(neighbour, s.searchNode)

	s.heap.Insert(neighbour)
}

func (s *searcher) checkAbilityToMove(parent, currentNode, neighbour *internal.Node) (float64, bool) {
	var newWeight float64
	var updateParent bool

	if s.isVisible(parent, neighbour) {
		newWeight = parent.FScore + s.grid.Haversine(parent, neighbour)
		updateParent = false
	} else {
		newWeight = currentNode.FScore + s.grid.Haversine(currentNode, neighbour)
		updateParent = true
	}

	return newWeight, updateParent
}

func (s *searcher) hasCandidate(node *internal.Node) bool {
	_, ok := s.candidateWeights[node]
	return ok
}

func (s *searcher) updateWeight(node *internal.Node, weight float64) {
	s.candidateWeights[node] = weight
}

func (s *searcher) removeCandidate(node *internal.Node) {
	delete(s.candidateWeights, node)
}

func (s *searcher) isVisible(currentNode, neighbour *internal.Node) bool {
	//return false

	//return pf.lineOfSight.lineOfSight(currentNode, neighbour, true, true)

	lat1, lon1 := s.grid.LatLon(currentNode)
	lat2, lon2 := s.grid.LatLon(neighbour)
	return fillGreatCircle(lat1, lon1, lat2, lon2, s.grid)
}
