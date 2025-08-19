package cmd

import (
    "fmt"
    "os"
    "text/tabwriter"
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
    Short: "List all available commands in nexuscli",
    Run: func(cmd *cobra.Command, args []string) {
		// Create a tabwriter for pretty output
        w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
        fmt.Fprintln(w, "COMMAND\tDESCRIPTION")

		// Get all root subcommands
        for _, c := range rootCmd.Commands() {
            fmt.Fprintf(w, "%s\t%s\n", c.Use, c.Short)

			// Print subcommands of each command too
            for _, sub := range c.Commands() {
                fmt.Fprintf(w, "  %s %s\t%s\n", c.Use, sub.Use, sub.Short)
            }
        }

        w.Flush()
    },
}

func init() {
    rootCmd.AddCommand(commandCmd)
    commandCmd.AddCommand(commandListCmd)
}
