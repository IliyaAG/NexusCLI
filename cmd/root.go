package cmd

import (
    "fmt"
    "os"
    "nexuscli/config"
    "nexuscli/nexus"
    "github.com/spf13/cobra"
)

var (
    nexusClient  *nexus.Client
    outputFormat string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
    Use:   "nexuscli",
    Short: "A command-line interface for Nexus Repository Manager",
    Long: `nexuscli is a powerful command-line tool to interact with Nexus Repository Manager.
It allows you to manage repositories, users, and more directly from your terminal.`,
    PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
        // Load config before any command runs
        if err := config.InitConfig(); err != nil {
            return fmt.Errorf("failed to initialize configuration: %w", err)
        }
        // Initialize Nexus client
        nexusClient = nexus.NewClient(config.GlobalConfig)
        return nil
    },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
    if err := rootCmd.Execute(); err != nil {
        os.Exit(1)
    }
}

func init() {
    // Add global output flag
    rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "table", "Output format: table|json|yaml|color")

    // Example: if you want to add global config file path later
    // rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.nexuscli/config.yaml)")
}
