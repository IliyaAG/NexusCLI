// cmd/repo.go
package cmd

import (
    "encoding/json"
    "fmt"
    "os"
    "text/tabwriter"
    "gopkg.in/yaml.v3"
    "github.com/spf13/cobra"
)

var (
    repoType string
)

// repoCmd represents the repo command
var repoCmd = &cobra.Command{
    Use:   "repo",
    Short: "Manage Nexus repositories",
    Long:  `Allows creating, deleting, and listing Nexus repositories.`,
}

// repoCreateCmd represents the create subcommand for repo
var repoCreateCmd = &cobra.Command{
    Use:   "create <repo_type> <repo_name>",
    Short: "Create a new Nexus repository",
    Args:  cobra.ExactArgs(2),
    Run: func(cmd *cobra.Command, args []string) {
        repoType := args[0]
        repoName := args[1]

        if err := nexusClient.CreateRepository(repoType, repoName); err != nil {
            fmt.Printf("Error creating repository '%s': %v\n", repoName, err)
            os.Exit(1)
        }
        fmt.Printf("Repository '%s' of type '%s' created successfully.\n", repoName, repoType)
    },
}

// repoDeleteCmd represents the delete subcommand for repo
var repoDeleteCmd = &cobra.Command{
    Use:   "delete <repo_name>",
    Short: "Delete a Nexus repository",
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        repoName := args[0]

        if err := nexusClient.DeleteRepository(repoName); err != nil {
            fmt.Printf("Error deleting repository '%s': %v\n", repoName, err)
            os.Exit(1)
        }
        fmt.Printf("Repository '%s' deleted successfully.\n", repoName)
    },
}

// repoListCmd represents the list subcommand for repo
var repoListCmd = &cobra.Command{
    Use:   "list",
    Short: "List all Nexus repositories",
    Run: func(cmd *cobra.Command, args []string) {
        repos, err := nexusClient.ListRepositories()
        if err != nil {
            fmt.Printf("Error listing repositories: %v\n", err)
            os.Exit(1)
        }

        if len(repos) == 0 {
            fmt.Println("No repositories found.")
            return
        }

        switch outputFormat {
        case "json":
            data, _ := json.MarshalIndent(repos, "", "  ")
            fmt.Println(string(data))
        case "yaml", "yml":
            data, _ := yaml.Marshal(repos)
            fmt.Println(string(data))
        case "color":
            for _, repo := range repos {
                fmt.Printf("\033[32m%s\033[0m\t%s\t%s\t%s\t%t\n",
                    repo["name"], repo["format"], repo["type"], repo["url"], repo["online"])
            }
        default:
            w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.Debug)
            fmt.Fprintln(w, "NAME\tFORMAT\tTYPE\tURL\tONLINE")
            for _, repo := range repos {
                name, _ := repo["name"].(string)
                format, _ := repo["format"].(string)
                repoType, _ := repo["type"].(string)
                url, _ := repo["url"].(string)

                online := false
                if v, ok := repo["online"].(bool); ok {
                    online = v
                }
                fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%t\n", name, format, repoType, url, online)
            }
            w.Flush()
        }
    },
}

func init() {
    rootCmd.AddCommand(repoCmd)
    repoCmd.AddCommand(repoCreateCmd)
    repoCmd.AddCommand(repoDeleteCmd)
    repoCmd.AddCommand(repoListCmd)
}
