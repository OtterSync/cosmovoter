package config

import (
	"encoding/json"
	"errors"
	"github.com/OtterSync/cosmovoter/crypto"
	"io/ioutil"
	"os"
	"path/filepath"
)

type ChainConfig struct {
	ChainName      string `json:"chain-name"`
	Executable     string `json:"executable"`
	WalletName     string `json:"wallet-name,omitempty"`
	WalletPassword string `json:"wallet-password,omitempty"`
	Salt           string `json:"salt,omitempty"`
	GasPrice       string `json:"gas-price,omitempty"`
	ChainID        string `json:"chain-id"`
	RPCNode        string `json:"rpc-node"`
}

func GetChainConfig(chainName string) (*ChainConfig, error) {
	configDir := GetConfigDir()
	configFilePath := filepath.Join(configDir, "config.json")

	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		return nil, errors.New("config file does not exist")
	}

	configFile, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return nil, err
	}

	var chainConfigs []ChainConfig
	if err := json.Unmarshal(configFile, &chainConfigs); err != nil {
		return nil, err
	}

	for _, chainConfig := range chainConfigs {
		if chainConfig.ChainName == chainName {
			if chainConfig.WalletPassword != "" && chainConfig.Salt != "" {
				key, _ := crypto.GenerateKey([]byte(chainConfig.WalletPassword), []byte(chainConfig.Salt), 32)
				decryptedPassword, _ := crypto.Decrypt([]byte(chainConfig.WalletPassword), key)
				chainConfig.WalletPassword = string(decryptedPassword)
			}
			return &chainConfig, nil
		}
	}

	return nil, errors.New("chain config not found")
}

func AddChainConfig(chainConfig *ChainConfig) error {
	configDir := GetConfigDir()
	configFilePath := filepath.Join(configDir, "config.json")

	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return err
		}
		if _, err := os.Create(configFilePath); err != nil {
			return err
		}
	}

	configFile, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return err
	}

	var chainConfigs []ChainConfig
	if err := json.Unmarshal(configFile, &chainConfigs); err != nil {
		return err
	}

	for _, existingConfig := range chainConfigs {
		if existingConfig.ChainName == chainConfig.ChainName {
			return errors.New("chain with the same name already exists")
		}
	}

	if chainConfig.WalletPassword != "" {
		key, _ := crypto.GenerateKey([]byte(chainConfig.WalletPassword), []byte(chainConfig.Salt), 32)
		encryptedPassword, _ := crypto.Encrypt([]byte(chainConfig.WalletPassword), key)
		chainConfig.WalletPassword = crypto.HexEncode(encryptedPassword)
	}

	chainConfigs = append(chainConfigs, *chainConfig)

	chainConfigsJSON, err := json.Marshal(chainConfigs)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(configFilePath, chainConfigsJSON, 0644)
	if err != nil {
		return err
	}

	return nil
}

func UpdateChain(chainName, executable, walletName, gasPrice, chainID, rpcNode string, encryptedPassword, salt []byte
