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

	task := knapsack.MakeMakeKnapsackProblemFromFile("example-tasks/knapsack/task_3_01_n5.txt")

	start := time.Now()
	//solution, criterion := task.RecursiveSolutionDefaultOrder(task.GetCompanyPerfomance(), true)

	//fmt.Println(solution, criterion)

	//criterion := task.TableSolutionByDefault(task.GetCompanyPerfomance())

	permutation := task.DefaultSort()

	criterion := task.TableSolution(permutation, task.GetCompanyPerfomance())

	fmt.Println(criterion)

	end := time.Now()
	fmt.Println("time: ", end.Sub(start))
}
