package configuration

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type KeycloakConfig struct {
	BaseUri  string `json:"baseUri"`
	ClientID string `json:"clientId"`
}

type ReconmapApiConfig struct {
	BaseUri string `json:"baseUri"`
}

type Config struct {
	KeycloakConfig    `json:"keycloak"`
	ReconmapApiConfig `json:"reconmapApi"`
}

const configFileName = "config.json"

func NewConfig() Config {
	return Config{
		KeycloakConfig: KeycloakConfig{
			BaseUri:  "http://localhost:8080/realms/reconmap",
			ClientID: "web-client",
		},
		ReconmapApiConfig: ReconmapApiConfig{
			BaseUri: "http://localhost:5510",
		},
	}
}

func GetReconmapConfigDirectory() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".reconmap"), nil

}

func SaveConfig(config Config) (string, error) {
	var reconmapConfigDir, err = GetReconmapConfigDirectory()

	if _, err := os.Stat(reconmapConfigDir); os.IsNotExist(err) {
		if err := os.MkdirAll(reconmapConfigDir, 0750); err != nil {
			return "", err
		}
	}

	jsondata, _ := json.MarshalIndent(config, "", " ")

	filepath := filepath.Join(reconmapConfigDir, configFileName)
	err = os.WriteFile(filepath, jsondata, 0400)

	return filepath, err
}

func ReadConfig() (*Config, error) {
	var reconmapConfigDir, err = GetReconmapConfigDirectory()
	if err != nil {
		return nil, err
	}
	path := filepath.Join(reconmapConfigDir, configFileName)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, err
	}

	jsonFile, err := os.Open(filepath.Clean(path))
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := jsonFile.Close(); err != nil {
			fmt.Printf("Error closing file: %s\n", err)
		}
	}()

	bytes, _ := io.ReadAll(jsonFile)

	config := Config{}
	err = json.Unmarshal(bytes, &config)

	return &config, nil
}

func HasConfig() bool {
	var reconmapConfigDir, err = GetReconmapConfigDirectory()
	if err != nil {
		return false
	}
	path := filepath.Join(reconmapConfigDir, configFileName)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	return true
}
