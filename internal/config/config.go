package config

import (
	"encoding/json"
	"fmt"
	"os"
)



type Config struct {
	Db_url string `json:"db_url"`
	Current_user_name string `json:"current_user_name"`
}

func (c Config) SetUser(username string) error {
	c.Current_user_name = username;
	return write(c)
}

const configFileName = ".gatorconfig.json"

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir();

	if err != nil {
		return "", fmt.Errorf("unable to get home dir %v", err)
	}

	entries, err := os.ReadDir(homeDir)

	if err != nil {
		return "", fmt.Errorf("unable to read home dir %v", err)
	}

	for _, entry := range entries {
		if entry.Name() != configFileName {
			continue;
		}

		fileInfo, err := entry.Info()

		if err != nil {
			return "", fmt.Errorf("unable to get file info for config file %v", err)
		}

		configPath := homeDir + "/" + fileInfo.Name();
		
		return configPath, nil;
	}

	return "", fmt.Errorf("unable to find config file");
}

func Read() (Config, error) {
	configFilePath, err := getConfigFilePath();

	if err != nil {
		return Config{}, err;
	}

	fileContents, err := os.ReadFile(configFilePath);

	if err != nil {
		return Config{}, fmt.Errorf("unable to read file contents %v", err)
	}

	var cnfg Config;
	if err := json.Unmarshal(fileContents, &cnfg); err != nil {
		return Config{}, fmt.Errorf("unable to parse config file contents %v", err)
	}

	return cnfg, nil;
}

func write(cnfg Config) error {
	configFilePath, err := getConfigFilePath();

	if err != nil {
		return err;
	}

	jsonData, err := json.Marshal(cnfg)

    if err != nil {
        return err;
    }

	err = os.WriteFile(configFilePath, jsonData, os.ModeDevice)

	if err != nil {
		return err;
	}

	return nil;
}