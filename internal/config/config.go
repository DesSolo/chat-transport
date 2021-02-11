package config

import (
	"time"

	"github.com/BurntSushi/toml"
)

// Chat ...
type Chat struct {
	Name           string
	Type           string
	Template       string
	ChatID         string `toml:"chat_id"`
	Token          string
	IgnoreAccounts []string `toml:"ignore_accounts"`
}

// Config ...
type Config struct {
	Interval time.Duration `toml:"update_interval"`
	Template string
	Src      map[string]*Chat
	Dst      map[string]*Chat
}

// NewConfig ...
func NewConfig(file string) (*Config, error) {
	var cfg Config

	if _, err := toml.DecodeFile(file, &cfg); err != nil {
		return nil, err
	}

	setDefault(&cfg)

	return &cfg, nil
}

func setDefault(c *Config) {
	setD := func(s map[string]*Chat) {
		for name, chat := range s {
			if chat.Name == "" {
				chat.Name = name
			}
			if chat.Template == "" {
				chat.Template = c.Template
			}
		}
	}

	setD(c.Src)
	setD(c.Dst)
}
