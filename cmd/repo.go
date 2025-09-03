package cmd

import (
    "fmt"
    "os"
    "nexuscli/internal/output"
    "github.com/spf13/cobra"
)

var repoCmd = &cobra.Command{
    Use:   "repo",
    Short: "Manage Nexus repositories",
    Long:  `Allows creating, deleting, and listing Nexus repositories.`,
}

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

        items := []map[string]interface{}{}
        for _, repo := range repos {
            items = append(items, map[string]interface{}{
                "NAME":   repo["name"],
                "FORMAT": repo["format"],
                "TYPE":   repo["type"],
                "URL":    repo["url"],
                "ONLINE": repo["online"],
            })
        }

        headers := []string{"NAME", "FORMAT", "TYPE", "URL", "ONLINE"}
        output.Render(items, outputFormat, headers, func(r map[string]interface{}) {
            fmt.Printf("\033[32m%s\033[0m\t%s\t%s\t%s\t%v\n",
                r["NAME"], r["FORMAT"], r["TYPE"], r["URL"], r["ONLINE"])
        })
    },
}

func init() {
    rootCmd.AddCommand(repoCmd)
    repoCmd.AddCommand(repoCreateCmd)
    repoCmd.AddCommand(repoDeleteCmd)
    repoCmd.AddCommand(repoListCmd)
}
