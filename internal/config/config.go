package config

import (
	"github.com/BurntSushi/toml"
)

// Chat ...
type Chat struct {
	Name           string
	Type           string
	ChatID         string `toml:"chat_id"`
	Token          string
	IgnoreAccounts []string `toml:"ignore_accounts"`
}

// Config ...
type Config struct {
	Interval int `toml:"update_interval"`
	Src      map[string]Chat
	Dst      map[string]Chat
}

// NewConfig ...
func NewConfig(file string) (*Config, error) {
	var config Config

	if _, err := toml.DecodeFile(file, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
