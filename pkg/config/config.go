package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// Config is configuration for request service.
type Config struct {
	RestListen string    `mapstructure:"rest_listen"`
	DB         *DBOption `mapstructure:"db"`
}

type DBOption struct {
	ConnectString string `mapstructure:"connect_string"`
}

var conf *Config

func InitConfig(configPath string) error {
	// Set viper path and read configuration.
	env := "development"
	if e := os.Getenv("ENVIRONMENT"); e != "" {
		env = e
	}

	viper.SetConfigName("config." + env)

	viper.AddConfigPath(configPath)

	replacer := strings.NewReplacer(".", "__")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv()

	// Handle errors reading the config file
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("fatal error config file-> %v", err)
	}

	// Handle errors unmarshalling the config file.
	if err := viper.Unmarshal(&conf); err != nil {
		return fmt.Errorf("fatal error to update config-> %v", err)
	}

	return nil
}

func GetConfig() (*Config, error) {
	if conf == nil {
		return nil, errors.New("config is not created")
	}
	return conf, nil
}

func GetDBOption() (*DBOption, error) {
	if conf == nil || conf.DB == nil {
		return nil, errors.New("config db is not created")
	}

	return conf.DB, nil
}
