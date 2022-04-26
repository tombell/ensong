package config

import (
	"fmt"

	"github.com/pelletier/go-toml"
)

type Config struct{}

func Load(path string) (*Config, error) {
	tree, err := toml.LoadFile(path)
	if err != nil {
		return nil, fmt.Errorf("toml load file failed: %w", err)
	}

	cfg := &Config{}
	if err := tree.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("toml unmarshal failed: %w", err)
	}

	return cfg, nil
}
