package internal

func (m *Map) isTraversableRect(node *Node, dl int) bool {
	for ix := -dl; ix <= dl; ix++ {
		for iy := -dl; iy <= dl; iy++ {
			if node.LonIdx+dl > 0 && node.LonIdx+dl < Width && node.LatIdx+dl > 0 && node.LatIdx+dl < height {
				if !m.Grid.IsTraversable(node.LonIdx+dl, node.LatIdx+dl) {
					return false
				}
			}
		}
	}

	return true
}

func (m *Map) checkPoint(lon, lat int) bool {
	if lon > 0 && lon < Width && lat > 0 && lat < height {
		if m.Grid.IsTraversable(lon, lat) && !m.GetNode(lat, lon).Visited {
			return true
		}
	}

	return false
}

func (m *Map) checkLine(lon1, lat1, lon2, lat2 int) bool {
	return TraversableSegmentOnPlane(lon1, lat1, lon2, lat2, m.Grid)
}

func (m *Map) addNonBlockedPoints(neighbours *[]*Node, lon1, lat1, lon2, lat2 int) {
	dx := Abs(lon1 - lon2)
	dy := Abs(lat1 - lat2)
	minLat := Min(lat1, lat2)
	minLon := Min(lon1, lon2)
	for iLon := 0; iLon <= dx; iLon += 1 {
		for iLat := 0; iLat <= dy; iLat += 1 {
			if iLat > 0 && iLat < height && iLon > 0 && iLon < Width {
				if m.Grid.IsTraversable(minLon+iLon, minLat+iLat) {
					*neighbours = append(*neighbours, m.GetNode(minLat+iLat, minLon+iLon))
				}
			}
		}
	}
}

func (m *Map) GetNeighboursAdoptedV4(neighbours *[]*Node, n *Node) {
	dl := 500

	p1Lon, p1Lat := n.LonIdx-dl, n.LatIdx
	p2Lon, p2Lat := n.LonIdx-dl, n.LatIdx+dl
	p3Lon, p3Lat := n.LonIdx, n.LatIdx+dl
	p4Lon, p4Lat := n.LonIdx+dl, n.LatIdx+dl
	p5Lon, p5Lat := n.LonIdx+dl, n.LatIdx
	p6Lon, p6Lat := n.LonIdx+dl, n.LatIdx-dl
	p7Lon, p7Lat := n.LonIdx, n.LatIdx-dl
	p8Lon, p8Lat := n.LonIdx-dl, n.LatIdx-dl

	if m.checkPoint(p1Lon, p1Lat) {
		*neighbours = append(*neighbours, m.GetNode(p1Lat, p1Lon))
	}

	if m.checkPoint(p2Lon, p2Lat) {
		*neighbours = append(*neighbours, m.GetNode(p2Lat, p2Lon))
	}

	if m.checkPoint(p3Lon, p3Lat) {
		*neighbours = append(*neighbours, m.GetNode(p3Lat, p3Lon))
	}

	if m.checkPoint(p4Lon, p4Lat) {
		*neighbours = append(*neighbours, m.GetNode(p4Lat, p4Lon))
	}

	if m.checkPoint(p5Lon, p5Lat) {
		*neighbours = append(*neighbours, m.GetNode(p5Lat, p5Lon))
	}

	if m.checkPoint(p6Lon, p6Lat) {
		*neighbours = append(*neighbours, m.GetNode(p6Lat, p6Lon))
	}

	if m.checkPoint(p7Lon, p7Lat) {
		*neighbours = append(*neighbours, m.GetNode(p7Lat, p7Lon))
	}

	if m.checkPoint(p8Lon, p8Lat) {
		*neighbours = append(*neighbours, m.GetNode(p8Lat, p8Lon))
	}

	for ix := -dl / 2; ix < dl/2; ix++ {
		for iy := -dl / 2; iy < dl/2; iy++ {
			m.GetNode(n.LatIdx-ix, n.LonIdx-iy).Visited = true
		}
	}
}
func (m *Map) GetNeighboursAdoptedV3(neighbours *[]*Node, n *Node) {
	dl := 1000

	p1Lon, p1Lat := n.LonIdx-dl, n.LatIdx
	p2Lon, p2Lat := n.LonIdx-dl, n.LatIdx+dl
	p3Lon, p3Lat := n.LonIdx, n.LatIdx+dl
	p4Lon, p4Lat := n.LonIdx+dl, n.LatIdx+dl
	p5Lon, p5Lat := n.LonIdx+dl, n.LatIdx
	p6Lon, p6Lat := n.LonIdx+dl, n.LatIdx-dl
	p7Lon, p7Lat := n.LonIdx, n.LatIdx-dl
	p8Lon, p8Lat := n.LonIdx-dl, n.LatIdx-dl

	if m.checkPoint(p1Lon, p1Lat) && m.checkLine(n.LonIdx, n.LatIdx, p1Lon, p1Lat) {
		*neighbours = append(*neighbours, m.GetNode(p1Lat, p1Lon))
	} else {
		m.addNonBlockedPoints(neighbours, n.LonIdx, n.LatIdx, p1Lon, p1Lat)
	}

	if m.checkPoint(p2Lon, p2Lat) {
		*neighbours = append(*neighbours, m.GetNode(p2Lat, p2Lon))
	} else {
		m.addNonBlockedPoints(neighbours, n.LonIdx, n.LatIdx, p2Lon, p2Lat)
	}

	if m.checkPoint(p3Lon, p3Lat) && m.checkLine(n.LonIdx, n.LatIdx, p3Lon, p3Lat) {
		*neighbours = append(*neighbours, m.GetNode(p3Lat, p3Lon))
	} else {
		m.addNonBlockedPoints(neighbours, n.LonIdx, n.LatIdx, p3Lon, p3Lat)
	}

	if m.checkPoint(p4Lon, p4Lat) {
		*neighbours = append(*neighbours, m.GetNode(p4Lat, p4Lon))
	} else {
		m.addNonBlockedPoints(neighbours, n.LonIdx, n.LatIdx, p4Lon, p4Lat)
	}

	if m.checkPoint(p5Lon, p5Lat) && m.checkLine(n.LonIdx, n.LatIdx, p5Lon, p5Lat) {
		*neighbours = append(*neighbours, m.GetNode(p5Lat, p5Lon))
	} else {
		m.addNonBlockedPoints(neighbours, n.LonIdx, n.LatIdx, p5Lon, p5Lat)
	}

	if m.checkPoint(p6Lon, p6Lat) {
		*neighbours = append(*neighbours, m.GetNode(p6Lat, p6Lon))
	} else {
		m.addNonBlockedPoints(neighbours, n.LonIdx, n.LatIdx, p6Lon, p6Lat)
	}

	if m.checkPoint(p7Lon, p7Lat) && m.checkLine(n.LonIdx, n.LatIdx, p7Lon, p7Lat) {
		*neighbours = append(*neighbours, m.GetNode(p7Lat, p7Lon))
	} else {
		m.addNonBlockedPoints(neighbours, n.LonIdx, n.LatIdx, p7Lon, p7Lat)
	}

	if m.checkPoint(p8Lon, p8Lat) {
		*neighbours = append(*neighbours, m.GetNode(p8Lat, p8Lon))
	} else {
		m.addNonBlockedPoints(neighbours, n.LonIdx, n.LatIdx, p8Lon, p8Lat)
	}

	for ix := -dl / 2; ix < dl/2; ix++ {
		for iy := -dl / 2; iy < dl/2; iy++ {
			m.GetNode(n.LatIdx-ix, n.LonIdx-iy).Visited = true
		}
	}
}

func (m *Map) addBlock(neighbours *[]*Node, lat, lon, size int) {
	if m.checkPoint(lon, lat) {
		n := m.GetNode(lat, lon)
		//n.Size = size
		*neighbours = append(*neighbours, n)
	} else {
		if size > 50 {
			dl := size / 8
			for i := -1; i <= 1; i++ {
				for j := -1; j <= 1; j++ {
					m.addBlock(neighbours, lat+dl*i, lon+dl*j, dl)
				}
			}
		}
	}
}

func (m *Map) GetNeighboursAdoptedV5(neighbours *[]*Node, n *Node) {
	dl := 1000

	p1Lon, p1Lat := n.LonIdx-dl, n.LatIdx
	p2Lon, p2Lat := n.LonIdx-dl, n.LatIdx+dl
	p3Lon, p3Lat := n.LonIdx, n.LatIdx+dl
	p4Lon, p4Lat := n.LonIdx+dl, n.LatIdx+dl
	p5Lon, p5Lat := n.LonIdx+dl, n.LatIdx
	p6Lon, p6Lat := n.LonIdx+dl, n.LatIdx-dl
	p7Lon, p7Lat := n.LonIdx, n.LatIdx-dl
	p8Lon, p8Lat := n.LonIdx-dl, n.LatIdx-dl

	m.addBlock(neighbours, p1Lat, p1Lon, dl)
	m.addBlock(neighbours, p2Lat, p2Lon, dl)
	m.addBlock(neighbours, p3Lat, p3Lon, dl)
	m.addBlock(neighbours, p4Lat, p4Lon, dl)
	m.addBlock(neighbours, p5Lat, p5Lon, dl)
	m.addBlock(neighbours, p6Lat, p6Lon, dl)
	m.addBlock(neighbours, p7Lat, p7Lon, dl)
	m.addBlock(neighbours, p8Lat, p8Lon, dl)

	for ix := -dl / 2; ix < dl/2; ix++ {
		for iy := -dl / 2; iy < dl/2; iy++ {
			if n.LatIdx-ix > 0 && n.LatIdx-ix < height && n.LonIdx-iy > 0 && n.LonIdx-iy < Width {
				m.GetNode(n.LatIdx-ix, n.LonIdx-iy).Visited = true
			}
		}
	}
}

func (m *Map) GetNeighboursAdoptedV6(neighbours *[]*Node, n *Node) {
	dl := 1000

	p1Lon, p1Lat := n.LonIdx-dl, n.LatIdx
	p2Lon, p2Lat := n.LonIdx-dl, n.LatIdx+dl
	p3Lon, p3Lat := n.LonIdx, n.LatIdx+dl
	p4Lon, p4Lat := n.LonIdx+dl, n.LatIdx+dl
	p5Lon, p5Lat := n.LonIdx+dl, n.LatIdx
	p6Lon, p6Lat := n.LonIdx+dl, n.LatIdx-dl
	p7Lon, p7Lat := n.LonIdx, n.LatIdx-dl
	p8Lon, p8Lat := n.LonIdx-dl, n.LatIdx-dl

	m.addBlock(neighbours, p1Lat, p1Lon, dl)
	m.addBlock(neighbours, p2Lat, p2Lon, dl)
	m.addBlock(neighbours, p3Lat, p3Lon, dl)
	m.addBlock(neighbours, p4Lat, p4Lon, dl)
	m.addBlock(neighbours, p5Lat, p5Lon, dl)
	m.addBlock(neighbours, p6Lat, p6Lon, dl)
	m.addBlock(neighbours, p7Lat, p7Lon, dl)
	m.addBlock(neighbours, p8Lat, p8Lon, dl)

	for ix := -dl / 2; ix < dl/2; ix++ {
		for iy := -dl / 2; iy < dl/2; iy++ {
			if n.LatIdx-ix > 0 && n.LatIdx-ix < height && n.LonIdx-iy > 0 && n.LonIdx-iy < Width {
				m.GetNode(n.LatIdx-ix, n.LonIdx-iy).Visited = true
			}
		}
	}
}
