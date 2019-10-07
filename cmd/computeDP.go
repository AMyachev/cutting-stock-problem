// compute delivery problem command
package cmd

import (
	"fmt"
	"time"
	"io/ioutil"
	"path/filepath"

	"github.com/spf13/cobra"
	log "github.com/sirupsen/logrus"

	"csp/exact_algorithms"
	"csp/task"
)

func init() {
	deliveryProblemsDir = "example-tasks/delivery-problem"

	rootCmd.AddCommand(computeDPCmd)

	computeDPCmd.Flags().StringVar(&deliveryProblemFile, "task-file", "", "path to file with delivery problem")
	computeDPCmd.Flags().StringVar(&branchingStrategy, "branching", "optimistic", "branching strategy for branch and bound method: {}")
	computeDPCmd.Flags().StringVar(&upperBoundStrategy, "upper-bound", "default", "strategy of finding upper bounds: {}")
	computeDPCmd.Flags().StringVar(&lowerBoundStrategy, "lower-bound", "default", "strategy of finding lower bounds: {}")
}

var deliveryProblemsDir string
var deliveryProblemFile string

var branchingStrategy string
var upperBoundStrategy string
var lowerBoundStrategy string

var computeDPCmd = &cobra.Command{
	Use:   "computeDP",
    Short: "compute delivery problems",
    Long:  `<add description>`,
    Run: func(cmd *cobra.Command, args []string) {
		if deliveryProblemFile != "" {
			ComputeDPProblem(deliveryProblemFile, branchingStrategy, lowerBoundStrategy, upperBoundStrategy)
		} else {
			files, err := ioutil.ReadDir(deliveryProblemsDir)
			if err != nil {
				log.WithFields(log.Fields{
					"deliveryProblemsDir": deliveryProblemsDir,
				}).Fatalf("computeDPCmd: %s", err)
			}
		
			for _, file := range files {
				ComputeDPProblem(filepath.Join(deliveryProblemsDir, file.Name()), branchingStrategy, lowerBoundStrategy, upperBoundStrategy)
			}
		}
    },
}

func ComputeDPProblem(taskFile, branchingStrategy, lowerBoundStrategy, upperBoundStrategy string) {
	start := time.Now()
	taskDP := task.MakeDeliveryProblemFromFile(taskFile, branchingStrategy, lowerBoundStrategy, upperBoundStrategy)
	bestVertex, countTraversedVertexes := exact_algorithms.BnB(taskDP)

	end := time.Now()
	elapsed := end.Sub(start)
	fmt.Printf("best Vertex: %v\n countTraversedVertexes: %d\n time: %v\n\n", bestVertex, countTraversedVertexes, elapsed)
}
