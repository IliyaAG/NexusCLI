// cmd/config.go
package cmd

import (
    "fmt"
    "os"
    "nexuscli/config"
    "github.com/spf13/cobra"
)

var (
    configKey   string
    configValue string
)

// configCmd represents the config command
var configCmd = &cobra.Command{
    Use:   "config",
    Short: "Manage Nexus CLI configuration",
    Long:  `Allows viewing and setting configuration values for nexuscli.`,
}

// configViewCmd represents the view subcommand for config
var configViewCmd = &cobra.Command{
    Use:   "view",
    Short: "View the current Nexus CLI configuration",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("Current Nexus CLI Configuration:")
        fmt.Printf("  Nexus URL: %s\n", config.GlobalConfig.Nexus.URL)
        fmt.Printf("  Nexus Token: %s\n", config.GlobalConfig.Nexus.Token)
        fmt.Printf("  Nexus Username: %s/n", config.GlobalConfig.Nexus.Username)
        fmt.Printf("  Nexus Password: %s/n", maskPassword(config.GlobalConfig.Nexus.Password))
        fmt.Printf("  Insecure Skip Verify: %t\n", config.GlobalConfig.Nexus.InsecureSkipVerify)
        fmt.Printf("  Timeout Seconds: %d\n", config.GlobalConfig.Nexus.TimeoutSeconds)
    },
}

// configSetCmd represents the set subcommand for config
var configSetCmd = &cobra.Command{
    Use:   "set <key> <value>",
    Short: "Set a configuration value",
    Long: `Set a configuration value.
Supported keys: nexus.url, nexus.token, nexus.username, nexus.password, nexus.insecureSkipVerify, nexus.timeoutSeconds.`,
    Args: cobra.ExactArgs(2),
    Run: func(cmd *cobra.Command, args []string) {
        key := args[0]
        value := args[1]

        switch key {
        case "nexus.url":
            config.GlobalConfig.Nexus.URL = value
        case "nexus.token":
            config.GlobalConfig.Nexus.Token = value
        case "nexus.username":
            config.GlobalConfig.Nexus.Username = value
        case "nexus.password":
            config.GlobalConfig.Nexus.Password = value
        case "nexus.insecureSkipVerify":
            if value == "true" {
                config.GlobalConfig.Nexus.InsecureSkipVerify = true
            } else if value == "false" {
                config.GlobalConfig.Nexus.InsecureSkipVerify = false
            } else {
                fmt.Printf("Error: Invalid value for nexus.insecureSkipVerify. Use 'true' or 'false'.\n")
                os.Exit(1)
            }
        case "nexus.timeoutSeconds":
            var timeout int
            _, err := fmt.Sscanf(value, "%d", &timeout)
            if err != nil {
                fmt.Printf("Error: Invalid value for nexus.timeoutSeconds. Must be an integer.\n")
                os.Exit(1)
            }
            config.GlobalConfig.Nexus.TimeoutSeconds = timeout
        default:
            fmt.Printf("Error: Unknown config key '%s'.\n", key)
            os.Exit(1)
        }

        if err := config.SaveConfig(config.GlobalConfig); err != nil {
            fmt.Printf("Error saving configuration: %v\n", err)
            os.Exit(1)
        }
        fmt.Printf("Configuration key '%s' set to '%s' successfully.\n", key, value)
    },
}

func init() {
    rootCmd.AddCommand(configCmd)
    configCmd.AddCommand(configViewCmd)
    configCmd.AddCommand(configSetCmd)
}
func maskPassword(pw string) string {
    if pw == "" {
        return ""
    }
    return "********"
}
