package cmd

import (
    "fmt"
    "github.com/spf13/cobra"
)

var commandCmd = &cobra.Command{
    Use:   "command",
    Short: "List all available commands",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("Available commands:")
        for _, c := range rootCmd.Commands() {
            fmt.Printf("  %s\t%s\n", c.Name(), c.Short)
        }
    },
}

func init() {
    rootCmd.AddCommand(commandCmd)
}
