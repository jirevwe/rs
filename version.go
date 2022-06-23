package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print out the cli version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s\n", GetVersion())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)

}
