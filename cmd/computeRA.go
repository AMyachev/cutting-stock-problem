package cmd

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"csp/task/resourceAllocation"
)

func init() {
	resourceAllocationProblemsDir = "example-tasks/resourceAllocation"

	rootCmd.AddCommand(computeRACmd)

	computeRACmd.Flags().StringVar(&resourceAllocationProblemFile, "task-file", "", "path to file with traveling salesman problem")
}

var resourceAllocationProblemsDir string
var resourceAllocationProblemFile string

var computeRACmd = &cobra.Command{
	Use:   "computeRA",
	Short: "compute resource allocation problems",
	Long:  `<add description>`,
	Run: func(cmd *cobra.Command, args []string) {
		if resourceAllocationProblemFile != "" {
			ComputeRAProblem(resourceAllocationProblemFile)
		} else {
			files, err := ioutil.ReadDir(resourceAllocationProblemsDir)
			if err != nil {
				log.WithFields(log.Fields{
					"resourceAllocationProblemsDir": resourceAllocationProblemsDir,
				}).Fatalf("computeRACmd: %s", err)
			}

			start := time.Now()
			for _, file := range files {
				ComputeRAProblem(filepath.Join(resourceAllocationProblemsDir, file.Name()))
			}
			end := time.Now()
			fmt.Println("alltime: ", end.Sub(start))
		}
	},
}

func ComputeRAProblem(taskFile string) int {
	task := resourceAllocation.MakeResourceAllocationTaskFromFile(taskFile)
	solution := task.Compute()

	_, taskFileName := filepath.Split(taskFile)
	taskFileName = strings.TrimSuffix(taskFileName, ".txt")

	fmt.Println("flow: ", solution)

	return solution
}
