package pathfinder

import (
	"image"
	"image/color"
	"math"
	"pathfinder/internal"
	"sync"
)

func fillGreatCircle(lat1, lon1, lat2, lon2 float64, grid *internal.Grid) bool {
	// AVG_EARTH_RADIUS_NM = 3440.0707986216
	step := grid.Step

	cosLat1 := math.Cos(lat1)
	sinLat1 := math.Sin(lat1)
	cosLon1 := math.Cos(lon1)
	sinLon1 := math.Sin(lon1)

	cosLat2 := math.Cos(lat2)
	sinLat2 := math.Sin(lat2)
	cosLon2 := math.Cos(lon2)
	sinLon2 := math.Sin(lon2)

	d := math.Acos(sinLat1*sinLat2 + cosLat1*cosLat2*math.Cos(lon1-lon2))
	sind := math.Sin(d)
	if sind == 0 {
		d += 1e-9
	}

	cosLat1CosLon1 := cosLat1 * cosLon1
	cosLat1SinLon1 := cosLat1 * sinLon1
	cosLat2CosLon2 := cosLat2 * cosLon2
	cosLat2SinLon2 := cosLat2 * sinLon2

	step /= d

	getGCLatLon := func(f float64) (float64, float64) {
		a := math.Sin((1 - f) * d)
		b := math.Sin(f * d)
		x := a*cosLat1CosLon1 + b*cosLat2CosLon2
		y := a*cosLat1SinLon1 + b*cosLat2SinLon2
		z := a*sinLat1 + b*sinLat2
		lat := math.Atan2(z, math.Pow(math.Pow(x, 2)+math.Pow(y, 2), 0.5))
		lon := math.Atan2(y, x)
		return lat, lon
	}

	//N := int(1 / step)

	//if N < 50 {
	for f := step; f <= 1+1e-9; f += step {
		lat, lon := getGCLatLon(f)
		latIdx := grid.BisecLat(lat)
		lonIdx := grid.BisecLon(lon)

		if !grid.IsTraversable(lonIdx, latIdx) {
			return false
		}

		if f > step {
			prevLat, prevLon := getGCLatLon(f - step)
			prevLatIdx := grid.BisecLat(prevLat)
			prevLonIdx := grid.BisecLon(prevLon)

			if !internal.TraversableSegmentOnPlane(lonIdx, latIdx, prevLonIdx, prevLatIdx, grid) {
				return false
			}
		}
	}

	return true
}

func fillGreatCircleGoroutines(lat1, lon1, lat2, lon2 float64, grid *internal.Grid) bool {
	// AVG_EARTH_RADIUS_NM = 3440.0707986216
	step := grid.Step

	cosLat1 := math.Cos(lat1)
	sinLat1 := math.Sin(lat1)
	cosLon1 := math.Cos(lon1)
	sinLon1 := math.Sin(lon1)

	cosLat2 := math.Cos(lat2)
	sinLat2 := math.Sin(lat2)
	cosLon2 := math.Cos(lon2)
	sinLon2 := math.Sin(lon2)

	d := math.Acos(sinLat1*sinLat2 + cosLat1*cosLat2*math.Cos(lon1-lon2))
	sind := math.Sin(d)
	if sind == 0 {
		d += 1e-9
	}

	cosLat1CosLon1 := cosLat1 * cosLon1
	cosLat1SinLon1 := cosLat1 * sinLon1
	cosLat2CosLon2 := cosLat2 * cosLon2
	cosLat2SinLon2 := cosLat2 * sinLon2

	step /= d

	getGCLatLon := func(f float64) (float64, float64) {
		a := math.Sin((1 - f) * d)
		b := math.Sin(f * d)
		x := a*cosLat1CosLon1 + b*cosLat2CosLon2
		y := a*cosLat1SinLon1 + b*cosLat2SinLon2
		z := a*sinLat1 + b*sinLat2
		lat := math.Atan2(z, math.Pow(math.Pow(x, 2)+math.Pow(y, 2), 0.5))
		lon := math.Atan2(y, x)
		return lat, lon
	}

	//N := int(1 / step)

	var wg sync.WaitGroup
	nonTraversablePoints := 0

	for f := step; f <= 1+1e-9; f += step {
		wg.Add(1)
		go func(wg *sync.WaitGroup, f float64) {
			defer wg.Done()

			lat, lon := getGCLatLon(f)
			latIdx := grid.BisecLat(lat)
			lonIdx := grid.BisecLon(lon)

			if !grid.IsTraversable(lonIdx, latIdx) {
				nonTraversablePoints++
			}

			if nonTraversablePoints > 0 {
				return
			}

			if f > step {
				prevLat, prevLon := getGCLatLon(f - step)
				prevLatIdx := grid.BisecLat(prevLat)
				prevLonIdx := grid.BisecLon(prevLon)

				if !internal.TraversableSegmentOnPlane(lonIdx, latIdx, prevLonIdx, prevLatIdx, grid) {
					nonTraversablePoints++
				}
			}
		}(&wg, f)
	}

	wg.Wait()

	return nonTraversablePoints == 0
}

func DrawGreatCircle(node1, node2 *internal.Node, grid *internal.Grid, img *image.RGBA) {
	// d between lats 0.0029102294150895602
	// gc step 0.07267292292361319
	// AVG_EARTH_RADIUS_NM = 3440.0707986216
	step := grid.Step

	lat1 := grid.Latitudes[node1.LatIdx]
	lon1 := grid.Longitudes[node1.LonIdx]
	lat2 := grid.Latitudes[node2.LatIdx]
	lon2 := grid.Longitudes[node2.LonIdx]

	cosLat1 := math.Cos(lat1)
	sinLat1 := math.Sin(lat1)
	cosLon1 := math.Cos(lon1)
	sinLon1 := math.Sin(lon1)

	cosLat2 := math.Cos(lat2)
	sinLat2 := math.Sin(lat2)
	cosLon2 := math.Cos(lon2)
	sinLon2 := math.Sin(lon2)

	d := math.Acos(sinLat1*sinLat2 + cosLat1*cosLat2*math.Cos(lon1-lon2))
	sind := math.Sin(d)
	if sind == 0 {
		d += 1e-9
	}

	cosLat1CosLon1 := cosLat1 * cosLon1
	cosLat1SinLon1 := cosLat1 * sinLon1
	cosLat2CosLon2 := cosLat2 * cosLon2
	cosLat2SinLon2 := cosLat2 * sinLon2

	step /= d

	c := color.RGBA{255, 255, 0, 255}
	dx := 2

	if internal.Width == 2160 {
		dx = 1
	}

	for f := step; f <= 1+1e-9; f += step {
		a := math.Sin((1 - f) * d)
		b := math.Sin(f * d)
		x := a*cosLat1CosLon1 + b*cosLat2CosLon2
		y := a*cosLat1SinLon1 + b*cosLat2SinLon2
		z := a*sinLat1 + b*sinLat2
		lat := math.Atan2(z, math.Pow(math.Pow(x, 2)+math.Pow(y, 2), 0.5))
		lon := math.Atan2(y, x)

		latIdx := grid.BisecLat(lat)
		lonIdx := grid.BisecLon(lon)

		for i := -dx; i < dx; i++ {
			for j := -dx; j < dx; j++ {
				img.Set(lonIdx+i, latIdx+j, c)
			}
		}
	}
}
