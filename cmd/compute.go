package cmd

import (
    "fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/spf13/cobra"
	log "github.com/sirupsen/logrus"

    "csp/heuristics"
    taskP "csp/task"
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
	task, err := taskP.MakeOneDimensionalCuttingStockProblemFromFile(taskFile)
	if err != nil {
		fmt.Println(err)
	}

	switch algorithm {
	case "greedy-for-ascending":
		solution := heuristics.GreedyAlgorithmForAscending(task)
		fmt.Printf("%s == %d\n", taskFile, solution.GetCountUsedMaterials())
		break
	case "greedy-for-descending":
		solution := heuristics.GreedyAlgorithmForDescending(task)
		fmt.Printf("%s == %d\n", taskFile, solution.GetCountUsedMaterials())
		break
	case "local-optim-descending":
		permutation := task.GetAllPiecesByProperty("descending")
		_, solution := heuristics.LocalOptimization(task, permutation)
		fmt.Printf("%s == %d\n", taskFile, solution.GetCountUsedMaterials())
		break
	default:
		panic("not implemented")
	}
}
