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

	task := knapsack.MakeMakeKnapsackProblemFromFile("example-tasks/knapsack/task_3_10_n1000.txt")

	start := time.Now()
	solution, criterion := task.RecursiveSolutionDefaultOrder(task.GetCompanyPerfomance(), true)

	fmt.Println(solution, criterion)
	end := time.Now()
	fmt.Println("time: ", end.Sub(start))
}
