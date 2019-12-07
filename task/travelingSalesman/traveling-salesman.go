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

func MakeTravelingSalesmanTaskTest(taskFile string) *travelingSalesmanTask {
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

	betweenTownsLength := computeEuclideanDistance(townsCoord)

	return &travelingSalesmanTask{
		countTown:          dimension,
		townsCoord:         townsCoord,
		betweenTownsLength: betweenTownsLength,
	}
}

func computeEuclideanDistance(townsCoord [][2]float64) [][]float64 {
	countTowns := len(townsCoord)

	// allocation
	betweenTownsLength := make([][]float64, countTowns)
	for i := 0; i < countTowns; i++ {
		betweenTownsLength[i] = make([]float64, countTowns)
	}

	// initialization
	for i := 0; i < countTowns-1; i++ {
		for j := i + 1; j < countTowns; j++ {
			xSqr := math.Pow(townsCoord[i][0]-townsCoord[j][0], 2)
			ySqr := math.Pow(townsCoord[i][1]-townsCoord[j][1], 2)
			betweenTownsLength[i][j] = math.Sqrt(xSqr + ySqr)
			// Euclidean distance is symmetric, so ...
			betweenTownsLength[j][i] = betweenTownsLength[i][j]
		}
	}

	return betweenTownsLength
}

func (task *travelingSalesmanTask) reducto(alpha int) []*travelingSalesmanTask {
	panic("not implemented")
}

func (task *travelingSalesmanTask) computeExternalTask(subTasks []*travelingSalesmanTask) []*travelingSalesmanTask {
	panic("not implemented")
}

func (task *travelingSalesmanTask) CountTown() int {
	return task.countTown
}

func (task *travelingSalesmanTask) Greedy() *travelingSalesmanSolution {
	panic("not implemented")
}

func (task *travelingSalesmanTask) ExhaustiveSearch() *travelingSalesmanSolution {
	panic("not implemented")
}

func (task *travelingSalesmanTask) CombineSolutions(solutions []*travelingSalesmanSolution) *travelingSalesmanSolution {
	panic("not implemented")
}

func (task *travelingSalesmanTask) Compute(alpha int, betta int) (solution *travelingSalesmanSolution) {
	solutions := []*travelingSalesmanSolution{}

	if alpha > alphaMax {
		panic(fmt.Sprintf("compute: alpha [%d] > max alpha [%d]", alpha, alphaMax))
	}

	subTasks := task.reducto(alpha)
	betta--

	orderedSubTasks := task.computeExternalTask(subTasks)

	for _, subTask := range orderedSubTasks {
		if betta != 0 {
			solutions = append(solutions, subTask.Compute(alpha, betta))
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
