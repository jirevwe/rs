package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringVar(&dbPath, "db-path", "data", "the mongodb dp path")
	runCmd.Flags().StringVar(&version, "version", "4.2.21", "specify the mongodb version")
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Configures and runs a mongodb version",
	RunE: func(cmd *cobra.Command, args []string) error {
		homedir, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		binDir := fmt.Sprintf("%s/.rs/%s/bin", homedir, version)
		dataDir := fmt.Sprintf("%s/.rs/%s/data", homedir, version)

		println(binDir, dataDir)

		return nil
	},
}
