package cmd

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"csp/task/travelingSalesman"
)

func init() {
	travelingSalesmanProblemsDir = "example-tasks/travelingSalesman"

	rootCmd.AddCommand(computeTSPCmd)

	computeTSPCmd.Flags().IntVar(&alpha, "count-cluster", 10, "count clusters per step")
	computeTSPCmd.Flags().IntVar(&betta, "depth", 3, "count os using reduction")
	computeTSPCmd.Flags().StringVar(&travelingSalesmanProblemFile, "task-file", "", "path to file with traveling salesman problem")
	computeTSPCmd.Flags().BoolVar(&getOptimDeviation, "deviation", true, "compute deviation relative to optimum value")
}

var travelingSalesmanProblemsDir string
var travelingSalesmanProblemFile string
var getOptimDeviation bool

var alpha int
var betta int

var computeTSPCmd = &cobra.Command{
	Use:   "computeTSP",
	Short: "compute traveling salesman problems",
	Long:  `<add description>`,
	Run: func(cmd *cobra.Command, args []string) {
		if travelingSalesmanProblemFile != "" {
			deviat := ComputeTSProblem(travelingSalesmanProblemFile, getOptimDeviation)
			if getOptimDeviation {
				fmt.Printf("deviation: %f", deviat)
			}
		} else {
			files, err := ioutil.ReadDir(travelingSalesmanProblemsDir)
			if err != nil {
				log.WithFields(log.Fields{
					"travelingSalesmanProblemsDir": travelingSalesmanProblemsDir,
				}).Fatalf("computeTSPCmd: %s", err)
			}

			var sumDeviat float64 = 0
			start := time.Now()
			for _, file := range files {
				deviat := ComputeTSProblem(filepath.Join(travelingSalesmanProblemsDir, file.Name()), getOptimDeviation)
				fmt.Printf("deviation: %f\n\n", deviat)
				sumDeviat += deviat
			}

			if getOptimDeviation {
				averageDeviation := sumDeviat / float64(len(files))
				fmt.Printf("averageDeviation: %f\n", averageDeviation)
			}
			end := time.Now()
			fmt.Println("alltime: ", end.Sub(start))
		}
	},
}

func ComputeTSProblem(taskFile string, getDeviation bool) float64 {
	task := travelingSalesman.MakeTravelingSalesmanTask(taskFile)
	solution := task.Compute("standard", "standard", alpha, betta)

	_, taskFileName := filepath.Split(taskFile)
	taskFileName = strings.TrimSuffix(taskFileName, ".txt")

	critValue := task.Criterion(solution)
	fmt.Println("found: ", critValue)
	fmt.Println("best: ", travelingSalesman.BestMinLength[taskFileName])
	return (critValue - travelingSalesman.BestMinLength[taskFileName]) / travelingSalesman.BestMinLength[taskFileName]
}
