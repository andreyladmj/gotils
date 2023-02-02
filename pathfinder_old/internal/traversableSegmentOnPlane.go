package internal

import (
	"math"
)

func TraversableSegmentOnPlane(sx, sy, ex, ey int, grid *Grid) bool {
	var dx, dy float64
	dx = float64(ex - sx)
	dy = float64(ey - sy)

	dist := math.Abs(dx) + math.Abs(dy) // Manhattan distance
	// Note: this function will fail if (sx, sy) == (ex, ey).
	// Division by zero must be verified / handled outside this function

	if dist == 0 {
		return true
	}

	dx /= dist
	dy /= dist

	var startIdx, endIdx, xSign, ySign, step int

	if dy == 0 { // Movement along horizontal axis
		if dx > 0 {
			startIdx = 0
			endIdx = int(dist)
			xSign = 1
		} else {
			startIdx = 1
			endIdx = int(dist + 1)
			xSign = -1
		}

		sym1 := sy - 1

		for i := startIdx; i < endIdx; i++ {
			if !grid.IsTraversable(sx+xSign*i, sy) || !grid.IsTraversable(sx+xSign*i, sym1) {
				return false
			}
		}

	} else if dx == 0 { // Movement along vertical axis
		if dy > 0 {
			startIdx = 0
			endIdx = int(dist)
			ySign = 1
		} else {
			startIdx = 1
			endIdx = int(dist + 1)
			ySign = -1
		}

		sxm1 := sx - 1

		for i := startIdx; i < endIdx; i++ {
			if !grid.IsTraversable(sx, sy+ySign*i) || !grid.IsTraversable(sxm1, sy+ySign*i) {
				return false
			}
		}
	} else { // Any angle movement
		if math.Abs(dx) == 0.5 && math.Abs(dy) == 0.5 {
			if dx > 0 && dy > 0 {
				startIdx = 1
				endIdx = int(dist)
				step = 2
			} else {
				startIdx = 1
				endIdx = int(dist + 1)
				step = 2
			}
		} else {
			if dx > 0 && dy > 0 {
				startIdx = 0
				endIdx = int(dist)
				step = 1
			} else {
				startIdx = 1
				endIdx = int(dist)
				step = 1
			}
		}

		for i := startIdx; i < endIdx; i += step {
			x := int(math.Floor(float64(sx) + dx*float64(i)))
			y := int(math.Floor(float64(sy) + dy*float64(i)))
			if !grid.IsTraversable(x, y) {
				return false
			}
		}
	}

	return true
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
func Min(x, y int) int {
	if x < y {
		return x
	}

	return y
}
