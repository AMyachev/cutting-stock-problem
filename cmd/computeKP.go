package cmd

import (
	"fmt"
	"io/ioutil"
	"math"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"csp/task/knapsack"
)

func init() {
	knapsackProblemsDir = "example-tasks/knapsack"

	rootCmd.AddCommand(computeKPCmd)

	computeKPCmd.Flags().StringVar(&knapsackProblemFile, "task-file", "", "path to file with knapsack problem")
	computeKPCmd.Flags().StringVar(&modification, "mod", "default", "task that will be computed {default, defaultSort, costSort}")
}

var knapsackProblemsDir string
var knapsackProblemFile string

var knapsackModification string

var computeKPCmd = &cobra.Command{
	Use:   "computeKP",
	Short: "compute knapsack problems",
	Long:  `<add description>`,
	Run: func(cmd *cobra.Command, args []string) {
		if knapsackProblemFile != "" {
			ComputeKProblem(knapsackProblemFile, modification)
		} else {
			files, err := ioutil.ReadDir(knapsackProblemsDir)
			if err != nil {
				log.WithFields(log.Fields{
					"knapsackProblemsDir": knapsackProblemsDir,
				}).Fatalf("computeKPCmd: %s", err)
			}

			start := time.Now()
			for _, file := range files {
				fmt.Println(file.Name() + ":")
				ComputeKProblem(filepath.Join(knapsackProblemsDir, file.Name()), modification)
				fmt.Println("")
			}
			end := time.Now()
			fmt.Println("alltime: ", end.Sub(start))
		}
	},
}

func computeLength(length int, ratio float64) int {
	if newLength := int(math.Round(float64(length) * ratio)); newLength != 0 {
		return newLength
	} else {
		return 1
	}
}

func ComputeKProblem(taskFile string, modification string) {
	permutation := []int{}
	task := knapsack.MakeKnapsackProblemFromFile(taskFile)
	switch modification {
	case "defaultSort":
		permutation = task.DefaultSort()
	case "costSort":
		permutation = task.CostSort()
	default:
		permutation = task.DefaultSort()
		_, crit := task.RecursiveSolution(permutation, task.GetCompanyPerfomance(), true)
		fmt.Println("crit: ", crit)
		return
	}

	length := len(permutation)
	fullCrit := task.TableSolution(permutation, task.GetCompanyPerfomance())
	crit10 := task.TableSolution(permutation[:computeLength(length, 0.1)], task.GetCompanyPerfomance())
	crit30 := task.TableSolution(permutation[:computeLength(length, 0.3)], task.GetCompanyPerfomance())
	crit50 := task.TableSolution(permutation[:computeLength(length, 0.5)], task.GetCompanyPerfomance())
	fmt.Println(fullCrit, crit10, crit30, crit50)
	deviat10 := float64(fullCrit-crit10) / float64(fullCrit)
	deviat30 := float64(fullCrit-crit30) / float64(fullCrit)
	deviat50 := float64(fullCrit-crit50) / float64(fullCrit)
	fmt.Printf("10%%: %f\t30%%: %f\t50%%: %f\n", deviat10, deviat30, deviat50)
}
