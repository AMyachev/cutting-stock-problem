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

	task := knapsack.MakeKnapsackProblemTest()

	order := []int{0, 1, 2, 3, 4}
	solution, criterion := task.RecursiveSolution(order, task.GetCompanyPerfomance(), false)

	fmt.Println(solution, criterion)
}
