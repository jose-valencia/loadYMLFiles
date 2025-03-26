package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config structure to map YAML content
type Config struct {
	Server struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"server"`
	Database struct {
		User     string `yaml:"username"`
		Password string `yaml:"password"`
		Name     string `yaml:"name"`
	} `yaml:"database"`
	Client struct {
		Name string `yaml:"name"`
		Id   string `yaml:"id"`
	} `yaml:"client"`
}

// LoadYAML loads a single YAML file into the config struct
func LoadYAML(filename string, config *Config) error {
	if filename == "" {
		return fmt.Errorf("filename cannot be empty")
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(data, config); err != nil {
		return fmt.Errorf("failed to unmarshal YAML from file %s: %w", filename, err)
	}
	return err
}

// LoadAllYAMLFiles loads all YAML files in a given folder
func LoadAllYAMLFiles(folder string) ([]Config, error) {
	var configs []Config

	// Walk through all files in the folder
	err := filepath.WalkDir(folder, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// Process only .yaml or .yml files
		if !d.IsDir() && (strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml")) {
			var config Config
			if err := LoadYAML(path, &config); err != nil {
				fmt.Printf("Error loading YAML file %s: %v\n", path, err)
			} else {
				configs = append(configs, config)
				fmt.Printf("Loaded config from: %s\n", path)
			}
		}
		return nil
	})

	return configs, err
}

func main() {
	folder := "config" // config folder path

	configs, err := LoadAllYAMLFiles(folder)
	if err != nil {
		log.Fatalf("Error loading YAML files: %v", err)
	}

	// Print all loaded configurations
	for i, config := range configs {
		fmt.Printf("---------------------------------------------\n")
		fmt.Printf("id %d: %+v\n", i+1, config.Client.Id)
		fmt.Printf("name %d: %+v\n", i+1, config.Client.Name)
		fmt.Printf("---------------------------------------------\n")
	}
}
