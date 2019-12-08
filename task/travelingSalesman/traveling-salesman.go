package travelingSalesman

import (
	"fmt"
	"io/ioutil"
	"math"
	"sort"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

const (
	alphaMax = 10
)

var BestMinLength map[string]float64 = map[string]float64{
	"DJ38":   6656.,
	"EI8246": 206171,
	"UY734":  79114.,
	"XQF131": 564.,
	"XQG237": 1019.,
	"XQL662": 2513.,
}

type travelingSalesmanSolution struct {
	towns []int
}

type travelingSalesmanTask struct {
	countTown          int
	townsCoord         [][2]float64
	betweenTownsLength [][]float64
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
				if strings.HasPrefix(coordsLine, "EOF") {
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

func (task *travelingSalesmanTask) Criterion(solution *travelingSalesmanSolution) float64 {
	return criterion(solution.towns, task.betweenTownsLength)
}

func (task *travelingSalesmanTask) Compute(reductoAlgoName, greedyAlgoName string, alpha, betta int) (solution *travelingSalesmanSolution) {

	// towns can be a subset of {0,1, ..., n-1}
	subTaskCountTown := task.countTown
	towns := make([]int, subTaskCountTown)
	for i := 0; i < subTaskCountTown; i++ {
		towns[i] = i
	}

	subTask := MakeTravelingSalesmanSubTask(towns, task.townsCoord, task.betweenTownsLength)

	var reductoAlgo func(*travelingSalesmanSubTask, int) []*travelingSalesmanSubTask
	var greedyAlgo func(*travelingSalesmanSubTask) *travelingSalesmanSolution

	switch reductoAlgoName {
	case "standard":
		reductoAlgo = standardReducto
	default:
		panic("not implemeted")
	}

	switch greedyAlgoName {
	case "standard":
		greedyAlgo = standardGreedy
	default:
		panic("not implemented")
	}

	return subTask.Compute(reductoAlgo, greedyAlgo, alpha, betta)
}

type travelingSalesmanSubTask struct {
	countTown          int
	towns              []int
	townsCoord         [][2]float64
	betweenTownsLength [][]float64
}

func MakeTravelingSalesmanSubTask(towns []int, townsCoord [][2]float64, betweenTownsLength [][]float64) *travelingSalesmanSubTask {
	return &travelingSalesmanSubTask{
		countTown:          len(towns),
		towns:              towns,
		townsCoord:         townsCoord,
		betweenTownsLength: betweenTownsLength,
	}
}

func (task *travelingSalesmanSubTask) length(firstTown int, secondTown int) float64 {
	return task.betweenTownsLength[task.towns[firstTown]][task.towns[secondTown]]
}

func (task *travelingSalesmanSubTask) coords(town int) [2]float64 {
	return task.townsCoord[task.towns[town]]
}

func (task *travelingSalesmanSubTask) computeClusterWeightCenter() [2]float64 {
	weightCenter := task.coords(0)
	for i := 1; i < task.countTown; i++ {
		weightCenter = computeWeightCenter(weightCenter, task.coords(i))
	}
	return weightCenter
}

func (task *travelingSalesmanSubTask) computeExternalTask(subTasks []*travelingSalesmanSubTask) []*travelingSalesmanSubTask {
	weightCenters := make([][2]float64, len(subTasks))
	for i, task := range subTasks {
		weightCenters[i] = task.computeClusterWeightCenter()
	}

	betweenClustersLength := computeAllEuclideanDistance(weightCenters)

	clusterOrder := make([]int, len(subTasks))
	for i := 0; i < len(subTasks); i++ {
		clusterOrder[i] = i
	}
	bestClusterOrder := bruteForce(clusterOrder, betweenClustersLength)

	// make result order
	resultOrderSubTask := make([]*travelingSalesmanSubTask, len(subTasks))
	for i := 0; i < len(subTasks); i++ {
		resultOrderSubTask[i] = subTasks[bestClusterOrder[i]]
	}

	return resultOrderSubTask
}

func (task *travelingSalesmanSubTask) bruteForce() *travelingSalesmanSolution {
	sort.Slice(task.towns, func(i int, j int) bool { return task.towns[i] < task.towns[j] })

	bestTownsOrder := bruteForce(task.towns, task.betweenTownsLength)
	return &travelingSalesmanSolution{
		towns: bestTownsOrder,
	}
}

func (task *travelingSalesmanSubTask) CombineSolutions(solutions []*travelingSalesmanSolution) *travelingSalesmanSolution {
	resultTowns := []int{}
	resultTowns = append(resultTowns, solutions[0].towns...)

	for i := 0; i < len(solutions)-1; i++ {
		minLength := math.Inf(1)
		bestJ := 0
		bestK := 0
		for j, town := range resultTowns {
			for k, nextTown := range solutions[i+1].towns {
				if length := task.betweenTownsLength[town][nextTown]; length < minLength {
					minLength = length
					bestJ = j
					bestK = k
				}
			}
		}
		mergeTowns := make([]int, len(resultTowns)+len(solutions[i+1].towns))
		pos := 0
		for j := 0; j <= bestJ; j++ {
			mergeTowns[pos] = resultTowns[j]
			pos++
		}
		for k := bestK; k < len(solutions[i+1].towns); k++ {
			mergeTowns[pos] = solutions[i+1].towns[k]
			pos++
		}
		for k := 0; k < bestK; k++ {
			mergeTowns[pos] = solutions[i+1].towns[k]
			pos++
		}
		for j := bestJ + 1; j < len(resultTowns); j++ {
			mergeTowns[pos] = resultTowns[j]
			pos++
		}
		resultTowns = mergeTowns
	}

	return &travelingSalesmanSolution{
		towns: resultTowns,
	}
}

func (task *travelingSalesmanSubTask) Compute(reducto func(*travelingSalesmanSubTask, int) []*travelingSalesmanSubTask,
	greedy func(*travelingSalesmanSubTask) *travelingSalesmanSolution,
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
			solutions = append(solutions, subTask.Compute(reducto, greedy, alpha, betta))
		} else {
			if subTask.countTown > alphaMax {
				solutions = append(solutions, greedy(subTask))
			} else {
				solutions = append(solutions, subTask.bruteForce())
			}
		}
	}

	solution = task.CombineSolutions(solutions)

	return solution
}
