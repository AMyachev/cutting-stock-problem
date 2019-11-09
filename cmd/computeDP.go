// compute delivery problem command
package cmd

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

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

	computeDPCmd.Flags().BoolVar(&getDeviation, "deviation", true, "compute deviation relative to default strategies")
}

var deliveryProblemsDir string
var deliveryProblemFile string

var branchingStrategy string
var upperBoundStrategy string
var lowerBoundStrategy string

var getDeviation bool

var computeDPCmd = &cobra.Command{
	Use:   "computeDP",
	Short: "compute delivery problems",
	Long:  `<add description>`,
	Run: func(cmd *cobra.Command, args []string) {
		if deliveryProblemFile != "" {
			deviat := ComputeDPProblem(deliveryProblemFile, branchingStrategy, lowerBoundStrategy, upperBoundStrategy, getDeviation)
			if getDeviation {
				fmt.Printf("deviation: %f", deviat)
			}
		} else {
			files, err := ioutil.ReadDir(deliveryProblemsDir)
			if err != nil {
				log.WithFields(log.Fields{
					"deliveryProblemsDir": deliveryProblemsDir,
				}).Fatalf("computeDPCmd: %s", err)
			}

			var sumDeviat float64 = 0
			for _, file := range files {
				deviat := ComputeDPProblem(filepath.Join(deliveryProblemsDir, file.Name()), branchingStrategy, lowerBoundStrategy, upperBoundStrategy, getDeviation)
				fmt.Printf("deviation: %f\n\n", deviat)
				sumDeviat += deviat
			}

			if getDeviation {
				averageDeviation := sumDeviat / float64(len(files))
				fmt.Printf("averageDeviation: %f", averageDeviation)
			}
		}
	},
}

func ComputeDPProblem(taskFile, branchingStrategy, lowerBoundStrategy, upperBoundStrategy string, getDeviation bool) float64 {
	start := time.Now()
	taskDP := task.MakeDeliveryProblemFromFile(taskFile, branchingStrategy, lowerBoundStrategy, upperBoundStrategy)
	bestVertex, countTraversedVertexes := exact_algorithms.BnB(taskDP)

	end := time.Now()
	elapsed := end.Sub(start)
	fmt.Printf("best Vertex: %v\n countTraversedVertexes: %d\n time: %v\n\n", bestVertex, countTraversedVertexes, elapsed)

	if !getDeviation {
		return 0
	}

	log.WithFields(log.Fields{
		"taskFile":           taskFile,
		"branchingStrategy":  "optimistic",
		"lowerBoundStrategy": "default",
		"upperBoundStrategy": "default",
	}).Infof("ComputeDPProblem: compute default task")
	defaultTaskDP := task.MakeDeliveryProblemFromFile(taskFile, "optimistic", "default", "default")
	_, defaultCountTraversedVertexes := exact_algorithms.BnB(defaultTaskDP)

	return float64(defaultCountTraversedVertexes-countTraversedVertexes) / float64(defaultCountTraversedVertexes)
}
