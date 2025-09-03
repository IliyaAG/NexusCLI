package config

import (
    "fmt"
    "os"
    "path/filepath"

    "github.com/spf13/viper"
)

type Config struct {
    URL      string `mapstructure:"url"`
    Username string `mapstructure:"username"`
    Password string `mapstructure:"password"`
    Token    string `mapstructure:"token"`
    Timeout  int    `mapstructure:"timeout"`
}

var Global Config

// InitViper sets up viper with defaults and config file
func InitViper() error {
    home, err := os.UserHomeDir()
    if err != nil {
        return fmt.Errorf("could not find home dir: %w", err)
    }

    cfgFile := filepath.Join(home, ".nexuscli.yaml")

    viper.SetConfigFile(cfgFile)
    viper.SetConfigType("yaml")

    // Defaults
    viper.SetDefault("url", "")
    viper.SetDefault("username", "")
    viper.SetDefault("password", "")
    viper.SetDefault("token", "")
    viper.SetDefault("timeout", 30)

    // ENV support (NEXUS_URL, NEXUS_USERNAME, ...)
    viper.SetEnvPrefix("NEXUS")
    viper.AutomaticEnv()

    // Load config if exists
    if err := viper.ReadInConfig(); err != nil {
        // ignore if file does not exist
        if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
            return fmt.Errorf("failed reading config: %w", err)
        }
    }

    return viper.Unmarshal(&Global)
}

// SaveConfig writes current viper state to file
func SaveConfig() error {
    home, err := os.UserHomeDir()
    if err != nil {
        return fmt.Errorf("could not find home dir: %w", err)
    }
    cfgFile := filepath.Join(home, ".nexuscli.yaml")

    if err := viper.WriteConfigAs(cfgFile); err != nil {
        return fmt.Errorf("could not write config: %w", err)
    }
    return nil
}
