package travelingSalesman

import (
	"fmt"
	"testing"
)

func formatErrorReport(funcName string, outputValue interface{}, expectedOutputValue interface{}) string {
	return fmt.Sprintf("from func: %s expected %v but got %v", funcName, expectedOutputValue, outputValue)
}

func TestComputeEuclideanDistance(t *testing.T) {
	firstPoint := [2]float64{1.0, 1.0}
	secondPoint := [2]float64{2.0, 2.0}

	expectedDistance := 0.
	if distance := computeEuclideanDistance(firstPoint, firstPoint); distance != expectedDistance {
		t.Error(formatErrorReport("computeEuclideanDistance", distance, expectedDistance))
	}

	expectedDistance = 1.4142135623730951
	if distance := computeEuclideanDistance(firstPoint, secondPoint); distance != expectedDistance {
		t.Error(formatErrorReport("computeEuclideanDistance", distance, expectedDistance))
	}
}

func TestComputeWeightCenter(t *testing.T) {
	firstPoint := [2]float64{1.0, 1.0}
	secondPoint := [2]float64{2.0, 2.0}

	expectedWeightCenter := [2]float64{1.5, 1.5}
	if weightCenter := computeWeightCenter(firstPoint, secondPoint); weightCenter != expectedWeightCenter {
		t.Error(formatErrorReport("computeWeightCenter", weightCenter, expectedWeightCenter))
	}
}
