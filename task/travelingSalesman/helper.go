package travelingSalesman

import (
	"math"

	log "github.com/sirupsen/logrus"
)

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

func swap(objectsOrder []int, i int, j int) {
	temp := objectsOrder[i]
	objectsOrder[i] = objectsOrder[j]
	objectsOrder[j] = temp
}

func nextOrder(objectsOrder []int) bool {
	countObjects := len(objectsOrder)
	// choose penult object id
	objectIdx := countObjects - 2
	for objectIdx != -1 && objectsOrder[objectIdx] >= objectsOrder[objectIdx+1] {
		objectIdx--
	}
	// next order not found
	if objectIdx == -1 {
		return false
	}

	// choose last object id
	objectIdxLast := countObjects - 1
	for objectsOrder[objectIdx] >= objectsOrder[objectIdxLast] {
		objectIdxLast--
	}

	swap(objectsOrder, objectIdx, objectIdxLast)

	// sort objects after objectIdx
	left := objectIdx + 1
	right := countObjects - 1
	for left < right {
		swap(objectsOrder, left, right)
		left++
		right--
	}
	return true
}

func copySlice(source []int, destination []int) {
	if len(source) != len(destination) {
		log.Fatalf("length should be the same: [%d], [%d]", len(source), len(destination))
	}

	for i := 0; i < len(source); i++ {
		destination[i] = source[i]
	}
}
