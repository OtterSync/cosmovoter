package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "voter-config",
	Short: "Voter-config manages the configuration for the Cosmovoter project",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(addChainCmd)
	rootCmd.AddCommand(updateChainCmd)
	rootCmd.AddCommand(removeChainCmd)
	rootCmd.AddCommand(listCmd)
}
