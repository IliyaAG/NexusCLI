package cmd

import (
    "encoding/json"
    "fmt"
    "os"
    "text/tabwriter"
    "gopkg.in/yaml.v3"
    "github.com/spf13/cobra"
)

var commandCmd = &cobra.Command{
    Use:   "command",
    Short: "Manage CLI commands",
    Long:  "List available CLI commands and their descriptions.",
}

var commandListCmd = &cobra.Command{
    Use:   "list",
    Short: "List all available commands",
    Run: func(cmd *cobra.Command, args []string) {
        commands := []map[string]string{}
        for _, c := range rootCmd.Commands() {
            commands = append(commands, map[string]string{
                "name":        c.Use,
                "description": c.Short,
            })
        }

        switch outputFormat {
        case "json":
            data, _ := json.MarshalIndent(commands, "", "  ")
            fmt.Println(string(data))
        case "yaml", "yml":
            data, _ := yaml.Marshal(commands)
            fmt.Println(string(data))
        case "color":
            for _, c := range commands {
                fmt.Printf("\033[34m%s\033[0m - %s\n", c["name"], c["description"])
            }
        default:
            w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
            fmt.Fprintln(w, "COMMAND\tDESCRIPTION")
            for _, c := range commands {
            fmt.Fprintf(w, "%s\t%s\n", c["name"], c["description"])
            }
            w.Flush()
        }
    },
}

func init() {
    rootCmd.AddCommand(commandCmd)
    commandCmd.AddCommand(commandListCmd)
}
