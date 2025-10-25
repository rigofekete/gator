package config

import  (
	"encoding/json"
	"os"
	"fmt"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DBURL 		string  `json:"db_url"`
	CurrentUserName string 	`json:"current_user_name"`
}


func (cfg *Config) SetUser(userName string) error {
	cfg.CurrentUserName = userName
	return write(*cfg) 
}

func GetConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()

	if err != nil {
		fmt.Println("home directory could not be found\n")
		return "", err
	}

	path := filepath.Join(home, ".config", "gator", configFileName)
	return path, nil
}

func Read() (Config, error) {
	fullPath, err := GetConfigFilePath()
	if err != nil {
		fmt.Println("error getting file path")
		return Config{}, err
	}

	file, err := os.Open(fullPath)
	if err != nil {
		fmt.Println("error opening file")
		return Config{}, err
	}
	defer file.Close()

	cfg := Config{}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&cfg)
	if err != nil {
		fmt.Println("error decoding file")
		return Config{}, err
	}

	return cfg, nil
}


func write(cfg Config) error {
	fullPath, err := GetConfigFilePath()
	if err != nil {
		return fmt.Errorf("Error getting file path: %v", err)
	}

	file, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(&cfg)
	if err != nil {
		return fmt.Errorf("error encoding to config file: %v", err)
	}

	return nil
}
