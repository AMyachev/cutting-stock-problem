package cmd

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"csp/task/resourceAllocation"
)

func init() {
	resourceAllocationProblemsDir = "example-tasks/resourceAllocation"

	rootCmd.AddCommand(computeRACmd)

	computeRACmd.Flags().StringVar(&resourceAllocationProblemFile, "task-file", "", "path to file with traveling salesman problem")
	computeRACmd.Flags().StringVar(&modification, "mod", "default", "task that will be computed {default, warehouse, minVolumeWarehouse, minCountWareHouseDefault, minCountWareHouseImprov}")
}

var resourceAllocationProblemsDir string
var resourceAllocationProblemFile string

var modification string

var computeRACmd = &cobra.Command{
	Use:   "computeRA",
	Short: "compute resource allocation problems",
	Long:  `<add description>`,
	Run: func(cmd *cobra.Command, args []string) {
		if resourceAllocationProblemFile != "" {
			ComputeRAProblem(resourceAllocationProblemFile, modification)
		} else {
			files, err := ioutil.ReadDir(resourceAllocationProblemsDir)
			if err != nil {
				log.WithFields(log.Fields{
					"resourceAllocationProblemsDir": resourceAllocationProblemsDir,
				}).Fatalf("computeRACmd: %s", err)
			}

			start := time.Now()
			for _, file := range files {
				fmt.Println(file.Name() + ":")
				ComputeRAProblem(filepath.Join(resourceAllocationProblemsDir, file.Name()), modification)
				fmt.Println("")
			}
			end := time.Now()
			fmt.Println("alltime: ", end.Sub(start))
		}
	},
}

func ComputeRAProblem(taskFile string, modification string) int {
	task := resourceAllocation.MakeResourceAllocationTaskFromFile(taskFile)
	solution := task.Compute(modification)

	fmt.Println("flow: ", solution)

	return solution
}
