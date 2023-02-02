package internal

import (
	"math"
)

const RoundingTolerance = 1e-9

type linePoint struct {
	lat    float64
	latIdx int
}

type genLonPoint struct {
	lon     float64
	cellLon float64
	lonIdx  int
}
type genLatPoint struct {
	lat    float64
	latIdx int
}

type LineOfSightChecking struct {
	grid *Grid
}

func NewLineOfSight(grid *Grid) *LineOfSightChecking {
	los := LineOfSightChecking{grid: grid}
	return &los
}

func (los *LineOfSightChecking) LineOfSight(startNode, endNode *Node, sourceIsCell, targetIsCell bool) bool {
	slat, slon := los.grid.LatLon(startNode)
	elat, elon := los.grid.LatLon(endNode)
	slatIdx, slonIdx := startNode.LatIdx, startNode.LonIdx
	elatIdx, elonIdx := endNode.LatIdx, endNode.LonIdx

	prevLat := slat
	prevLatIdx, prevLonIdx := slatIdx, slonIdx

	if math.Abs(slon-elon) < RoundingTolerance {
		cellLon := los.grid.Longitudes[prevLonIdx]
		for _, p := range genLatsMeridian(slat, elat, slatIdx, elatIdx, sourceIsCell, targetIsCell) {
			_ = Haversine(prevLat, slon, p.lat, slon)

			if !los.isTraversable(prevLatIdx, prevLonIdx, cellLon, prevLat, slon, p.lat, slon) {
				return false
			}

			prevLat = p.lat
			prevLatIdx = p.latIdx
		}
	} else {
		cellSlon := los.grid.Longitudes[slonIdx]
		cellElon := los.grid.Longitudes[elonIdx]
		nLon := len(los.grid.Longitudes)

		crossAntimeridian, lonDir, startLon, lastLon, n := rotateGc(slon, elon, slonIdx, elonIdx, cellSlon, cellElon, nLon)
		prevLon := startLon

		gc := NewGreatCircle(slat, startLon, elat, lastLon)
		c := 0

		for _, p := range genLons(prevLon, lastLon, slonIdx, elonIdx, sourceIsCell, targetIsCell, crossAntimeridian, lonDir, nLon) {
			ilat := gc.intersectionWithLongitude(p.lon)
			ilatIdx := BisecLat(los.grid.Latitudes, ilat)

			var lon_idx_m1 int

			if lonDir {
				lon_idx_m1 = normaliseLonIndex(p.lonIdx-1, nLon)
			} else {
				lon_idx_m1 = p.lonIdx
			}

			k := (p.lon - prevLon) / (ilat - prevLat + 1e-9)
			_lon := prevLon
			_lat := prevLat

			for _, pLat := range _gen_lats(prevLat, prevLatIdx, ilatIdx, sourceIsCell, c) {
				ilon := prevLon + k*(pLat.lat-prevLat)
				_ = Haversine(_lat, _lon, pLat.lat, ilon)

				if !los.isTraversable(pLat.latIdx, lon_idx_m1, p.cellLon, _lat, _lon, pLat.lat, ilon) {
					return false
				}

				_lat = pLat.lat
				_lon = ilon
			}

			var lat_dir bool

			if ilat > 0 {
				lat_dir = _lat < ilat
			} else {
				lat_dir = _lat > ilat
			}

			if targetIsCell && c == n && ((lat_dir && lonDir) || (!lat_dir && !lonDir)) {

			} else {
				_ = Haversine(_lat, _lon, ilat, p.lon)

				if !los.isTraversable(ilatIdx, lon_idx_m1, p.cellLon, _lat, _lon, ilat, p.lon) {
					return false
				}

			}

			c += 1

			prevLon = p.lon
			prevLat = ilat
			prevLatIdx = ilatIdx
		}
	}

	return true
}

func (los *LineOfSightChecking) isTraversable(lat_idx, lon_idx int, cell_lon, lat1, lon1, lat2, lon2 float64) bool {
	if los.grid.IsTraversable(lon_idx, lat_idx) {
		return true
	}

	return false
}

func genLons(slon, elon float64, slonIdx, elonIdx int, sourceIsCell, targetIsCell, crossAntimeridian, lonDir bool, nLon int) []genLonPoint {
	var rotate bool
	delta := 0.00872665
	lonPoints := []genLonPoint{}

	if crossAntimeridian {
		if lonDir {
			cellLon := -math.Pi
			lon := -math.Pi
			slonIdx += 1

			if slonIdx >= nLon {
				slonIdx -= nLon
				rotate = true
			} else {
				rotate = false
			}

			for (slonIdx < nLon && !rotate) || (slonIdx <= elonIdx && rotate) {
				lon += delta
				lonPoints = append(lonPoints, genLonPoint{
					lon:     lon,
					cellLon: cellLon,
					lonIdx:  slonIdx,
				})

				cellLon = lon
				slonIdx += 1

				if slonIdx >= nLon {
					slonIdx -= nLon
					rotate = true
				}
			}

		} else {
			rotate = false

			for (slonIdx >= 0 && !rotate) || (slonIdx > elonIdx && rotate) {
				lonPoints = append(lonPoints, genLonPoint{
					lon:     slon,
					cellLon: slon,
					lonIdx:  slonIdx,
				})

				slon -= delta
				slonIdx -= 1
				if slonIdx < 0 {
					slonIdx += nLon
					rotate = true
				}
			}

		}
	} else {
		if lonDir {
			cellLon := -math.Pi
			lon := -math.Pi
			slonIdx += 1

			for slonIdx <= elonIdx {
				lon += delta

				lonPoints = append(lonPoints, genLonPoint{
					lon:     lon,
					cellLon: cellLon,
					lonIdx:  slonIdx,
				})

				cellLon = lon
				slonIdx += 1
			}

		} else {
			lon := slon - math.Mod(slon, delta)

			if sourceIsCell {
				if slon-lon < 1e-15 {
					lon = math.Round(slon/delta)*delta - delta
				}
				slonIdx -= 1
			}

			for slonIdx > elonIdx {

				lonPoints = append(lonPoints, genLonPoint{
					lon:     lon,
					cellLon: lon,
					lonIdx:  slonIdx,
				})

				slonIdx -= 1
				lon -= delta
			}

			if targetIsCell {
				lonPoints = append(lonPoints, genLonPoint{
					lon:     lon,
					cellLon: lon,
					lonIdx:  slonIdx,
				})
			} else {

				lonPoints = append(lonPoints, genLonPoint{
					lon:     elon,
					cellLon: lon,
					lonIdx:  slonIdx,
				})
			}

		}
	}

	return lonPoints
}

func rotateGc(slon, elon float64, slonIdx, elonIdx int, cellSlon, cellElon float64, nLon int) (bool, bool, float64, float64, int) {
	var crossAntimeridian bool
	var lonDir bool
	var n int
	var endLon, startLon float64
	delta := 0.00872665

	if math.Abs(elon-slon) > math.Pi {
		crossAntimeridian = true

		if slon <= elon {
			lonDir = false
			n = nLon - elonIdx + slonIdx
			endLon = -math.Pi + (elon - cellElon)
			startLon = -math.Pi + float64(n)*delta + math.Abs(slon-cellSlon)
			n -= 1
		} else {
			lonDir = true
			n = nLon - slonIdx + elonIdx
			startLon = -math.Pi + (slon - cellSlon)
			endLon = -math.Pi + float64(n)*delta + math.Abs(elon-cellElon)
		}
	} else {
		crossAntimeridian = false
		if slon <= elon {
			lonDir = true
			n = elonIdx - slonIdx
			startLon = -math.Pi + (slon - cellSlon)
			endLon = -math.Pi + float64(n)*delta + math.Abs(elon-cellElon)
		} else {
			lonDir = false
			n = slonIdx - elonIdx
			startLon = -math.Pi + float64(n)*delta + (slon - cellSlon)
			endLon = -math.Pi + math.Abs(elon-cellElon)
			n -= 1
		}
	}

	return crossAntimeridian, lonDir, startLon, endLon, n
}

func genLatsMeridian(slat, elat float64, slatIdx, elatIdx int, sourceIsCell, targetIsCell bool) []linePoint {
	delta := 0.00872665
	slat -= math.Mod(slat, delta)
	var arr []linePoint

	if slat < elat {
		for slatIdx < elatIdx {
			slat += delta
			slatIdx += 1
			arr = append(arr, linePoint{
				lat:    slat,
				latIdx: slatIdx,
			})
		}
	} else {
		if sourceIsCell {
			slatIdx -= 1
		}
		for slatIdx > elatIdx {
			slatIdx -= 1
			arr = append(arr, linePoint{
				lat:    slat,
				latIdx: slatIdx,
			})
			slat -= delta
		}
		if targetIsCell {
			arr = append(arr, linePoint{
				lat:    elat,
				latIdx: elatIdx,
			})
		}
	}

	return arr
}

func BisecLat(a [height]float64, x float64) int {
	lo := 0
	hi := len(a)
	for lo < hi {
		mid := int((lo + hi) / 2)
		if x < a[mid] {
			hi = mid
		} else {
			lo = mid + 1
		}
	}

	if lo > 0 {
		return lo - 1
	}

	return 0
}

func _gen_lats(slat float64, slatIdx, elatIdx int, sourceIsCell bool, c int) []genLatPoint {
	delta := 0.00872665
	arr := []genLatPoint{}

	slat -= math.Mod(slat, delta)

	if slatIdx <= elatIdx {
		for slatIdx <= elatIdx {
			slat += delta
			arr = append(arr, genLatPoint{
				lat:    slat,
				latIdx: slatIdx,
			})
			slatIdx += 1
		}
	} else {
		if sourceIsCell && c == 0 {
			slatIdx -= 1
		}
		slat += delta
		for slatIdx > elatIdx {
			slat -= delta
			arr = append(arr, genLatPoint{
				lat:    slat,
				latIdx: slatIdx,
			})
			slatIdx -= 1
		}
	}

	return arr
}

func normaliseLonIndex(idx, nLon int) int {
	if idx < 0 {
		idx += nLon
	} else if idx >= nLon {
		idx -= nLon
	}

	return idx
}
