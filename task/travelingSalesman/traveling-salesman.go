package travelingSalesman

import (
	"fmt"
	"io/ioutil"
	"math"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

const (
	alphaMax = 10
)

type travelingSalesmanTask struct {
	countTown          int
	townsCoord         [][2]float64
	betweenTownsLength [][]float64
}

type travelingSalesmanSolution struct {
	towns []int
}

func MakeTravelingSalesmanTask(taskFile string) *travelingSalesmanTask {
	content, err := ioutil.ReadFile(taskFile)
	if err != nil {
		log.Fatal(err)
	}

	townsCoord := make([][2]float64, 0)
	dimension := 0

	contentLines := strings.Split(string(content), "\r\n")
	for i, contentLine := range contentLines {

		if strings.HasPrefix(contentLine, "DIMENSION") {
			// line with dimension was found
			// line format: "DIMENSION: [NUMBER]"
			splitted := strings.Split(contentLine, ":")
			dimension, err = strconv.Atoi(strings.Trim(splitted[1], " "))
			if err != nil {
				log.Fatal(err)
			}
			townsCoord = make([][2]float64, dimension)
		}

		if strings.HasPrefix(contentLine, "NODE_COORD_SECTION") {
			// init line with town's coords was found
			// lines format: "[town number; start with 1] [town's X coord] [town's Y coord]"
			for j, coordsLine := range contentLines[i+1:] {
				if strings.HasPrefix(coordsLine, "END") {
					break
				}

				splitted := strings.Split(coordsLine, " ")
				x, parseXErr := strconv.ParseFloat(splitted[1], 64)
				y, parseYErr := strconv.ParseFloat(splitted[2], 64)
				if parseXErr != nil || parseYErr != nil {
					log.Fatal(err)
				}
				townsCoord[j] = [2]float64{x, y}
			}
			break
		}
	}

	betweenTownsLength := computeAllEuclideanDistance(townsCoord)

	return &travelingSalesmanTask{
		countTown:          dimension,
		townsCoord:         townsCoord,
		betweenTownsLength: betweenTownsLength,
	}
}

// point format: x,y
func computeEuclideanDistance(firstPoint [2]float64, secondPoint [2]float64) float64 {
	xSqr := math.Pow(firstPoint[0]-secondPoint[0], 2)
	ySqr := math.Pow(firstPoint[1]-secondPoint[1], 2)
	return math.Sqrt(xSqr + ySqr)
}

func computeWeightCenter(firstPoint [2]float64, secondPoint [2]float64) [2]float64 {
	return [2]float64{math.Abs(firstPoint[0]-secondPoint[0]) / 2, math.Abs(firstPoint[1]-secondPoint[1]) / 2}
}

func computeAllEuclideanDistance(townsCoord [][2]float64) [][]float64 {
	countTowns := len(townsCoord)

	// allocation
	betweenTownsLength := make([][]float64, countTowns)
	for i := 0; i < countTowns; i++ {
		betweenTownsLength[i] = make([]float64, countTowns)
	}

	// initialization
	for i := 0; i < countTowns-1; i++ {
		for j := i + 1; j < countTowns; j++ {
			betweenTownsLength[i][j] = computeEuclideanDistance(townsCoord[i], townsCoord[j])
			// Euclidean distance is symmetric, so ...
			betweenTownsLength[j][i] = betweenTownsLength[i][j]
		}
	}

	return betweenTownsLength
}

func (task *travelingSalesmanTask) Compute(reductoAlgoName string, alpha int, betta int) (solution *travelingSalesmanSolution) {

	// towns can be a subset of {1,2, ..., n}
	subTaskCountTown := task.countTown
	towns := make([]int, subTaskCountTown)
	for i := 0; i < subTaskCountTown; i++ {
		towns[i] = i
	}

	subTask := MakeTravelingSalesmanSubTask(task, towns)

	var reductoAlgo func(*travelingSalesmanSubTask, int) []*travelingSalesmanSubTask

	switch reductoAlgoName {
	case "standard":
		reductoAlgo = reducto
	default:
		panic("not implemeted")
	}
	return subTask.Compute(reductoAlgo, alpha, betta)
}

type travelingSalesmanSubTask struct {
	countTown          int
	towns              []int
	townsCoord         [][2]float64
	betweenTownsLength [][]float64
}

func MakeTravelingSalesmanSubTask(task *travelingSalesmanTask, towns []int) *travelingSalesmanSubTask {
	return &travelingSalesmanSubTask{
		countTown:          task.countTown,
		towns:              towns,
		townsCoord:         task.townsCoord,
		betweenTownsLength: task.betweenTownsLength,
	}
}

func MakeTravelingSalesmanSubTaskFromSubTask(task *travelingSalesmanSubTask, towns []int) *travelingSalesmanSubTask {
	return &travelingSalesmanSubTask{
		countTown:          len(towns),
		towns:              towns,
		townsCoord:         task.townsCoord,
		betweenTownsLength: task.betweenTownsLength,
	}
}

func (task *travelingSalesmanSubTask) length(firstTown int, secondTown int) float64 {
	return task.betweenTownsLength[task.towns[firstTown]][task.towns[secondTown]]
}

func (task *travelingSalesmanSubTask) coords(town int) [2]float64 {
	return task.townsCoord[task.towns[town]]
}

func reducto(task *travelingSalesmanSubTask, alpha int) []*travelingSalesmanSubTask {
	if alpha <= 0 {
		log.WithField("alpha", alpha).Fatal("alpha <= 0")
	}

	if alpha == 1 {
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

	result := make([]*travelingSalesmanSubTask, alpha)
	for i := 0; i < alpha; i++ {
		result[i] = MakeTravelingSalesmanSubTaskFromSubTask(task, towns[i])
	}

	return result
}

func (task *travelingSalesmanSubTask) Compute(reducto func(*travelingSalesmanSubTask, int) []*travelingSalesmanSubTask,
	alpha int, betta int) (solution *travelingSalesmanSolution) {

	solutions := []*travelingSalesmanSolution{}

	if alpha > alphaMax {
		panic(fmt.Sprintf("compute: alpha [%d] > max alpha [%d]", alpha, alphaMax))
	}

	subTasks := reducto(task, alpha)
	betta--

	orderedSubTasks := task.computeExternalTask(subTasks)

	for _, subTask := range orderedSubTasks {
		if betta != 0 {
			solutions = append(solutions, subTask.Compute(reducto, alpha, betta))
		} else {
			if subTask.CountTown() > alphaMax {
				solutions = append(solutions, subTask.Greedy())
			} else {
				solutions = append(solutions, subTask.ExhaustiveSearch())
			}
		}
	}

	solution = task.CombineSolutions(solutions)

	return solution
}

func (task *travelingSalesmanSubTask) computeExternalTask(subTasks []*travelingSalesmanSubTask) []*travelingSalesmanSubTask {
	panic("not implemented")
}

func (task *travelingSalesmanSubTask) CountTown() int {
	return task.countTown
}

func (task *travelingSalesmanSubTask) Greedy() *travelingSalesmanSolution {
	panic("not implemented")
}

func (task *travelingSalesmanSubTask) ExhaustiveSearch() *travelingSalesmanSolution {
	panic("not implemented")
}

func (task *travelingSalesmanSubTask) CombineSolutions(solutions []*travelingSalesmanSolution) *travelingSalesmanSolution {
	panic("not implemented")
}
