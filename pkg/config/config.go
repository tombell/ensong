package config

import (
	"fmt"

	"github.com/pelletier/go-toml"
)

type Config struct {
	Directories []string `toml:"directories"`

	Backup string `toml:"backup"`

	Metadata struct {
		Names struct {
			Monday    string `toml:"monday"`
			Tuesday   string `toml:"tuesday"`
			Wednesday string `toml:"wednesday"`
			Thursday  string `toml:"thursday"`
			Friday    string `toml:"friday"`
			Saturday  string `toml:"saturday"`
			Sunday    string `toml:"sunday"`
		} `toml:"names"`

		Pictures struct {
			Monday    string `toml:"monday"`
			Tuesday   string `toml:"tuesday"`
			Wednesday string `toml:"wednesday"`
			Thursday  string `toml:"thursday"`
			Friday    string `toml:"friday"`
			Saturday  string `toml:"saturday"`
			Sunday    string `toml:"sunday"`
		} `toml:"pictures"`
	} `toml:"metadata"`

	GitHub struct {
		Token string `toml:"token"`

		Repo struct {
			Owner  string `toml:"owner"`
			Name   string `toml:"name"`
			Branch string `toml:"branch"`
		} `toml:"repo"`

		Committer struct {
			Name  string `toml:"name"`
			Email string `toml:"email"`
		} `toml:"committer"`
	} `toml:"github"`

	Mixcloud struct {
		Token string `toml:"token"`
	} `toml:"mixcloud"`
}

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
