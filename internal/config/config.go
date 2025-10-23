package config

import (
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	Airport string    `json:"airport"`
	Modules ModuleCfg `json:"modules"`
}

type ModuleCfg struct {
	METAR  bool `json:"metar"`
	TAF    bool `json:"taf"`
	AFD    bool `json:"discussion"`
	AIRMET bool `json:"airmet"`
	PIREP  bool `json:"pirep"`
}

func getConfigFile() (string, error) {
	configDir := os.Getenv("XDG_CONFIG_HOME")

	if configDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("unable to determine home dir: %w", err)
		}

		configDir = filepath.Join(homeDir, ".config")
	}

	configFile := filepath.Join(configDir, "pilot-bar/config.json")
	return configFile, nil
}

// Load - satisfied LSP noise
func Load() (*Config, error) {
	// config file doesn't exist; build it

	// config file is unreadable

	// config file is broken/invalid in some way

}
