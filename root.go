package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "rs",
	Short: "rs is a zero config mongodb replica set runner",
	Long:  `rs is a zero config mongodb replica set runner; it downloads mongodb and runs it as a replica set`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func init() {

}