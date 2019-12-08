package travelingSalesman

import (
	"math"

	log "github.com/sirupsen/logrus"
)

// input order should be sorted by increasing
// note: objectsOrder will be modificated
func bruteForce(objectsOrder []int, betweenObjectsLength [][]float64) []int {
	bestObjectsOrder := make([]int, len(objectsOrder))
	minCritValue := math.Inf(1)

	for nextOrder(objectsOrder) {
		if critValue := criterion(objectsOrder, betweenObjectsLength); critValue < minCritValue {
			minCritValue = critValue
			copySlice(objectsOrder, bestObjectsOrder)
		}
	}
	return bestObjectsOrder
}

// compute cycle length
func criterion(objectOrder []int, betweenObjectsLength [][]float64) float64 {
	critValue := 0.
	for i := 0; i < len(objectOrder)-1; i++ {
		critValue += betweenObjectsLength[objectOrder[i]][objectOrder[i+1]]
	}
	critValue += betweenObjectsLength[objectOrder[len(objectOrder)-1]][objectOrder[0]]
	return critValue
}

func standardReducto(task *travelingSalesmanSubTask, alpha int) []*travelingSalesmanSubTask {
	if alpha <= 0 {
		log.WithField("alpha", alpha).Fatal("alpha <= 0")
	}

	if alpha == 1 || task.CountTown() <= alphaMax {
		return []*travelingSalesmanSubTask{task}
	}

	clustersCenters := []int{-1, -1}

	// find first 2 centers
	maxLength := 0.
	for i := 0; i < task.countTown; i++ {
		for j := i + 1; j < task.countTown; j++ {
			if length := task.length(i, j); length > maxLength {
				maxLength = length
				clustersCenters[0] = i
				clustersCenters[1] = j
			}
		}
	}

	// helper function
	find := func(clustersCenters []int, town int) bool {
		for _, centerTown := range clustersCenters {
			if town == centerTown {
				return true
			}
		}
		return false
	}

	// find other centers
	for clusterNumber := 2; clusterNumber < alpha; clusterNumber++ {
		maxSumLength := 0.
		nearestTown := 0
		for town := 0; town < task.countTown; town++ {
			// center?
			if find(clustersCenters, town) {
				continue
			}

			sumLength := 0.
			for _, centerTown := range clustersCenters {
				sumLength += task.length(town, centerTown)
			}

			if sumLength > maxSumLength {
				maxSumLength = sumLength
				nearestTown = town
			}
		}
		clustersCenters = append(clustersCenters, nearestTown)
	}

	// fill clusters with other towns

	// fill with cluster's centers
	towns := make([][]int, alpha)
	for i, clusterCenter := range clustersCenters {
		towns[i] = []int{clusterCenter}
	}

	clustersCenterWeight := make([][2]float64, alpha)
	for i := 0; i < alpha; i++ {
		clustersCenterWeight[i] = task.coords(clustersCenters[i])
	}

	// fill with others
	for town := 0; town < task.countTown; town++ {
		// center?
		if find(clustersCenters, town) {
			continue
		}

		minLength := math.Inf(1)
		clusterCenterIdx := 0
		for i, clusterCenterWeight := range clustersCenterWeight {
			if length := computeEuclideanDistance(task.coords(town), clusterCenterWeight); length < minLength {
				minLength = length
				clusterCenterIdx = i
			}
		}
		towns[clusterCenterIdx] = append(towns[clusterCenterIdx], town)
		// update weight of cluster center
		clustersCenterWeight[clusterCenterIdx] = computeWeightCenter(clustersCenterWeight[clusterCenterIdx], task.coords(town))
	}

	// fast fix

	result := make([]*travelingSalesmanSubTask, alpha)
	for i := 0; i < alpha; i++ {
		townsNumbers := make([]int, len(towns[i]))
		for j := 0; j < len(townsNumbers); j++ {
			townsNumbers[j] = task.towns[towns[i][j]]
		}
		result[i] = MakeTravelingSalesmanSubTask(townsNumbers, task.townsCoord, task.betweenTownsLength)
	}

	return result
}

func standardGreedy(task *travelingSalesmanSubTask) *travelingSalesmanSolution {
	resultOrder := []int{task.towns[0]}
	remainingTowns := make([]int, len(task.towns)-1)
	copySlice(task.towns[1:], remainingTowns)

	for i := 0; i < task.countTown-1; i++ {
		lastTown := resultOrder[len(resultOrder)-1]
		minLength := math.Inf(1)
		townPos := 0
		for pos, remainingTown := range remainingTowns {
			if length := task.betweenTownsLength[lastTown][remainingTown]; length < minLength {
				minLength = length
				townPos = pos
			}
		}
		resultOrder = append(resultOrder, remainingTowns[townPos])
		remainingTowns = append(remainingTowns[:townPos], remainingTowns[townPos+1:]...)
	}

	return &travelingSalesmanSolution{
		towns: resultOrder,
	}
}
