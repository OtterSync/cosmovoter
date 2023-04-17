package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/OtterSync/cosmovoter/config"
	"github.com/spf13/cobra"
)

var addChainCmd = &cobra.Command{
	Use:   "add-chain",
	Short: "Add a new chain to the configuration",
	Long: `Add a new chain to the configuration file, including the chain name, executable, chain ID,
RPC node, wallet name, and gas price. The wallet password will be prompted for securely.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		reader := bufio.NewReader(os.Stdin)

		// Prompt the user for the chain configuration details
		fmt.Print("Enter chain name: ")
		chainName, _ := reader.ReadString('\n')
		fmt.Print("Enter chain executable path: ")
		executablePath, _ := reader.ReadString('\n')
		fmt.Print("Enter chain ID: ")
		chainID, _ := reader.ReadString('\n')
		fmt.Print("Enter chain RPC node: ")
		rpcNode, _ := reader.ReadString('\n')
		fmt.Print("Enter wallet name: ")
		walletName, _ := reader.ReadString('\n')
		fmt.Print("Enter wallet password: ")
		walletPassword, _ := reader.ReadString('\n')
		fmt.Print("Enter gas price: ")
		var gasPrice uint64
		_, err := fmt.Scanln(&gasPrice)
		if err != nil {
			return fmt.Errorf("failed to read gas price: %w", err)
		}

		// Create the new chain configuration
		chainConfig, err := ChainConfigFromFields(
			strings.TrimSpace(chainName),
			strings.TrimSpace(executablePath),
			strings.TrimSpace(chainID),
			strings.TrimSpace(rpcNode),
			strings.TrimSpace(walletName),
			strings.TrimSpace(walletPassword),
			gasPrice,
		)
		if err != nil {
			return err
		}

		// Save the new chain configuration
		err = chainConfig.Save()
		if err != nil {
			return err
		}

		fmt.Printf("Chain %s added to config file\n", chainConfig.Name)
		return nil
	},
}

func ChainConfigFromFields(name, executable, chainID, rpcNode, walletName, walletPassword string, gasPrice uint64) (*config.ChainConfig, error) {
	// Encrypt the wallet password
	walletPwdHash, err := config.EncryptPassword(strings.TrimSpace(walletPassword))
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt wallet password: %w", err)
	}

	// Create the new chain configuration
	chainConfig := &config.ChainConfig{
		Name:          strings.TrimSpace(name),
		Executable:    strings.TrimSpace(executable),
		ChainID:       strings.TrimSpace(chainID),
		RPCNode:       strings.TrimSpace(rpcNode),
		WalletName:    strings.TrimSpace(walletName),
		WalletPwdHash: walletPwdHash,
		GasPrice:      gasPrice,
	}

	return chainConfig, nil
}
