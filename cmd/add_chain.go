package cmd

import (
	"fmt"
	"github.com/OtterSync/cosmovoter/config"
	"github.com/spf13/cobra"
)

var addChainCmd = &cobra.Command{
	Use:   "add-chain",
	Short: "Add a new chain to the configuration file",
	Run: func(cmd *cobra.Command, args []string) {
		chainName, _ := cmd.Flags().GetString("chain-name")
		executable, _ := cmd.Flags().GetString("executable")
		walletName, _ := cmd.Flags().GetString("wallet-name")
		gasPrice, _ := cmd.Flags().GetString("gas-price")
		chainID, _ := cmd.Flags().GetString("chain-id")
		rpcNode, _ := cmd.Flags().GetString("rpc-node")

		if chainName == "" || executable == "" || walletName == "" || gasPrice == "" || chainID == "" || rpcNode == "" {
			fmt.Println("Error: All flags are required.")
			return
		}

		err := config.AddChain(chainName, executable, walletName, gasPrice, chainID, rpcNode)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		} else {
			fmt.Println("Chain added successfully.")
		}
	},
}

func init() {
	addChainCmd.Flags().StringP("chain-name", "n", "", "Name of the chain")
	addChainCmd.Flags().StringP("executable", "e", "", "Chain executable program")
	addChainCmd.Flags().StringP("wallet-name", "w", "", "Name of the wallet")
	addChainCmd.Flags().StringP("gas-price", "g", "", "Gas price")
	addChainCmd.Flags().StringP("chain-id", "i", "", "Chain ID")
	addChainCmd.Flags().StringP("rpc-node", "r", "", "RPC node")
	addChainCmd.MarkFlagRequired("chain-name")
	addChainCmd.MarkFlagRequired("executable")
	addChainCmd.MarkFlagRequired("wallet-name")
	addChainCmd.MarkFlagRequired("gas-price")
	addChainCmd.MarkFlagRequired("chain-id")
	addChainCmd.MarkFlagRequired("rpc-node")
	rootCmd.AddCommand(addChainCmd)
}
