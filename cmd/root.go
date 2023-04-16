package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "cosmovoter",
	Short: "Cosmovoter is a tool to manage Cosmos ecosystem validators",
	Long: `Cosmovoter is a CLI tool designed to help companies running Cosmos ecosystem validators
to easily vote on different proposals on different chains with different executables.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
