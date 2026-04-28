package config

import (
	"fmt"
	"sync"

	"github.com/spf13/viper"
)

var (
	cfg  *Config
	once sync.Once
)

type Config struct {
	Server struct {
		Port      int `mapstructure:"port"`
		MaxTables int `mapstructure:"max_tables"`
	} `mapstructure:"server"`
	Game struct {
		InitialPoints int `mapstructure:"initial_points"`
		MinPay       int `mapstructure:"min_pay"`
		MaxPlayers  int `mapstructure:"max_players"`
	} `mapstructure:"game"`
}

func Load() *Config {
	once.Do(func() {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		if err := viper.ReadInConfig(); err != nil {
			panic(fmt.Errorf("failed to read config: %w", err))
		}
		cfg = &Config{}
		if err := viper.Unmarshal(cfg); err != nil {
			panic(fmt.Errorf("failed to unmarshal config: %w", err))
		}
	})
	return cfg
}

func Get() *Config {
	if cfg == nil {
		return Load()
	}
	return cfg
}