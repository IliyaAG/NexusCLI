// cmd/repo.go
package cmd

import (
    "fmt"
    "os"
    "text/tabwriter"
    "github.com/spf13/cobra"
)

var (
    repoType string // Flag for repository type (e.g., docker-proxy, docker-hosted)
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

        w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.Debug)
        fmt.Fprintln(w, "NAME\tFORMAT\tTYPE\tURL\tONLINE")
        for _, repo := range repos {
            name := repo["name"].(string)
            format := repo["format"].(string)
            repoType := repo["type"].(string)
            url := repo["url"].(string)
            online := repo["online"].(bool)
            fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%t\n", name, format, repoType, url, online)
        }
        w.Flush()
    },
}

func init() {
    rootCmd.AddCommand(repoCmd)
    repoCmd.AddCommand(repoCreateCmd)
    repoCmd.AddCommand(repoDeleteCmd)
    repoCmd.AddCommand(repoListCmd)
}
