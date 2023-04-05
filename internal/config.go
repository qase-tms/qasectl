package internal

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const configFilename = ".qaseconfig.json"

// permissions so only user allowed to work with Config
const configPermissions = 0600

type Config struct {
	Token       string `json:"token"`
	ProjectCode string `json:"projectCode"`
	RunId       int    `json:"runId"`
}

func UpdateConfig(apply func(Config) Config) error {
	cfg, err := getConfig()
	// config might not exist - it's not an error then
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	cfg = apply(cfg)

	return saveConfig(cfg)
}

func getConfig() (Config, error) {
	path, err := configPath()
	if err != nil {
		return Config{}, err
	}

	bytes, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	var cfg Config
	err = json.Unmarshal(bytes, &cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func saveConfig(cfg Config) error {
	path, err := configPath()
	if err != nil {
		return err
	}

	bytes, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	// or create
	fd, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, configPermissions)
	if err != nil {
		return err
	}
	defer fd.Close()

	_, err = fd.Write(bytes)
	return err
}

func configPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	path := filepath.Join(homeDir, configFilename)
	return path, nil
}
