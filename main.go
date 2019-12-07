package main

import (
	"os"

	"csp/task/travelingSalesman"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.WarnLevel)
}

func main() {
	//cmd.Execute()

	task := travelingSalesman.MakeTravelingSalesmanTaskTest("example-tasks/travelingSalesman/DJ38.txt")
	print(task)
}
