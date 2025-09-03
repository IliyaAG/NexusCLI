package cmd

import (
    "fmt"
    "os"
    "nexuscli/config"
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
)

var (
    cfgURL     string
    cfgUsername string
    cfgPassword string
    cfgToken string
    cfgTimeout int
)

var configCmd = &cobra.Command{
    Use:   "config",
    Short: "Manage Nexus CLI configuration",
}

var configViewCmd = &cobra.Command{
    Use:   "view",
    Short: "View current configuration",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Printf("URL: %s\n", viper.GetString("url"))
        fmt.Printf("Username: %s\n", viper.GetString("username"))
        fmt.Printf("Password: %s\n", mask(viper.GetString("password")))
        fmt.Printf("Token: %s\n", mask(viper.GetString("token")))
        fmt.Printf("Timeout: %d\n", viper.GetInt("timeout"))
    },
}

var configSetCmd = &cobra.Command{
    Use:   "set",
    Short: "Set configuration values",
    Run: func(cmd *cobra.Command, args []string) {
        if cfgURL != "" {
            viper.Set("url", cfgURL)
        }
        if cfgUsername != "" {
            viper.Set("username", cfgUsername)
        }
        if cfgPassword != "" {
            viper.Set("password", cfgPassword)
        }
        if cfgToken != "" {
            viper.Set("token", cfgToken)
        }
        if cfgTimeout > 0 {
            viper.Set("timeout", cfgTimeout)
        }

        if err := config.SaveConfig(); err != nil {
            fmt.Printf("Error saving config: %v\n", err)
            os.Exit(1)
        }

        fmt.Println("Configuration updated successfully.")
    },
}

func init() {
    rootCmd.AddCommand(configCmd)
    configCmd.AddCommand(configViewCmd)
    configCmd.AddCommand(configSetCmd)

    configSetCmd.Flags().StringVar(&cfgURL, "url", "", "Nexus server URL")
    configSetCmd.Flags().StringVar(&cfgUsername, "username", "", "Nexus username")
    configSetCmd.Flags().StringVar(&cfgPassword, "password", "", "Nexus password")
    configSetCmd.Flags().StringVar(&cfgToken, "token", "", "Nexus API token")
    configSetCmd.Flags().IntVar(&cfgTimeout, "timeout", 0, "Request timeout in seconds")
}

func mask(s string) string {
    if s == "" {
        return ""
    }
    if len(s) <= 4 {
        return "****"
    }
    return s[:2] + "****" + s[len(s)-2:]
}
