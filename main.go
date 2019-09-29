package main
import (
	"os"

	//"csp/cmd"
	"csp/exact_algorithms"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.WarnLevel)
}

func main() {
	//cmd.Execute()
	exact_algorithms.BnB()
}
