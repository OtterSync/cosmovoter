package utils

import (
	"os"
	"path/filepath"
)

func GetConfigDir() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		configDir = "."
	}
	return filepath.Join(configDir, "cosmovoter")
}
