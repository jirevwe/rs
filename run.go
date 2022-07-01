package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"time"

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
		mongodDir := fmt.Sprintf("%s/%s", binDir, "mongod")
		mongoshDir := fmt.Sprintf("%s/%s", binDir, "mongo")

		err = os.MkdirAll(dataDir, os.ModePerm)
		if err != nil {
			return err
		}

		mongod := exec.Command(mongodDir, "--dbpath", dataDir, "--replSet", "localhost")
		err = mongod.Start()
		if err != nil {
			return err
		}
		fmt.Printf("Mongo Daemon started Successfully\n\n")

		time.Sleep(1 * time.Second)

		mongoShell := exec.Command(mongoshDir)
		mongoShell.Stdout = os.Stdout

		input := &bytes.Buffer{}
		input.Write([]byte("rs.initiate()\nshow dbs"))
		mongoShell.Stdin = input

		err = mongoShell.Run()
		if err != nil {
			return err
		}

		if err := mongod.Wait(); err != nil {
			fmt.Printf("\n\n[DAEMON:] %s\n\n", err.Error())
			return err
		}

		return nil
	},
}
