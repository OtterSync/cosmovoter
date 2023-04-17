package config

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"golang.org/x/crypto/scrypt"
)

// Config stores the configuration information for each chain
type Config struct {
	Chains []ChainConfig `json:"chains"`
}

// ChainConfig stores the configuration information for a single chain
type ChainConfig struct {
	Name          string `json:"name"`
	Executable    string `json:"executable"`
	ChainID       string `json:"chain_id"`
	RPCNode       string `json:"rpc_node"`
	WalletName    string `json:"wallet_name"`
	WalletPwdHash string `json:"wallet_pwd_hash"`
	GasPrice      uint64 `json:"gas_price"`
}

const (
	saltSize = 32
	keySize  = 32
)

var (
	configFilePath string
)

// SetConfigFilePath sets the path to the config file
func SetConfigFilePath(path string) {
	configFilePath = path
}

// LoadConfig loads the configuration from the file
func LoadConfig() (Config, error) {
	var config Config
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		// Config file does not exist, create a new one
		config.Chains = []ChainConfig{}
		return config, SaveConfig(config)
	}
	configBytes, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return config, fmt.Errorf("failed to read config file: %w", err)
	}
	err = json.Unmarshal(configBytes, &config)
	if err != nil {
		return config, fmt.Errorf("failed to unmarshal config file: %w", err)
	}
	return config, nil
}

// SaveConfig saves the configuration to the file
func SaveConfig(config Config) error {
	configBytes, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config file: %w", err)
	}
	err = ioutil.WriteFile(configFilePath, configBytes, 0644)
	if err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	return nil
}

// EncryptPassword encrypts the wallet password using AES-256 encryption with a random salt
func EncryptPassword(password string) (string, error) {
	salt := make([]byte, saltSize)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	// Derive the encryption key from the password using scrypt
	key, err := scrypt.Key([]byte(password), salt, 1<<14, 8, 1, keySize)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Encrypt the password using AES-256 encryption
	ciphertext := make([]byte, aes.BlockSize+len(password))
	iv := ciphertext[:aes.BlockSize]
	if _, err := rand.Read(iv); err != nil {
		return "", err
	}
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(password))

	// Encode the encrypted password and salt in base64 format
	encrypted := base64.URLEncoding.EncodeToString(append(salt, ciphertext...))
	return encrypted, nil
}

// DecryptPassword decrypts the encrypted wallet password using AES-256 decryption
func DecryptPassword(encrypted string, password string) (string, error) {
	// Decode the encrypted password and salt from base64 format
	encryptedBytes, err := base64.URLEncoding.DecodeString(encrypted)
	if err != nil {
		return "", fmt.Errorf("failed to decode encrypted password: %w", err)
	}

	salt := encryptedBytes[:saltSize]
	ciphertext := encryptedBytes[saltSize:]

	// Derive the encryption key from the password and salt using scrypt
	key, err := scrypt.Key([]byte(password), salt, 1<<14, 8, 1, keySize)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Decrypt the password using AES-256 decryption
	plaintext := make([]byte, len(ciphertext)-aes.BlockSize)
	iv := ciphertext[:aes.BlockSize]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(plaintext, ciphertext[aes.BlockSize:])

	return string(plaintext), nil
}
