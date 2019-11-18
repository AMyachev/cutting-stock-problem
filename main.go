package main

import (
	"fmt"
	"os"
	"time"

	"csp/task/knapsack"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.WarnLevel)
}

func main() {
	//cmd.Execute()

	start := time.Now()

	RecursiveSolution()

	end := time.Now()
	fmt.Println("time: ", end.Sub(start))
}

// RecursiveSolution compute with default order
func RecursiveSolution() {
	task := knapsack.MakeMakeKnapsackProblemFromFile("example-tasks/knapsack/task_3_09_n1000.txt")
	solution, criterion := task.RecursiveSolutionDefaultOrder(task.GetCompanyPerfomance(), true)
	fmt.Println(solution, criterion)
}

// TableSolution compute with default order
func TableSolution() {
	task := knapsack.MakeMakeKnapsackProblemFromFile("example-tasks/knapsack/task_3_09_n1000.txt")
	criterion := task.TableSolutionByDefault(task.GetCompanyPerfomance())
	fmt.Println(criterion)
}

// TableSolutionWithDifferentStrategy compute with modificated order (1% object)
func TableSolutionWithDifferentStrategy() {
	task := knapsack.MakeMakeKnapsackProblemFromFile("example-tasks/knapsack/task_3_09_n1000.txt")
	defaultPermutation := task.DefaultSort()
	costPermutation := task.CostSort()

	defaultCriterion := task.TableSolution(defaultPermutation[:10], task.GetCompanyPerfomance())
	costCriterion := task.TableSolution(costPermutation[:10], task.GetCompanyPerfomance())

	fmt.Println("default: ", defaultCriterion, "cost: ", costCriterion)
}
