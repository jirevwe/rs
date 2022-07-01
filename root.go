package main

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "rs",
	Short: "rs is a zero config mongodb replica set runner",
	Example: `  rs download
  rs download 4.2.0

  rs run
  rs run 4.2.0 
	`,
	Long: `rs is a zero config mongodb replica set runner. It downloads mongodb and runs it as a replica set`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

}
