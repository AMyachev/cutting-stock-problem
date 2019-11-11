package knapsack

import (
	"io/ioutil"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

type knapsackProblem struct {
	companyPerformance int
	countOrders        int
	complexityOrders   []int
	costOrders         []int
}

func MakeKnapsackProblemTest() *knapsackProblem {
	return &knapsackProblem{
		companyPerformance: 10,
		countOrders:        5,
		complexityOrders:   []int{4, 5, 4, 3, 2},
		costOrders:         []int{2, 3, 2, 3, 1},
	}
}

func MakeMakeKnapsackProblemFromFile(taskFile string) *knapsackProblem {
	content, err := ioutil.ReadFile(taskFile)
	if err != nil {
		log.Fatal(err)
	}

	contentLines := strings.Split(string(content), "\n")

	companyPerformance, err := strconv.Atoi(strings.TrimSpace(contentLines[0]))
	if err != nil {
		log.Fatal(err)
	}

	countOrders, err := strconv.Atoi(strings.TrimSpace(contentLines[1]))
	if err != nil {
		log.Fatal(err)
	}

	complexityOrdersString := strings.Split(strings.TrimSpace(contentLines[2]), " ")
	complexityOrdersInt := make([]int, countOrders)
	for pos, value := range complexityOrdersString {
		complexityOrdersInt[pos], err = strconv.Atoi(value)
		if err != nil {
			log.Fatal(err)
		}
	}

	costOrdersString := strings.Split(strings.TrimSpace(contentLines[3]), " ")
	costOrdersInt := make([]int, countOrders)
	for pos, value := range costOrdersString {
		costOrdersInt[pos], err = strconv.Atoi(value)
		if err != nil {
			log.Fatal(err)
		}
	}

	return &knapsackProblem{
		companyPerformance: companyPerformance,
		countOrders:        countOrders,
		complexityOrders:   complexityOrdersInt,
		costOrders:         costOrdersInt,
	}
}

func (problem *knapsackProblem) RecursiveSolution(permutation []int, remainingPerformance int, doCache bool) (solution []bool, criterion int) {
	if doCache {
		panic("not implemented")
	}

	// dimension ~ k
	dimension := len(permutation)

	if dimension == 0 {
		return []bool{}, 0
	}

	// complexityByDimension ~ a_k+1
	complexityByDimension := problem.complexityOrders[dimension-1]

	// case when the order is not taken
	solution, criterion = problem.RecursiveSolution(permutation[:dimension-1], remainingPerformance, doCache)
	solution = append(solution, false)

	if complexityByDimension <= remainingPerformance {
		// case when the order is taken
		secondSolution, secondCriterion := problem.RecursiveSolution(permutation[:dimension-1], remainingPerformance-complexityByDimension, doCache)
		secondSolution = append(secondSolution, true)
		lastElem := permutation[dimension-1]

		// ~ += c_k+1
		secondCriterion += problem.costOrders[lastElem]

		if secondCriterion >= criterion {
			// second solution is best
			solution = secondSolution
			criterion = secondCriterion
		}
	}

	return solution, criterion
}

func (problem *knapsackProblem) GetCountOrders() int {
	return problem.countOrders
}

func (problem *knapsackProblem) GetCompanyPerfomance() int {
	return problem.companyPerformance
}
