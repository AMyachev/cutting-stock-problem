package travelingSalesman

import "math"

// point format: x,y
func computeEuclideanDistance(firstPoint [2]float64, secondPoint [2]float64) float64 {
	xSqr := math.Pow(firstPoint[0]-secondPoint[0], 2)
	ySqr := math.Pow(firstPoint[1]-secondPoint[1], 2)
	return math.Sqrt(xSqr + ySqr)
}
