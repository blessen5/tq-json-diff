package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// DefaultConfigFileName is the standard configuration filename.
const DefaultConfigFileName = ".jdiff.json"

// ConfigFile represents the schema of a .jdiff.json configuration file.
type ConfigFile struct {
	Ignore []string `json:"ignore"`
}

// Load reads and parses a JSON configuration file from disk.
func Load(path string) (*ConfigFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read configuration file: %w", err)
	}

	var cfg ConfigFile
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse configuration file: %w", err)
	}

	return &cfg, nil
}

// LoadDefault attempts to load .jdiff.json from the current working directory.
// Returns (cfg, found, err). If the default file does not exist, it returns (nil, false, nil).
func LoadDefault() (*ConfigFile, bool, error) {
	if _, err := os.Stat(DefaultConfigFileName); err != nil {
		if os.IsNotExist(err) {
			return nil, false, nil
		}
		return nil, false, err
	}

	cfg, err := Load(DefaultConfigFileName)
	if err != nil {
		return nil, true, err
	}
	return cfg, true, nil
}

// Merge combines CLI ignore rules and configuration file ignore rules, preserving order and deduplicating.
func Merge(cliRules, cfgRules []string) []string {
	seen := make(map[string]struct{})
	var result []string

	for _, r := range cliRules {
		if _, ok := seen[r]; !ok && r != "" {
			seen[r] = struct{}{}
			result = append(result, r)
		}
	}

	for _, r := range cfgRules {
		if _, ok := seen[r]; !ok && r != "" {
			seen[r] = struct{}{}
			result = append(result, r)
		}
	}

	return result
}
