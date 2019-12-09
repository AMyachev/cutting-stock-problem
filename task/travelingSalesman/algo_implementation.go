package travelingSalesman

import (
	"math"
	"sort"

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

	if alpha == 1 || task.countTown <= alphaMax {
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
	find := func(clustersCenters []int, townIdx int) bool {
		for _, centerTownIdx := range clustersCenters {
			if townIdx == centerTownIdx {
				return true
			}
		}
		return false
	}

	// find other centers
	for clusterNumber := 2; clusterNumber < alpha; clusterNumber++ {
		maxSumLength := 0.
		nearestTownIdx := 0
		for townIdx := 0; townIdx < task.countTown; townIdx++ {
			// center?
			if find(clustersCenters, townIdx) {
				continue
			}

			sumLength := 0.
			for _, centerTownIdx := range clustersCenters {
				sumLength += task.length(townIdx, centerTownIdx)
			}

			if sumLength > maxSumLength {
				maxSumLength = sumLength
				nearestTownIdx = townIdx
			}
		}
		clustersCenters = append(clustersCenters, nearestTownIdx)
	}

	// fill clusters with other towns

	// fill with cluster's centers
	townsIdx := make([][]int, alpha)
	for i, clusterCenter := range clustersCenters {
		townsIdx[i] = []int{clusterCenter}
	}

	clustersCenterWeight := make([][2]float64, alpha)
	for i := 0; i < alpha; i++ {
		clustersCenterWeight[i] = task.coords(clustersCenters[i])
	}

	// fill with others
	for townIdx := 0; townIdx < task.countTown; townIdx++ {
		// center?
		if find(clustersCenters, townIdx) {
			continue
		}

		minLength := math.Inf(1)
		clusterCenterIdx := 0
		for i, clusterCenterWeight := range clustersCenterWeight {
			if length := computeEuclideanDistance(task.coords(townIdx), clusterCenterWeight); length < minLength {
				minLength = length
				clusterCenterIdx = i
			}
		}

		// update weight of cluster center
		clustersCenterWeight[clusterCenterIdx] = computeWeightCenter(clustersCenterWeight[clusterCenterIdx], len(townsIdx[clusterCenterIdx]), task.coords(townIdx))

		townsIdx[clusterCenterIdx] = append(townsIdx[clusterCenterIdx], townIdx)
	}

	// conversion of the city index to its number
	result := make([]*travelingSalesmanSubTask, alpha)
	for i := 0; i < alpha; i++ {
		towns := make([]int, len(townsIdx[i]))
		for j := 0; j < len(towns); j++ {
			towns[j] = task.towns[townsIdx[i][j]]
		}
		result[i] = MakeTravelingSalesmanSubTask(towns, task.townsCoord, task.betweenTownsLength)
	}

	return result
}

func modifReducto(task *travelingSalesmanSubTask, alpha int) []*travelingSalesmanSubTask {
	subTasks := standardReducto(task, alpha)

	sort.Slice(subTasks, func(i, j int) bool { return subTasks[i].countTown < subTasks[j].countTown })

	subTasksWeightCenters := make([][2]float64, len(subTasks))
	for i := 0; i < len(subTasks); i++ {
		subTasksWeightCenters[i] = subTasks[i].computeClusterWeightCenter()
	}

	// filter
	remove := make([]bool, len(subTasks))

	for i := 0; i < len(subTasks)-1; i++ {
		if subTasks[i].countTown <= alphaMax {
			// found nearest cluster
			minLenght := math.Inf(1)
			nearestCluster := 0
			for j := i + 1; j < len(subTasks); j++ {
				if length := computeEuclideanDistance(subTasksWeightCenters[i], subTasksWeightCenters[j]); length < minLenght {
					minLenght = length
					nearestCluster = j
				}
			}

			subTasks[nearestCluster].towns = append(subTasks[nearestCluster].towns, subTasks[i].towns...)
			subTasks[nearestCluster].countTown += len(subTasks[i].towns)
			// modificate weight center
			subTasksWeightCenters[nearestCluster] = subTasks[nearestCluster].computeClusterWeightCenter()
			remove[i] = true
		}
	}

	result := []*travelingSalesmanSubTask{}
	for i := 0; i < len(subTasks); i++ {
		if !remove[i] {
			result = append(result, subTasks[i])
		}
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
