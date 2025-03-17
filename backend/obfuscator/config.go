package obfuscator

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	SingleLineOutput bool `json:"single_line_output"`
	StringLiteral    bool `json:"string_literal"`
	LoopStatement    bool `json:"loop_statement"`
	IfStatement      bool `json:"if_statement"`
	ConstantName     bool `json:"constant_name"`
	VariableName     bool `json:"variable_name"`
	FunctionName     bool `json:"function_name"`
	RemoveComments   bool `json:"remove_comments"`
	path             string
}

func NewConfig() *Config {
	var config Config

	err := config.load()
	if os.IsNotExist(err) {
		if err = config.save(); err != nil {
			log.Fatal(err)
		}
	} else if err != nil {
		log.Fatal(err)
	}

	return &config
}

func (config *Config) load() error {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	appConfigDir := filepath.Join(configDir, "hyperion")
	_, err = os.Stat(appConfigDir)
	if os.IsNotExist(err) {
		if err = os.MkdirAll(appConfigDir, os.ModePerm); err != nil {
			return err
		}
	}

	config.path = filepath.Join(appConfigDir, "config.json")
	if _, err = os.Stat(config.path); err != nil {
		return err
	}

	data, err := os.ReadFile(config.path)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(data, &config); err != nil {
		return err
	}

	return nil
}

func (config *Config) save() error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(config.path, data, os.ModePerm)
}

func (config *Config) Save(c Config) {
	log.Println("Saving obfuscation configuration", c)
}

func (config *Config) GetConfig() Config {
	return *config
}

func (config *Config) GetPath() string {
	return config.path
}
