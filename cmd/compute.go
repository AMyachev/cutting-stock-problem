package cmd

import (
	"fmt"
	"time"
	"io/ioutil"
	"path/filepath"

	"github.com/spf13/cobra"
	log "github.com/sirupsen/logrus"

    "csp/heuristics"
	taskP "csp/task"
	solutionP "csp/solution"
)

func init() {
	taskDir = "example-tasks/1DCSP"

	rootCmd.AddCommand(computeCmd)
	computeCmd.Flags().StringVar(&taskFile, "task-file", "", "path to task-file")
	computeCmd.Flags().StringVar(&algorithm, "algo", "greedy-for-ascending", "algorithm that would be used for computation: {greedy-for-ascending, greedy-for-descending, local-optim-descending, search}")
	computeCmd.Flags().StringVar(&deviation, "deviation", "lower-bound", "the value relative to which the deviation will be calculated: {lower-bound, greedy-value}")
}

var taskDir string
var taskFile string
var algorithm string
var deviation string

var computeCmd = &cobra.Command{
	Use:   "compute",
    Short: "Computes tasks provided by ",
    Long:  `<add description>`,
    Run: func(cmd *cobra.Command, args []string) {
		if taskFile != "" {
			computeProblem(taskFile, algorithm, deviation)
		} else {
			files, err := ioutil.ReadDir(taskDir)
			if err != nil {
				log.WithFields(log.Fields{
					"error": err,
					"taskDir": taskDir,
				}).Fatal("Compute command errored")
			}
		
			for _, file := range files {
				computeProblem(filepath.Join(taskDir, file.Name()), algorithm, deviation)
			}
		}
    },
}

func computeProblem(taskFile string, algorithm string, deviation string) {
	var solution solutionP.Solution

	task, err := taskP.MakeOneDimensionalCuttingStockProblemFromFile(taskFile)
	if err != nil {
		fmt.Println(err)
	}

	log.WithFields(log.Fields{
		"taskFile": taskFile,
		"algorithm": algorithm,
	}).Info("Computing problem ...")

	start := time.Now()

	switch algorithm {
	case "greedy-for-ascending":
		solution = heuristics.GreedyAlgorithmForAscending(task)
		break
	case "greedy-for-descending":
		solution = heuristics.GreedyAlgorithmForDescending(task)
		break
	case "local-optim-descending":
		permutation := task.GetAllPiecesByProperty("descending")
		_, solution = heuristics.LocalOptimization(task, permutation)
		break
	case "search":
		solution = heuristics.Search(task)
		break
	default:
		panic("not implemented")
	}

	end := time.Now()
	elapsed := end.Sub(start)

	criterion := solution.GetCountUsedMaterials()
	var deviationValue float64

	switch deviation {
	case "lower-bound":
		value := task.ComputeLowerBound()
		deviationValue = float64(criterion - value) / float64(value)
	case "greedy-value":
		greedySolution := heuristics.GreedyAlgorithmForAscending(task)
		value := greedySolution.GetCountUsedMaterials()
		deviationValue =  float64(value - criterion) / float64(value)
	default:
		panic("not implemented")
	}

	fmt.Printf("%s \t %s: %d; \t deviation: %f; \t lower bound: %d; \t time: %v;\n", taskFile, algorithm,
				solution.GetCountUsedMaterials(), deviationValue, task.ComputeLowerBound(), elapsed)
}
