package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() Config {
	configfilepath, err := getConfigfilePath()
	if err != nil {
		return Config{}
	}

	content, err := os.ReadFile(configfilepath)
	if err != nil {
		return Config{}
	}

	var config Config
	if err = json.Unmarshal(content, &config); err != nil {
		return Config{}
	}

	return config
}

func (c *Config) SetUser(name string) error {
	c.CurrentUserName = name
	res, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	path, err := getConfigfilePath()
	if err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(res)

	if err == nil {
		return err
	}
	return nil
}

func getConfigfilePath() (string, error) {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v/%v", homePath, configFileName), nil
}
