package main
import (
	"os"
	"fmt"

	taskP "csp/task"
	"csp/heuristics"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func main() {
	fmt.Println("there1")
	test1_01_n10()
}

func test1_01_n10() {
	materialLength := 18
	piecesLength := []int{10, 8, 7, 8, 3, 5, 12, 5, 8, 4}

	task, err := taskP.MakeOneDimensionalCuttingStockProblem(materialLength, piecesLength)
	if err != nil {
		log.WithFields(log.Fields{
			"materialLength": materialLength,
			"piecesLength": piecesLength,
		}).Fatal("MakeOneDimensionalCuttingStockProblem error")
	}

	solution := heuristics.GreedyAlgorithmByAscending(task)
	fmt.Println(solution)
}
