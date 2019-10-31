package main
import (
	"os"

	"csp/cmd"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.WarnLevel)
}

func main() {
	cmd.Execute()
	//cmd.ComputeDPProblem("example-tasks/delivery-problem/task_2_01_n3.txt", "optimistic", "default", "default-extra-step")
}
