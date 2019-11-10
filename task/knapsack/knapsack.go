package knapsack

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
