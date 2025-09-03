package cmd

import (
    "fmt"
    "os"
    "nexuscli/config"
    "nexuscli/internal/client"
    "github.com/spf13/cobra"
)

var (
    outputFormat string
    nexusClient  *client.NexusClient
)

var rootCmd = &cobra.Command{
    Use:   "nexuscli",
    Short: "CLI tool for Sonatype Nexus",
    PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
        if err := config.InitViper(); err != nil {
            return err
        }

        nexusClient = client.NewNexusClient(
            config.Global.URL,
            config.Global.Username,
            config.Global.Password,
            config.Global.Token,
            config.Global.Timeout,
        )
        return nil
    },
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}

func init() {
    rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "table",
        "Output format: table, json, yaml, color")
}
