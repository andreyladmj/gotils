package internal

import "math"

type GreatCircle struct {
	slon         float64
	elon         float64
	sin1Cos2DivC float64
	sin2Cos1DivC float64
}

func NewGreatCircle(slat, slon, elat, elon float64) *GreatCircle {
	sinLon12 := math.Sin(slon - elon)

	gc := GreatCircle{
		slon:         slat,
		elon:         elon,
		sin1Cos2DivC: math.Tan(slat) / (sinLon12 + 1e-9),
		sin2Cos1DivC: math.Tan(elat) / (sinLon12 + 1e-9),
	}

	return &gc
}

func (gc *GreatCircle) intersectionWithLongitude(lon float64) float64 {
	if gc.slon == gc.elon {
		if lon == gc.elon {
			return 1.0
		}

		return 0.0
	}

	return math.Atan(gc.sin1Cos2DivC*math.Sin(lon-gc.elon)) - gc.sin2Cos1DivC*math.Sin(lon-gc.slon)
}
