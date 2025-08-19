// config/config.go
package config

import (
    "fmt"
    "os"
    "path/filepath"
    "time"
    "github.com/spf13/viper"
)

// Config holds all application configuration
type Config struct {
    Nexus struct {
        URL                string `mapstructure:"url"`
        Token              string `mapstructure:"token"`
        Username           string `mapstructure:"username"`
        Password           string `mapstructure:"password"`
        InsecureSkipVerify bool   `mapstructure:"insecureSkipVerify"` 
        TimeoutSeconds     int    `mapstructure:"timeoutSeconds"`
    } `mapstructure:"nexus"`
}

// GlobalConfig holds the loaded configuration
var GlobalConfig Config

// InitConfig initializes the configuration from file or environment variables
func InitConfig() error {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return fmt.Errorf("failed to get user home directory: %w", err)
    }

    configDir := filepath.Join(homeDir, ".nexuscli")
    configFilePath := filepath.Join(configDir, "config.yaml")

    viper.AddConfigPath(configDir)
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")

 // Set default values
    viper.SetDefault("nexus.url", "http://localhost:8081")
    viper.SetDefault("nexus.token", "")
    viper.SetDefault("nexus.username", "")
    viper.SetDefault("nexus.password", "")
    viper.SetDefault("nexus.insecureSkipVerify", false)
    viper.SetDefault("nexus.timeoutSeconds", 30)

 // Read environment variables (e.g., NEXUS_URL, NEXUS_TOKEN)
    viper.SetEnvPrefix("NEXUS")
    viper.AutomaticEnv()

 // Create config directory if it doesn't exist
    if _, err := os.Stat(configDir); os.IsNotExist(err) {
        if err := os.MkdirAll(configDir, 0700); err != nil {
            return fmt.Errorf("failed to create config directory %s: %w", configDir, err)
        }
    }

 // Try to read the config file
    if err := viper.ReadInConfig(); err != nil {
        if _, ok := err.(viper.ConfigFileNotFoundError); ok {
            fmt.Printf("Config file not found at %s. Creating a default one.\n", configFilePath)
            if err := SaveConfig(GlobalConfig); err != nil {
                return fmt.Errorf("failed to save default config: %w", err)
            }
         } else {
             return fmt.Errorf("failed to read config file: %w", err)
           }
    }

 // Unmarshal the config into our struct
    if err := viper.Unmarshal(&GlobalConfig); err != nil {
        return fmt.Errorf("failed to unmarshal config: %w", err)
    }
    return nil
}

// SaveConfig saves the current configuration to the config file
func SaveConfig(cfg Config) error {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return fmt.Errorf("failed to get user home directory: %w", err)
    }

    configDir := filepath.Join(homeDir, ".nexuscli")
    configFilePath := filepath.Join(configDir, "config.yaml")

 // Ensure the directory exists
    if err := os.MkdirAll(configDir, 0700); err != nil {
        return fmt.Errorf("failed to create config directory %s: %w", configDir, err)
    }

 // Write the config to file
    viper.Set("nexus", cfg.Nexus)
    if err := viper.WriteConfigAs(configFilePath); err != nil {
        return fmt.Errorf("failed to write config file %s: %w", configFilePath, err)
    }
    return nil
}

// GetTimeout returns the timeout duration
func (c *Config) GetTimeout() time.Duration {
    return time.Duration(c.Nexus.TimeoutSeconds) * time.Second
}
