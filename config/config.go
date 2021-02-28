package config

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Config holds configuration of the database and the server.
// Another change
type Config struct {
	DBDriver      string `mapstructure:"DB_DRIVER"`
	DBSource      string `mapstructure:"DB_SOURCE"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
}

// LoadConfig get Config from file, environment variables and actively
// check for updates.
func LoadConfig(path string) (*Config, chan struct{}, error) {
	viper.SetDefault("DB_DRIVER", "postgres")
	viper.SetDefault("DB_SOURCE", "postgresql://postgres:postgres@localhost:5432/simple_bank?sslmode=disable")
	viper.SetDefault("SERVER_ADDRESS", "0.0.0.0:8080")

	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)

	if err := viper.ReadInConfig(); err != nil {
		return nil, nil, fmt.Errorf("failed to read configuration file: %w", err)
	}

	viper.AutomaticEnv()

	// watch live changes
	viper.WatchConfig()

	upd := make(chan struct{})
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
		upd <- struct{}{}
	})

	// unmarshal viper store into a Config
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, nil, fmt.Errorf("cannot unmarshal configuration file: %w", err)
	}

	return &config, upd, nil
}
