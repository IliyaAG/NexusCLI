package cmd

import (
    "fmt"
    "github.com/spf13/cobra"
)

// commandCmd represents the "command" parent command
var commandCmd = &cobra.Command{
    Use:   "command",
    Short: "Manage and inspect available CLI commands",
    Long:  `The 'command' group helps you to inspect available commands of nexuscli.`,
}

// commandListCmd represents the "list" subcommand
var commandListCmd = &cobra.Command{
    Use:   "list",
    Short: "List all available commands in nexuscli (tree view)",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("Available Commands:")
        printCommands(rootCmd, 0)
    },
}

// Recursive function to print commands with indent
func printCommands(c *cobra.Command, indent int) {
    prefix := ""
    for i := 0; i < indent; i++ {
        prefix += "  " // two spaces indent
    }

	// Skip printing the root itself, just its children
    if c != rootCmd {
        fmt.Printf("%s- %s : %s\n", prefix, c.Use, c.Short)
    }

    for _, child := range c.Commands() {
        printCommands(child, indent+1)
    }
}

func init() {
    rootCmd.AddCommand(commandCmd)
    commandCmd.AddCommand(commandListCmd)
}
