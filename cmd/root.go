package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	log "github.com/sirupsen/logrus"
)

var rootCmd = &cobra.Command{
	Use:   "csp",
	Short: "csp - the tool for solving the 1D Cutting-Stock problem",
	Long: `csp - the tool for solving the 1D Cutting-Stock problem`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("csp - the tool for solving the 1D Cutting-Stock problem\nfor more information use --help")
	},
  }
  
  func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("rootCmd.Execute() errored")
	}
  }