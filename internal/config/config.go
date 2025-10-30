package config

import (
	"fmt"
	"log/slog"
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
	configFile, err := getConfigFile()
	if err != nil {
		return nil, fmt.Errorf("error getting config file: %w", err)
	}
	slog.Info("found config file", "file", configFile)

	// TODO: find a way to navigate things like:
	// 		- config file doesn't exist
	//		- config file is unreadable
	//		- config file is broken/invalid JSON

	cfg, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}
	slog.Info("read config file", "file", configFile)

	return &Config{
		Airport: string(cfg),
		Modules: ModuleCfg{
			METAR: true,
		},
	}, nil
}
