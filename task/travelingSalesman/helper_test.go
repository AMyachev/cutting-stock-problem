package travelingSalesman

import (
	"fmt"
	"math"
	"testing"
)

func formatErrorReport(funcName string, outputValue interface{}, expectedOutputValue interface{}) string {
	return fmt.Sprintf("from func: %s expected %v but got %v", funcName, expectedOutputValue, outputValue)
}

func isEqualFloat64(first, second float64) bool {
	lambda := 0.000001
	return math.Abs(first-second) < lambda
}

func isEqualIntSlices(first, second []int) bool {
	if len(first) != len(second) {
		return false
	}

	for i := 0; i < len(first); i++ {
		if first[i] != second[i] {
			return false
		}
	}
	return true
}

func isEqualFloat64Matrixes(firstMatrix [][]float64, secondMatrix [][]float64) bool {
	if len(firstMatrix) != len(secondMatrix) {
		return false
	}

	for i := 0; i < len(firstMatrix); i++ {
		if len(firstMatrix[i]) != len(secondMatrix[i]) {
			return false
		}
		for j := 0; j < len(firstMatrix[i]); j++ {
			if !isEqualFloat64(firstMatrix[i][j], secondMatrix[i][j]) {
				return false
			}
		}
	}
	return true
}

func TestComputeEuclideanDistance(t *testing.T) {
	firstPoint := [2]float64{1.0, 1.0}
	secondPoint := [2]float64{2.0, 2.0}

	expectedDistance := 0.
	if distance := computeEuclideanDistance(firstPoint, firstPoint); !isEqualFloat64(distance, expectedDistance) {
		t.Error(formatErrorReport("computeEuclideanDistance", distance, expectedDistance))
	}

	expectedDistance = 1.4142135623730951
	if distance := computeEuclideanDistance(firstPoint, secondPoint); !isEqualFloat64(distance, expectedDistance) {
		t.Error(formatErrorReport("computeEuclideanDistance", distance, expectedDistance))
	}
}

func TestComputeWeightCenter(t *testing.T) {
	firstPoint := [2]float64{1.0, 1.0}
	secondPoint := [2]float64{2.0, 2.0}

	expectedWeightCenter := [2]float64{1.5, 1.5}
	if weightCenter := computeWeightCenter(firstPoint, 1, secondPoint); weightCenter != expectedWeightCenter {
		t.Error(formatErrorReport("computeWeightCenter", weightCenter, expectedWeightCenter))
	}
}

func TestComputeAllEuclideanDistance(t *testing.T) {
	firstPoint := [2]float64{1.0, 1.0}
	secondPoint := [2]float64{2.0, 2.0}
	thirdPoint := [2]float64{2.0, 1.0}

	points := [][2]float64{firstPoint, secondPoint, thirdPoint}

	expectedDistances := [][]float64{
		[]float64{0., 1.4142135623730951, 1.0},
		[]float64{1.4142135623730951, 0., 1.0},
		[]float64{1.0, 1.0, 0.},
	}
	if distances := computeAllEuclideanDistance(points); !isEqualFloat64Matrixes(distances, expectedDistances) {
		t.Error(formatErrorReport("computeAllEuclideanDistance", distances, expectedDistances))
	}
}

func TestNextOrder(t *testing.T) {
	currentOrder := []int{1, 2, 3}

	expectedOrders := [][]int{
		[]int{1, 3, 2},
		[]int{2, 1, 3},
		[]int{2, 3, 1},
		[]int{3, 1, 2},
		[]int{3, 2, 1},
	}

	for i := 0; i < len(expectedOrders); i++ {
		nextOrder(currentOrder)
		if !isEqualIntSlices(currentOrder, expectedOrders[i]) {
			t.Error(formatErrorReport("nextOder", currentOrder, expectedOrders[i]))
		}
	}
}
