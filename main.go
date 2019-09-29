package main
import (
	"os"

	"csp/cmd"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func main() {
	//cmd.Execute()
	cmd.ComputeDPProblem("example-tasks/delivery-problem/task_2_10_n50.txt", "optimistic", "default", "default")
}
