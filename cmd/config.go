package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/OtterSync/cosmovoter/internal/cosmos"
	"github.com/OtterSync/cosmovoter/internal/crypto"
)

const configFile = "config.json"

func LoadConfig() ([]cosmos.ChainConfig, error) {
	file, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	var configs []cosmos.ChainConfig
	if err := json.Unmarshal(file, &configs); err != nil {
		return nil, err
	}

	return configs, nil
}

func SaveConfig(configs []cosmos.ChainConfig) error {
	data, err := json.MarshalIndent(configs, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(configFile, data, 0644)
}

func promptUser(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	value, _ := reader.ReadString('\n')
	return strings.TrimSpace(value)
}

func AddChain() error {
	chainConfig := cosmos.ChainConfig{}

	chainConfig.ChainName = promptUser("Enter Chain Name: ")
	chainConfig.Executable = promptUser("Enter Executable Name: ")
	chainConfig.WalletName = promptUser("Enter Wallet Name: ")

	password := promptUser("Enter Wallet Password: ")
	salt, err := crypto.GenerateSalt()
	if err != nil {
		return err
	}
	encryptedPassword, err := crypto.EncryptWithPassword(password, password, salt)
	if err != nil {
		return err
	}

	chainConfig.Password = string(encryptedPassword)
	chainConfig.Salt = string(salt)

	chainConfig.GasPrices = promptUser("Enter Gas Prices: ")

	// Fetch Chain ID and RPC node from directory lookup
	chainID, rpcNode, err := cosmos.FetchChainDetails(chainConfig.ChainName)
	if err != nil {
		return err
	}
	chainConfig.ChainID = chainID
	chainConfig.RPCNode = rpcNode

	return AddChainConfig(chainConfig)
}

func AddChainConfig(chainConfig cosmos.ChainConfig) error {
	configs, err := LoadConfig()
	if err != nil {
		return err
	}

	configs = append(configs, chainConfig)
	return SaveConfig(configs)
}

func ListChains() error {
	configs, err := LoadConfig()
	if err != nil {
		return err
	}

	if len(configs) == 0 {
		fmt.Println("No chains configured.")
		return nil
	}

	fmt.Println("Configured chains:")
	for i, config := range configs {
		fmt.Printf("%d. Chain Name: %s, Executable: %s, Wallet: %s, Chain ID: %s, RPC Node: %s\n",
			i+1, config.ChainName, config.Executable, config.WalletName, config.ChainID, config.RPCNode)
	}

	return nil
}

func RemoveChain() error {
	if err := ListChains(); err != nil {
		return err
	}

	chainNumber := promptUser("Enter the number of the chain to remove: ")
	index, err := strconv.Atoi(chainNumber)
	if err != nil {
		return fmt.Errorf("Invalid chain number")
	}

	configs, err := LoadConfig()
	if err != nil {
		return err
	}

	if index < 1 || index > len(configs) {
		return fmt.Errorf("Chain number out of range")
	}

	confirm := promptUser(fmt.Sprintf("Are you sure you want to delete the chain with the name %s? (y/n): ", configs[index-1].ChainName))
	if strings.ToLower(confirm) != "y" {
		fmt.Println("Chain removal cancelled.")
		return nil
	}

	return RemoveChainByIndex(index - 1)
}

func RemoveChainByIndex(index int) error {
	configs, err := LoadConfig()
	if err != nil {
		return err
	}

	configs = append(configs[:index], configs[index+1:]...)
	return SaveConfig(configs)
}
