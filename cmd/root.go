package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Long: "CLI for the supermarket API. From here we can start up a new daemon or use the client(s) to interact with it",
}

func init() {
	rootCmd.AddCommand(daemonCmd)
	rootCmd.AddCommand(produceClientCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
