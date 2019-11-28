package salesman

import (
	"fmt"
	"math"
)

const (
	alphaMax = 10
)

type travelingSalesmanTask struct {
	countTown int
	// x,y coords
	townsCoord  [][2]int
	townsLength [][]float64
}

type travelingSalesmanSolution struct {
	towns []int
}

func MakeTravelingSalesmanTaskTest() *travelingSalesmanTask {
	townsCoord := [][2]int{
		{2, 4},
		{4, 2},
		{5, 4},
		{6, 2},
		{7, 4},
		{5, 6},
		{3, 6},
	}

	townsLength := computeEuclideanDistance(townsCoord)

	return &travelingSalesmanTask{
		countTown:   7,
		townsCoord:  townsCoord,
		townsLength: townsLength,
	}
}

func computeEuclideanDistance(townsCoord [][2]int) [][]float64 {
	countTowns := len(townsCoord)

	// allocation
	townsLength := make([][]float64, countTowns)
	for i := 0; i < countTowns; i++ {
		townsLength[i] = make([]float64, countTowns)
	}

	// initialization
	for i := 0; i < countTowns-1; i++ {
		for j := i + 1; j < countTowns; j++ {
			firstCoord := math.Pow(float64(townsCoord[i][0]-townsCoord[j][0]), 2)
			secondCoord := math.Pow(float64(townsCoord[i][1]-townsCoord[j][1]), 2)
			townsLength[i][j] = math.Sqrt(firstCoord + secondCoord)
			// Euclidean distance is symmetric, so ...
			townsLength[j][i] = townsLength[i][j]
		}
	}

	return townsLength
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
