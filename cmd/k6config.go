package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// structure of the config file with the fields that are used by k6exec
type k6configFile struct {
	Collectors struct {
		Cloud struct {
			Token string `json:"token"`
		} `json:"cloud"`
	} `json:"collectors"`
}

// loadConfig loads the k6 config file from the given path or the default locations.
// if using the default locations, if the file does not exist or can't be read, it returns an empty config.
// default locations are k6/config.json and loadimpact/k6/config.json under the user config directory.
func loadConfig(configPath string) (k6configFile, error) {
	if configPath != "" {
		return readConfig(configPath)
	}

	var config k6configFile
	homeDir, err := os.UserConfigDir() //nolint:forbidigo
	if err != nil {
		return config, fmt.Errorf("failed to get user config directory: %w", err)
	}
	for _, location := range []string{"", "loadimpact"} {
		configPath = filepath.Join(homeDir, location, "k6", "config.json")
		config, err = readConfig(configPath)
		if err == nil {
			break
		}
	}

	// if using default locations we don't return errors if we can't read the file
	return config, nil
}

func readConfig(configPath string) (k6configFile, error) {
	buffer, err := os.ReadFile(configPath) //nolint:forbidigo,gosec
	if err != nil {
		return k6configFile{}, fmt.Errorf("failed to read config file %q: %w", configPath, err)
	}

	var config k6configFile
	err = json.Unmarshal(buffer, &config)
	if err != nil {
		// ensure we return an empty config if the file is not valid
		return k6configFile{}, fmt.Errorf("failed to parse config file %q: %w", configPath, err)
	}

	return config, nil
}
