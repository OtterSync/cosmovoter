package cmd

import (
	"fmt"
	"github.com/OtterSync/cosmovoter/config"
	"github.com/spf13/cobra"
)

var updateChainCmd = &cobra.Command{
	Use:   "update-chain",
	Short: "Update a chain's configuration",
	Run: func(cmd *cobra.Command, args []string) {
		chainName, _ := cmd.Flags().GetString("chain-name")
		executable, _ := cmd.Flags().GetString("executable")
		walletName, _ := cmd.Flags().GetString("wallet
