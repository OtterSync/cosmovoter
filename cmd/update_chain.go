package cmd

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"github.com/OtterSync/cosmovoter/config"
	"github.com/OtterSync/cosmovoter/crypto"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var updateChainCmd = &cobra.Command{
	Use:   "update-chain",
	Short: "Update an existing chain configuration in the configuration file",
	Run: func(cmd *cobra.Command, args []string) {
		chainName, _ := cmd.Flags().GetString("chain-name")
		executable, _ := cmd.Flags().GetString("executable")
		walletName, _ := cmd.Flags().GetString("wallet-name")
		gasPrice, _ := cmd.Flags().GetString("gas-price")
		chainID, _ := cmd.Flags().GetString("chain-id")
		rpcNode, _ := cmd.Flags().GetString("rpc-node")

		if chainName == "" {
			fmt.Println("Error: The chain-name flag is required.")
			return
		}

		var encryptedPassword []byte
		var salt []byte

		if walletName != "" {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Enter new wallet password: ")
			password, _ := reader.ReadString('\n')
			password = strings.TrimSpace(password)

			// Generate encryption key and salt
			salt, _ = crypto.GenerateSalt()
			key, _ := crypto.GenerateKey([]byte(password), salt, 32)

			// Encrypt the wallet password
			encryptedPassword, err = crypto.Encrypt([]byte(password), key)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
		}

		// Update the chain configuration
		err := config.UpdateChain(chainName, executable, walletName, gasPrice, chainID, rpcNode, encryptedPassword, salt)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		} else {
			fmt.Println("Chain updated successfully.")
		}
	},
}

func init() {
	updateChainCmd.Flags().StringP("chain-name", "n", "", "Name of the chain to update (required)")
	updateChainCmd.Flags().StringP("executable", "e", "", "Chain executable program")
	updateChainCmd.Flags().StringP("wallet-name", "w", "", "Name of the wallet (optional)")
	updateChainCmd.Flags().StringP("gas-price", "g", "", "Gas price (optional)")
	updateChainCmd.Flags().StringP("chain-id", "i", "", "Chain ID")
	updateChainCmd.Flags().StringP("rpc-node", "r", "", "RPC node")
	rootCmd.AddCommand(updateChainCmd)
}
