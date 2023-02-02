package internal

const INFINITY = 999999.9

type Map struct {
	nodes [Width][height]*Node
	Grid  *Grid
}

func (m *Map) ExistNode(latIdx, lonIdx int) bool {
	return m.nodes[lonIdx][latIdx] != nil
}

func (m *Map) GetNode(latIdx, lonIdx int) *Node {
	lonIdx = normalizeLon(lonIdx)

	if m.nodes[lonIdx][latIdx] == nil {
		m.nodes[lonIdx][latIdx] = NewNode(latIdx, lonIdx)
	}

	return m.nodes[lonIdx][latIdx]
}
func (m *Map) GetNeighbours(n *Node) []*Node {
	neighbours := make([]*Node, 0)

	if n.LatIdx > 0 && m.Grid.IsTraversable(n.LonIdx, n.LatIdx-1) {
		neighbours = append(neighbours, m.GetNode(n.LatIdx-1, n.LonIdx))
	}

	if n.LatIdx < height-1 && m.Grid.IsTraversable(n.LonIdx, n.LatIdx+1) {
		neighbours = append(neighbours, m.GetNode(n.LatIdx+1, n.LonIdx))
	}

	if n.LonIdx > 0 && m.Grid.IsTraversable(n.LonIdx-1, n.LatIdx) {
		neighbours = append(neighbours, m.GetNode(n.LatIdx, n.LonIdx-1))
	}

	if n.LonIdx < Width-1 && m.Grid.IsTraversable(n.LonIdx+1, n.LatIdx) {
		neighbours = append(neighbours, m.GetNode(n.LatIdx, n.LonIdx+1))
	}

	if n.LatIdx > 0 && n.LonIdx > 0 && m.Grid.IsTraversable(n.LonIdx-1, n.LatIdx-1) {
		neighbours = append(neighbours, m.GetNode(n.LatIdx-1, n.LonIdx-1))
	}

	if n.LatIdx < height-1 && n.LonIdx < Width-1 && m.Grid.IsTraversable(n.LonIdx+1, n.LatIdx+1) {
		neighbours = append(neighbours, m.GetNode(n.LatIdx+1, n.LonIdx+1))
	}

	if n.LatIdx > 0 && n.LonIdx < Width-1 && m.Grid.IsTraversable(n.LonIdx+1, n.LatIdx-1) {
		neighbours = append(neighbours, m.GetNode(n.LatIdx-1, n.LonIdx+1))
	}

	if n.LatIdx < height-1 && n.LonIdx > 0 && m.Grid.IsTraversable(n.LonIdx-1, n.LatIdx+1) {
		neighbours = append(neighbours, m.GetNode(n.LatIdx+1, n.LonIdx-1))
	}

	return neighbours
}

func (m *Map) GetNeighboursAdopted(neighbours *[]*Node, n *Node) {
	if n.LatIdx > 0 && m.Grid.IsTraversable(n.LonIdx, n.LatIdx-1) {
		n := m.GetNode(n.LatIdx-1, n.LonIdx)
		if !n.Visited {
			*neighbours = append(*neighbours, n)
		}
	}

	if n.LatIdx < height-1 && m.Grid.IsTraversable(n.LonIdx, n.LatIdx+1) {
		n := m.GetNode(n.LatIdx+1, n.LonIdx)
		if !n.Visited {
			*neighbours = append(*neighbours, n)
		}
	}

	if n.LonIdx > 0 && m.Grid.IsTraversable(n.LonIdx-1, n.LatIdx) {
		n := m.GetNode(n.LatIdx, n.LonIdx-1)
		if !n.Visited {
			*neighbours = append(*neighbours, n)
		}
	}

	if n.LonIdx < Width-1 && m.Grid.IsTraversable(n.LonIdx+1, n.LatIdx) {
		n := m.GetNode(n.LatIdx, n.LonIdx+1)
		if !n.Visited {
			*neighbours = append(*neighbours, n)
		}
	}

	if n.LatIdx > 0 && n.LonIdx > 0 && m.Grid.IsTraversable(n.LonIdx-1, n.LatIdx-1) {
		n := m.GetNode(n.LatIdx-1, n.LonIdx-1)
		if !n.Visited {
			*neighbours = append(*neighbours, n)
		}
	}

	if n.LatIdx < height-1 && n.LonIdx < Width-1 && m.Grid.IsTraversable(n.LonIdx+1, n.LatIdx+1) {
		n := m.GetNode(n.LatIdx+1, n.LonIdx+1)
		if !n.Visited {
			*neighbours = append(*neighbours, n)
		}
	}

	if n.LatIdx > 0 && n.LonIdx < Width-1 && m.Grid.IsTraversable(n.LonIdx+1, n.LatIdx-1) {
		n := m.GetNode(n.LatIdx-1, n.LonIdx+1)
		if !n.Visited {
			*neighbours = append(*neighbours, n)
		}
	}

	if n.LatIdx < height-1 && n.LonIdx > 0 && m.Grid.IsTraversable(n.LonIdx-1, n.LatIdx+1) {
		n := m.GetNode(n.LatIdx+1, n.LonIdx-1)
		if !n.Visited {
			*neighbours = append(*neighbours, n)
		}
	}
}

func (m *Map) GetMinNeighbour(neighbours *[]*Node, n *Node) {
	var minNode *Node

	if n.LatIdx > 0 && m.Grid.IsTraversable(n.LonIdx, n.LatIdx-1) {
		n := m.GetNode(n.LatIdx-1, n.LonIdx)
		if !n.Visited {

			if minNode == nil {
				minNode = n
			}

		}
	}

	if n.LatIdx < height-1 && m.Grid.IsTraversable(n.LonIdx, n.LatIdx+1) {
		n := m.GetNode(n.LatIdx+1, n.LonIdx)
		if !n.Visited {

			if minNode == nil || minNode.GetScore() > n.GetScore() {
				minNode = n
			}
			//*neighbours = append(*neighbours, n)
		}
	}

	if n.LonIdx > 0 && m.Grid.IsTraversable(n.LonIdx-1, n.LatIdx) {
		n := m.GetNode(n.LatIdx, n.LonIdx-1)
		if !n.Visited {
			//*neighbours = append(*neighbours, n)

			if minNode == nil || minNode.GetScore() > n.GetScore() {
				minNode = n
			}
		}
	}

	if n.LonIdx < Width-1 && m.Grid.IsTraversable(n.LonIdx+1, n.LatIdx) {
		n := m.GetNode(n.LatIdx, n.LonIdx+1)
		if !n.Visited {
			//*neighbours = append(*neighbours, n)

			if minNode == nil || minNode.GetScore() > n.GetScore() {
				minNode = n
			}
		}
	}

	if n.LatIdx > 0 && n.LonIdx > 0 && m.Grid.IsTraversable(n.LonIdx-1, n.LatIdx-1) {
		n := m.GetNode(n.LatIdx-1, n.LonIdx-1)
		if !n.Visited {
			//*neighbours = append(*neighbours, n)

			if minNode == nil || minNode.GetScore() > n.GetScore() {
				minNode = n
			}
		}
	}

	if n.LatIdx < height-1 && n.LonIdx < Width-1 && m.Grid.IsTraversable(n.LonIdx+1, n.LatIdx+1) {
		n := m.GetNode(n.LatIdx+1, n.LonIdx+1)
		if !n.Visited {
			//*neighbours = append(*neighbours, n)

			if minNode == nil || minNode.GetScore() > n.GetScore() {
				minNode = n
			}
		}
	}

	if n.LatIdx > 0 && n.LonIdx < Width-1 && m.Grid.IsTraversable(n.LonIdx+1, n.LatIdx-1) {
		n := m.GetNode(n.LatIdx-1, n.LonIdx+1)
		if !n.Visited {
			//*neighbours = append(*neighbours, n)
			if minNode == nil || minNode.GetScore() > n.GetScore() {
				minNode = n
			}
		}
	}

	if n.LatIdx < height-1 && n.LonIdx > 0 && m.Grid.IsTraversable(n.LonIdx-1, n.LatIdx+1) {
		n := m.GetNode(n.LatIdx+1, n.LonIdx-1)
		if !n.Visited {
			//*neighbours = append(*neighbours, n)
			if minNode == nil || minNode.GetScore() > n.GetScore() {
				minNode = n
			}
		}
	}

	if minNode == nil {
		m.GetMinNeighbour(neighbours, n.Parent)
	} else {
		*neighbours = append(*neighbours, n)
	}

}

func (m *Map) GetNeighboursAdoptedFarAway(neighbours *[]*Node, n *Node) {
	dl := 10

	for i := dl; i > 0; i-- {
		lat := n.LatIdx - i
		lon := n.LonIdx

		if lat > 0 && m.Grid.IsTraversable(lon, lat) {
			n := m.GetNode(lat, lon)
			if !n.Visited {
				*neighbours = append(*neighbours, n)
				break
			}
		}
	}

	for i := dl; i > 0; i-- {
		lat := n.LatIdx + i
		lon := n.LonIdx

		if lat < height-1 && m.Grid.IsTraversable(lon, lat) {
			n := m.GetNode(lat, lon)
			if !n.Visited {
				*neighbours = append(*neighbours, n)
				break
			}
		}
	}

	for i := dl; i > 0; i-- {
		lat := n.LatIdx
		lon := n.LonIdx - i

		if lon > 0 && m.Grid.IsTraversable(lon, lat) {
			n := m.GetNode(lat, lon)
			if !n.Visited {
				*neighbours = append(*neighbours, n)
				break
			}
		}
	}

	for i := dl; i > 0; i-- {
		lat := n.LatIdx
		lon := n.LonIdx + i

		if lon < Width-1 && m.Grid.IsTraversable(lon, lat) {
			n := m.GetNode(lat, lon)
			if !n.Visited {
				*neighbours = append(*neighbours, n)
				break
			}
		}
	}

	for i := dl; i > 0; i-- {
		lat := n.LatIdx - i
		lon := n.LonIdx - i

		if lon > 0 && lat > 0 && m.Grid.IsTraversable(lon, lat) {
			n := m.GetNode(lat, lon)
			if !n.Visited {
				*neighbours = append(*neighbours, n)
				break
			}
		}
	}

	for i := dl; i > 0; i-- {
		lat := n.LatIdx + i
		lon := n.LonIdx + i

		if lon < Width-1 && lat < height-1 && m.Grid.IsTraversable(lon, lat) {
			n := m.GetNode(lat, lon)
			if !n.Visited {
				*neighbours = append(*neighbours, n)
				break
			}
		}
	}

	for i := dl; i > 0; i-- {
		lat := n.LatIdx - i
		lon := n.LonIdx + i

		if lon < Width-1 && lat > 0 && m.Grid.IsTraversable(lon, lat) {
			n := m.GetNode(lat, lon)
			if !n.Visited {
				*neighbours = append(*neighbours, n)
				break
			}
		}
	}

	for i := dl; i > 0; i-- {
		lat := n.LatIdx + i
		lon := n.LonIdx - i

		if lon > 0 && lat < height-1 && m.Grid.IsTraversable(lon, lat) {
			n := m.GetNode(lat, lon)
			if !n.Visited {
				*neighbours = append(*neighbours, n)
				break
			}
		}
	}

	for i := -dl + 1; i < dl; i++ {
		for j := -dl + 1; j < dl; j++ {
			lat := n.LatIdx + i
			lon := n.LonIdx + j

			if lon > 0 && lat > 0 && lat < height-1 && lon < Width-1 && m.Grid.IsTraversable(lon, lat) {
				m.GetNode(lat, lon).Visited = true
			}
		}
	}

}

func (m *Map) InsideMap(lat, lon int) bool {
	lon = normalizeLon(lon)
	if lon > 0 && lat > 0 && lat < height-1 && lon < Width-1 {
		return true
	}

	return false
}
func (m *Map) Valid(lat, lon int) bool {
	if m.InsideMap(lat, lon) && m.Grid.IsTraversable(lon, lat) {
		return true
	}

	return false
}

func (m *Map) GetNeighboursAdoptedFarAwayIncreasing(neighbours *[]*Node, n *Node) {
	dl := 10

	n.Current = true

	signs := [][2]int{
		{-1, 0},
		{1, 0},
		{0, -1},
		{0, 1},
		{1, 1},
		{1, -1},
		{-1, 1},
		{-1, -1},
	}

	for _, sign := range signs {
		prevLat, prevLon := 0, 0

		for i := 1; i <= dl; i++ {
			lat := n.LatIdx + i*sign[0]
			lon := n.LonIdx + i*sign[1]

			if m.Valid(lat, lon) && !m.GetNode(lat, lon).Visited {
				prevLat, prevLon = lat, lon

				if i != dl {
					continue
				}
			}

			if prevLat != 0 && prevLon != 0 {
				c := 0

				sLon := intSign(prevLon, n.LonIdx)
				sLat := intSign(prevLat, n.LatIdx)
				dLon := Abs(prevLon - n.LonIdx)
				dLat := Abs(prevLat - n.LatIdx)

				if dLon > 0 {
					dLon--
				}
				if dLat > 0 {
					dLat--
				}

				for iLon := 0; iLon <= dLon; iLon += 1 {
					for iLat := 0; iLat <= dLat; iLat += 1 {
						visitedNode := m.GetNode(n.LatIdx+iLat*sLat, n.LonIdx+iLon*sLon)
						if !visitedNode.InNeibers {
							visitedNode.Visited = true
						}
						c++
					}
				}

				n2 := m.GetNode(prevLat, prevLon)
				n2.Visited = false
				n2.InNeibers = true
				*neighbours = append(*neighbours, n2)
			}

			break
		}
	}

}

func (m *Map) GetNeighboursAdoptedFarAwayIncreasing2(neighbours *[]*Node, n *Node) {
	dl := 10

	n.Current = true

	signs := [][2]int{
		{-1, 0},
		{1, 0},
		{0, -1},
		{0, 1},
		{1, 1},
		{1, -1},
		{-1, 1},
		{-1, -1},
	}

	for _, sign := range signs {
		lat := n.LatIdx + dl*sign[0]
		lon := n.LonIdx + dl*sign[1]

		if m.Valid(lat, lon) && !m.GetNode(lat, lon).Visited {
			for iLon := 0; iLon < dl; iLon += 1 {
				for iLat := 0; iLat < dl; iLat += 1 {
					n2 := m.GetNode(n.LatIdx+iLat*sign[0], n.LonIdx+iLon*sign[1])
					if !n2.InNeibers {
						n2.Visited = true
					}
				}
			}

			n2 := m.GetNode(lat, lon)
			n2.Visited = false
			n2.InNeibers = true
			*neighbours = append(*neighbours, n2)
		} else {
			lat = n.LatIdx + 1*sign[0]
			lon = n.LonIdx + 1*sign[1]
			if m.Valid(lat, lon) && !m.GetNode(lat, lon).Visited {
				n2 := m.GetNode(lat, lon)
				n2.InNeibers = true
				*neighbours = append(*neighbours, n2)
			}
		}
	}
}

func normalizeLon(lonIdx int) int {
	if lonIdx < 0 {
		lonIdx = Width - lonIdx
	}
	if lonIdx >= Width {
		lonIdx = lonIdx - Width
	}
	return lonIdx
}

func (m *Map) GetNeighboursAdoptedFarAwayIncreasing3(neighbours *[]*Node, n *Node) {
	dl := 100

	n.Current = true

	closeToLand := 0

	signs := [][2]int{
		{-1, 0},
		{1, 0},
		{0, -1},
		{0, 1},
		{1, 1},
		{1, -1},
		{-1, 1},
		{-1, -1},
	}

	for _, sign := range signs {
		lat := n.LatIdx + 3*sign[0]
		lon := n.LonIdx + 3*sign[1]
		lon = normalizeLon(lon)

		if m.InsideMap(lat, lon) && !m.Grid.IsTraversable(lon, lat) {
			closeToLand++
		}
	}

	if closeToLand > 3 {
		dl = 1
	}

	for _, sign := range signs {
		prevLat, prevLon := 0, 0

		for i := 1; i <= dl; i++ {
			lat := n.LatIdx + i*sign[0]
			lon := n.LonIdx + i*sign[1]
			lon = normalizeLon(lon)

			if i == 1 && m.Valid(lat, lon) {
				nextLat := n.LatIdx + (i+1)*sign[0]
				nextLon := n.LonIdx + (i+1)*sign[1]

				if m.InsideMap(nextLat, nextLon) && !m.Grid.IsTraversable(nextLon, nextLat) {
					n2 := m.GetNode(lat, lon)
					if !n2.Current {
						n2.Visited = false
					}
					n2.InNeibers = true
					*neighbours = append(*neighbours, n2)
					break
				}
			}

			if m.Valid(lat, lon) && !m.GetNode(lat, lon).Visited {
				prevLat, prevLon = lat, lon

				if i != dl {
					continue
				}
			}

			if prevLat != 0 && prevLon != 0 {
				c := 0

				sLon := intSign(prevLon, n.LonIdx)
				sLat := intSign(prevLat, n.LatIdx)
				dLon := Abs(prevLon - n.LonIdx)
				dLat := Abs(prevLat - n.LatIdx)

				if dLon > 0 {
					dLon--
				}
				if dLat > 0 {
					dLat--
				}

				for iLon := 0; iLon <= dLon; iLon += 1 {
					for iLat := 0; iLat <= dLat; iLat += 1 {
						visitedNode := m.GetNode(n.LatIdx+iLat*sLat, n.LonIdx+iLon*sLon)
						if !visitedNode.InNeibers {
							visitedNode.Visited = true
						}
						c++
					}
				}

				n2 := m.GetNode(prevLat, prevLon)
				n2.Visited = false
				n2.InNeibers = true

				*neighbours = append(*neighbours, n2)
			}

			break
		}
	}

}

func (m *Map) GetNeighboursAdoptedFarAwayIncreasing4ForSmallMap(neighbours *[]*Node, n *Node) {
	dl := 25

	signs := [][2]int{
		{-1, 0},
		{1, 0},
		{0, -1},
		{0, 1},
		{1, 1},
		{1, -1},
		{-1, 1},
		{-1, -1},
	}

	lastI := 1

	for i := 1; i <= dl; i++ {
		allTraversable := true
		for _, sign := range signs {
			lat := n.LatIdx + i*sign[0]
			lon := n.LonIdx + i*sign[1]
			if m.Valid(lat, lon) {
				//*neighbours = append(*neighbours, m.getNode(lat, lon))
			} else {
				allTraversable = false
				break
			}
		}

		if !allTraversable {
			break
		}

		lastI = i
	}

	for ix := -lastI + 1; ix < lastI; ix++ {
		for iy := -lastI + 1; iy < lastI; iy++ {
			if m.Valid(n.LatIdx+iy, n.LonIdx+ix) {
				m.GetNode(n.LatIdx+iy, n.LonIdx+ix).Visited = true
			}
		}
	}

	for latI := -lastI; latI <= lastI; latI++ {
		lat := n.LatIdx + latI
		lon := n.LonIdx - lastI
		n2 := m.GetNode(lat, lon)
		if !n2.Visited && m.Valid(lat, lon) {
			*neighbours = append(*neighbours, n2)
		}
	}
	for latI := -lastI; latI <= lastI; latI++ {
		lat := n.LatIdx + latI
		lon := n.LonIdx + lastI
		n2 := m.GetNode(lat, lon)
		if !n2.Visited && m.Valid(lat, lon) {
			*neighbours = append(*neighbours, n2)
		}
	}

	for lonI := -lastI; lonI <= lastI; lonI++ {
		lat := n.LatIdx + lastI
		lon := n.LonIdx + lonI
		n2 := m.GetNode(lat, lon)
		if !n2.Visited && m.Valid(lat, lon) {
			*neighbours = append(*neighbours, n2)
		}
	}
	for lonI := -lastI; lonI <= lastI; lonI++ {
		lat := n.LatIdx - lastI
		lon := n.LonIdx + lonI
		n2 := m.GetNode(lat, lon)
		if !n2.Visited && m.Valid(lat, lon) {
			*neighbours = append(*neighbours, n2)
		}
	}

}

func (m *Map) GetNeighboursAdoptedFarAwayIncreasing4(neighbours *[]*Node, n *Node, dl int) {
	n.Current = true

	signs := [][2]int{
		{-1, 0},
		{1, 0},
		{0, -1},
		{0, 1},
		{1, 1},
		{1, -1},
		{-1, 1},
		{-1, -1},
	}

	maxDl := dl

	if dl > 4 {
		for _, sign := range signs {
			for i := 1; i <= dl; i++ {
				lat := n.LatIdx + i*sign[0]
				lon := n.LonIdx + i*sign[1]

				if !m.Valid(lon, lat) && i < maxDl {
					maxDl = i
					break
				}
			}
		}
	} else {
		maxDl = 1
	}

	for _, sign := range signs {
		lat := n.LatIdx + maxDl*sign[0]
		lon := n.LonIdx + maxDl*sign[1]

		n2 := m.GetNode(lat, lon)
		n2.Visited = false
		n2.InNeibers = true
		*neighbours = append(*neighbours, n2)
	}

	for i := 1; i < maxDl; i++ {
		for j := 1; j < maxDl; j++ {
			visitedNode := m.GetNode(n.LatIdx+i, n.LonIdx+j)
			if !visitedNode.InNeibers {
				visitedNode.Visited = true
			}
		}
	}

}

func (m *Map) GetNeighboursAdoptedFarAwayIncreasing5(neighbours *[]*Node, n *Node) {
	dl := 30

	n.Current = true

	signs := [][2]int{
		{-1, 0},
		{1, 0},
		{0, -1},
		{0, 1},
		{1, 1},
		{1, -1},
		{-1, 1},
		{-1, -1},
	}

	maxDl := dl

	if dl > 4 {
		for _, sign := range signs {
			for i := 1; i <= dl; i++ {
				lat := n.LatIdx + i*sign[0]
				lon := n.LonIdx + i*sign[1]

				if !m.Valid(lon, lat) && i < maxDl {
					maxDl = i
					break
				}
			}
		}
	} else {
		maxDl = 1
	}

	for _, sign := range signs {
		lat := n.LatIdx + maxDl*sign[0]
		lon := n.LonIdx + maxDl*sign[1]

		n2 := m.GetNode(lat, lon)
		n2.Visited = false
		n2.InNeibers = true
		*neighbours = append(*neighbours, n2)
	}

	for i := 1; i < maxDl; i++ {
		for j := 1; j < maxDl; j++ {
			visitedNode := m.GetNode(n.LatIdx+i, n.LonIdx+j)
			if !visitedNode.InNeibers {
				visitedNode.Visited = true
			}
		}
	}

}

func intSign(x1, x2 int) int {
	if x1 >= x2 {
		return 1
	}
	return -1
}

func (m *Map) GetNeighboursAdoptedFarAwayIncreasingOLD(neighbours *[]*Node, n *Node) {
	dl := 25
	prevLat, prevLon := 0, 0

	signs := [][2]int8{
		{-1, 0},
		{1, 0},
		{0, -1},
		{0, 1},
		{1, 1},
		{1, -1},
		{-1, 1},
		{-1, -1},
	}

	for _, sign := range signs {
		sign[0]++
		for i := 1; i < dl; i++ {
			lat := n.LatIdx - i
			lon := n.LonIdx

			if lat > 0 && m.Grid.IsTraversable(lon, lat) && !m.GetNode(lat, lon).Visited {
				prevLat, prevLon = lat, lon
				continue
			}

			*neighbours = append(*neighbours, m.GetNode(prevLat, prevLon))
			break
		}
	}

	for i := 1; i < dl; i++ {
		lat := n.LatIdx - i
		lon := n.LonIdx

		if lat > 0 && m.Grid.IsTraversable(lon, lat) && !m.GetNode(lat, lon).Visited {
			prevLat, prevLon = lat, lon
			continue
		}

		*neighbours = append(*neighbours, m.GetNode(prevLat, prevLon))
		break
	}

	for i := 1; i < dl; i++ {
		lat := n.LatIdx + i
		lon := n.LonIdx

		if lat < height-1 && m.Grid.IsTraversable(lon, lat) && !m.GetNode(lat, lon).Visited {
			prevLat, prevLon = lat, lon
			continue
		}

		*neighbours = append(*neighbours, m.GetNode(prevLat, prevLon))
		break
	}

	for i := 0; i < dl; i++ {
		lat := n.LatIdx
		lon := n.LonIdx - i

		if lon > 0 && m.Grid.IsTraversable(lon, lat) && !m.GetNode(lat, lon).Visited {
			prevLat, prevLon = lat, lon
			continue
		}

		*neighbours = append(*neighbours, m.GetNode(prevLat, prevLon))
		break
	}

	for i := 0; i < dl; i++ {
		lat := n.LatIdx
		lon := n.LonIdx + i

		if lon < Width-1 && m.Grid.IsTraversable(lon, lat) && !m.GetNode(lat, lon).Visited {
			prevLat, prevLon = lat, lon
			continue
		}

		*neighbours = append(*neighbours, m.GetNode(prevLat, prevLon))
		break
	}

	for i := 0; i < dl; i++ {
		lat := n.LatIdx - i
		lon := n.LonIdx - i

		if lon > 0 && lat > 0 && m.Grid.IsTraversable(lon, lat) && !m.GetNode(lat, lon).Visited {
			prevLat, prevLon = lat, lon
			continue
		}

		*neighbours = append(*neighbours, m.GetNode(prevLat, prevLon))
		break
	}

	for i := 0; i < dl; i++ {
		lat := n.LatIdx + i
		lon := n.LonIdx + i

		if lon < Width-1 && lat < height-1 && m.Grid.IsTraversable(lon, lat) && !m.GetNode(lat, lon).Visited {
			prevLat, prevLon = lat, lon
			continue
		}

		*neighbours = append(*neighbours, m.GetNode(prevLat, prevLon))
		break
	}

	for i := 0; i < dl; i++ {
		lat := n.LatIdx - i
		lon := n.LonIdx + i

		if lon < Width-1 && lat > 0 && m.Grid.IsTraversable(lon, lat) && !m.GetNode(lat, lon).Visited {
			prevLat, prevLon = lat, lon
			continue
		}

		*neighbours = append(*neighbours, m.GetNode(prevLat, prevLon))
		break
	}

	for i := 0; i < dl; i++ {
		lat := n.LatIdx + i
		lon := n.LonIdx - i

		if lon > 0 && lat < height-1 && m.Grid.IsTraversable(lon, lat) && !m.GetNode(lat, lon).Visited {
			prevLat, prevLon = lat, lon
			continue
		}

		*neighbours = append(*neighbours, m.GetNode(prevLat, prevLon))
		break
	}

	for i := -dl + 1; i < dl; i++ {
		for j := -dl + 1; j < dl; j++ {
			lat := n.LatIdx + i
			lon := n.LonIdx + j

			if lon > 0 && lat > 0 && lat < height-1 && lon < Width-1 && m.Grid.IsTraversable(lon, lat) {
				m.GetNode(lat, lon).Visited = true
			}
		}
	}

}

func (m *Map) GetNeighboursAdoptedV2(neighbours *[]*Node, n *Node) {
	dl := 100

	iLat := dl
	for j := -dl; j <= dl; j++ {
		if n.LatIdx+iLat >= 0 && n.LatIdx+iLat < height && n.LonIdx+j >= 0 && n.LonIdx+j < Width {
			if m.Grid.IsTraversable(n.LonIdx+j, n.LatIdx+iLat) {
				*neighbours = append(*neighbours, m.GetNode(n.LatIdx+iLat, n.LonIdx+j))
			}
		}
	}

	iLat = -dl
	for j := -dl; j <= dl; j++ {
		if n.LatIdx+iLat >= 0 && n.LatIdx+iLat < height && n.LonIdx+j >= 0 && n.LonIdx+j < Width {
			if m.Grid.IsTraversable(n.LonIdx+j, n.LatIdx+iLat) {
				*neighbours = append(*neighbours, m.GetNode(n.LatIdx+iLat, n.LonIdx+j))
			}
		}
	}

	iLon := -dl
	for i := -dl - 1; i <= dl-1; i++ {
		if n.LatIdx+i >= 0 && n.LatIdx+i < height && n.LonIdx+iLon >= 0 && n.LonIdx+iLon < Width {
			if m.Grid.IsTraversable(n.LonIdx+iLon, n.LatIdx+i) {
				*neighbours = append(*neighbours, m.GetNode(n.LatIdx+i, n.LonIdx+iLon))
			}
		}
	}

	iLon = dl
	for i := -dl - 1; i <= dl-1; i++ {
		if n.LatIdx+i >= 0 && n.LatIdx+i < height && n.LonIdx+iLon >= 0 && n.LonIdx+iLon < Width {
			if m.Grid.IsTraversable(n.LonIdx+iLon, n.LatIdx+i) {
				*neighbours = append(*neighbours, m.GetNode(n.LatIdx+i, n.LonIdx+iLon))
			}
		}
	}

	for ix := -dl + 1; ix < dl; ix++ {
		for iy := -dl + 1; iy < dl; iy++ {
			m.GetNode(n.LatIdx+ix, n.LonIdx+iy).Visited = true
		}
	}

	return

	if n.LatIdx > 0 && m.Grid.IsTraversable(n.LonIdx, n.LatIdx-1) {
		*neighbours = append(*neighbours, m.GetNode(n.LatIdx-1, n.LonIdx))
	}

	if n.LatIdx < height-1 && m.Grid.IsTraversable(n.LonIdx, n.LatIdx+1) {
		*neighbours = append(*neighbours, m.GetNode(n.LatIdx+1, n.LonIdx))
	}

	if n.LonIdx > 0 && m.Grid.IsTraversable(n.LonIdx-1, n.LatIdx) {
		*neighbours = append(*neighbours, m.GetNode(n.LatIdx, n.LonIdx-1))
	}

	if n.LonIdx < Width-1 && m.Grid.IsTraversable(n.LonIdx+1, n.LatIdx) {
		*neighbours = append(*neighbours, m.GetNode(n.LatIdx, n.LonIdx+1))
	}

	if n.LatIdx > 0 && n.LonIdx > 0 && m.Grid.IsTraversable(n.LonIdx-1, n.LatIdx-1) {
		*neighbours = append(*neighbours, m.GetNode(n.LatIdx-1, n.LonIdx-1))
	}

	if n.LatIdx < height-1 && n.LonIdx < Width-1 && m.Grid.IsTraversable(n.LonIdx+1, n.LatIdx+1) {
		*neighbours = append(*neighbours, m.GetNode(n.LatIdx+1, n.LonIdx+1))
	}

	if n.LatIdx > 0 && n.LonIdx < Width-1 && m.Grid.IsTraversable(n.LonIdx+1, n.LatIdx-1) {
		*neighbours = append(*neighbours, m.GetNode(n.LatIdx-1, n.LonIdx+1))
	}

	if n.LatIdx < height-1 && n.LonIdx > 0 && m.Grid.IsTraversable(n.LonIdx-1, n.LatIdx+1) {
		*neighbours = append(*neighbours, m.GetNode(n.LatIdx+1, n.LonIdx-1))
	}
}

type Node struct {
	FScore    float64
	GScore    float64
	LatIdx    int
	LonIdx    int
	Parent    *Node
	_         int
	Visited   bool
	Current   bool
	InNeibers bool
	_         bool
}

func (n *Node) GetScore() float64 {
	return n.GScore // - n.FScore/5
}

func NewNode(latIdx, lonIdx int) *Node {
	n := Node{
		FScore:    INFINITY,
		GScore:    INFINITY,
		LatIdx:    latIdx,
		LonIdx:    lonIdx,
		Parent:    nil,
		Visited:   false,
		Current:   false,
		InNeibers: false,
	}
	return &n
}

func (n *Node) GetNeighbours() []*Node {
	neighbours := make([]*Node, 0)

	if n.LatIdx > 0 {
		neighbours = append(neighbours, NewNode(n.LatIdx-1, n.LonIdx))
	}

	if n.LatIdx < height-1 {
		neighbours = append(neighbours, NewNode(n.LatIdx+1, n.LonIdx))
	}

	if n.LonIdx > 0 {
		neighbours = append(neighbours, NewNode(n.LatIdx, n.LonIdx-1))
	}

	if n.LonIdx < Width-1 {
		neighbours = append(neighbours, NewNode(n.LatIdx, n.LonIdx+1))
	}

	if n.LatIdx > 0 && n.LonIdx > 0 {
		neighbours = append(neighbours, NewNode(n.LatIdx-1, n.LonIdx-1))
	}

	if n.LatIdx < height-1 && n.LonIdx < Width-1 {
		neighbours = append(neighbours, NewNode(n.LatIdx+1, n.LonIdx+1))
	}

	if n.LatIdx > 0 && n.LonIdx < Width-1 {
		neighbours = append(neighbours, NewNode(n.LatIdx-1, n.LonIdx+1))
	}

	if n.LatIdx < height-1 && n.LonIdx > 0 {
		neighbours = append(neighbours, NewNode(n.LatIdx+1, n.LonIdx-1))
	}

	return neighbours
}
