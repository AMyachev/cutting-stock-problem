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
	taskDir = "example-tasks"

	rootCmd.AddCommand(computeCmd)
	computeCmd.Flags().StringVar(&taskFile, "task-file", "", "path to task-file (required)")
	computeCmd.Flags().StringVar(&algorithm, "algo", "greedy-for-ascending", "algorithm that would be used for computation")
}

var taskDir string
var taskFile string
var algorithm string

var computeCmd = &cobra.Command{
	Use:   "compute",
    Short: "Computes tasks provided by ",
    Long:  `<add description>`,
    Run: func(cmd *cobra.Command, args []string) {
		if taskFile != "" {
			computeProblem(taskFile, algorithm)
		} else {
			files, err := ioutil.ReadDir(taskDir)
			if err != nil {
				log.WithFields(log.Fields{
					"error": err,
					"taskDir": taskDir,
				}).Fatal("Compute command errored")
			}
		
			for _, file := range files {
				computeProblem(filepath.Join(taskDir, file.Name()), algorithm)
			}
		}
    },
}

func computeProblem(taskFile string, algorithm string) {
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

	fmt.Printf("%s == %d; time - %v\n", taskFile, solution.GetCountUsedMaterials(), elapsed)
	fmt.Println(solution, task.ComputeLowerBound())
}
