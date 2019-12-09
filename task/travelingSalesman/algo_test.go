package travelingSalesman

import "testing"

func TestCriterion(t *testing.T) {
	firstPoint := [2]float64{1.0, 1.0}
	secondPoint := [2]float64{2.0, 2.0}
	thirdPoint := [2]float64{2.0, 1.0}

	points := [][2]float64{firstPoint, secondPoint, thirdPoint}

	distances := computeAllEuclideanDistance(points)

	expectedCritValue := 3.414213562373095
	if critValue := criterion([]int{0, 1, 2}, distances); !isEqualFloat64(critValue, expectedCritValue) {
		t.Error(formatErrorReport("criterion", critValue, expectedCritValue))
	}
}
