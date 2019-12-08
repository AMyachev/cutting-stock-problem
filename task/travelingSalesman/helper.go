package travelingSalesman

import "math"

// point format: x,y
func computeEuclideanDistance(firstPoint [2]float64, secondPoint [2]float64) float64 {
	xSqr := math.Pow(firstPoint[0]-secondPoint[0], 2)
	ySqr := math.Pow(firstPoint[1]-secondPoint[1], 2)
	return math.Sqrt(xSqr + ySqr)
}

func computeWeightCenter(firstPoint [2]float64, secondPoint [2]float64) [2]float64 {
	return [2]float64{(firstPoint[0] + secondPoint[0]) / 2, (firstPoint[1] + secondPoint[1]) / 2}
}

func computeAllEuclideanDistance(clusterPoints [][2]float64) [][]float64 {
	countPoints := len(clusterPoints)

	// allocation
	betweenPointsLength := make([][]float64, countPoints)
	for i := 0; i < countPoints; i++ {
		betweenPointsLength[i] = make([]float64, countPoints)
	}

	// initialization
	for i := 0; i < countPoints-1; i++ {
		for j := i + 1; j < countPoints; j++ {
			betweenPointsLength[i][j] = computeEuclideanDistance(clusterPoints[i], clusterPoints[j])
			// Euclidean distance is symmetric, so ...
			betweenPointsLength[j][i] = betweenPointsLength[i][j]
		}
	}

	return betweenPointsLength
}
