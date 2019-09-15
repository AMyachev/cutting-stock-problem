package main
import (
	"os"
	"fmt"

	taskP "csp/task"
	"csp/heuristics"
	"csp/cmd"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.WarnLevel)
}

func main() {
	//test1_01_n10_file()
	cmd.Execute()
}


func test1_01_n10_file() {
	fileName := "C:/Users/amyachev/Desktop/UNN/cutting-stock-problem/task_1_02_n10.txt"
	task, err := taskP.MakeOneDimensionalCuttingStockProblemFromFile(fileName)
	if err != nil {
		fmt.Println(err)
	}
	solution := heuristics.GreedyAlgorithmForAscending(task)
	fmt.Println(solution)
	fmt.Println(task.ComputeLowerBound())

	permutation := task.GetAllPiecesByProperty("descending")
	_, solution2 := heuristics.LocalOptimization(task, permutation)
	fmt.Println(solution2)
}
