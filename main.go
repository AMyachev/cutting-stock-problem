package main

import (
	"fmt"
	"os"

	"csp/task/knapsack"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.WarnLevel)
}

func main() {
	//cmd.Execute()

	task := knapsack.MakeMakeKnapsackProblemFromFile("example-tasks/knapsack/task_3_02_n5.txt")

	solution, criterion := task.RecursiveSolutionDefaultOrder(task.GetCompanyPerfomance(), false)

	fmt.Println(solution, criterion)
}
