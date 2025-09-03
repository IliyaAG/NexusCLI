package cmd

import (
    "fmt"
    "os"
    "strings"
    "nexuscli/internal/output"
    "github.com/spf13/cobra"
)

var (
    userPassword  string
    userFirstName string
    userLastName  string
    userEmail     string
    userRoles     []string
)

var userCmd = &cobra.Command{
    Use:   "user",
    Short: "Manage Nexus users",
    Long:  `Allows creating, deleting, and listing Nexus users.`,
}

var userCreateCmd = &cobra.Command{
    Use:   "create <username>",
    Short: "Create a new Nexus user",
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        username := args[0]

        if userPassword == "" || userEmail == "" {
            fmt.Println("Error: --password and --email are required.")
            _ = cmd.Help()
            os.Exit(1)
        }

        if err := nexusClient.CreateUser(username, userPassword, userFirstName, userLastName, userEmail, userRoles); err != nil {
            fmt.Printf("Error creating user '%s': %v\n", username, err)
            os.Exit(1)
        }
        fmt.Printf("User '%s' created successfully.\n", username)
    },
}

var userDeleteCmd = &cobra.Command{
    Use:   "delete <username>",
    Short: "Delete a Nexus user",
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        username := args[0]

        if err := nexusClient.DeleteUser(username); err != nil {
            fmt.Printf("Error deleting user '%s': %v\n", username, err)
            os.Exit(1)
        }
        fmt.Printf("User '%s' deleted successfully.\n", username)
    },
}

var userListCmd = &cobra.Command{
    Use:   "list",
    Short: "List all Nexus users",
    Run: func(cmd *cobra.Command, args []string) {
        users, err := nexusClient.ListUsers()
        if err != nil {
            fmt.Printf("Error listing users: %v\n", err)
            os.Exit(1)
        }

        if len(users) == 0 {
            fmt.Println("No users found.")
            return
        }

        items := []map[string]interface{}{}
        for _, user := range users {
            roles := []string{}
            if r, ok := user["roles"].([]interface{}); ok {
                for _, role := range r {
                    if s, ok := role.(string); ok {
                        roles = append(roles, s)
                    }
                }
            }
            items = append(items, map[string]interface{}{
                "USER ID":    user["userId"],
                "FIRST NAME": user["firstName"],
                "LAST NAME":  user["lastName"],
                "EMAIL":      user["emailAddress"],
                "STATUS":     user["status"],
                "ROLES":      strings.Join(roles, ", "),
            })
        }

        headers := []string{"USER ID", "FIRST NAME", "LAST NAME", "EMAIL", "STATUS", "ROLES"}
        output.Render(items, outputFormat, headers, func(u map[string]interface{}) {
            fmt.Printf("\033[36m%s\033[0m\t%s\t%s\n",
                u["USER ID"], u["FIRST NAME"], u["STATUS"])
        })
    },
}

func init() {
    rootCmd.AddCommand(userCmd)
    userCmd.AddCommand(userCreateCmd)
    userCmd.AddCommand(userDeleteCmd)
    userCmd.AddCommand(userListCmd)

    userCreateCmd.Flags().StringVarP(&userPassword, "password", "p", "", "Password for the new user (required)")
    userCreateCmd.Flags().StringVarP(&userFirstName, "first-name", "f", "", "First name of the new user")
    userCreateCmd.Flags().StringVarP(&userLastName, "last-name", "l", "", "Last name of the new user")
    userCreateCmd.Flags().StringVarP(&userEmail, "email", "e", "", "Email address of the new user (required)")
    userCreateCmd.Flags().StringSliceVarP(&userRoles, "roles", "r", []string{"nx-anonymous"}, "Comma-separated list of roles for the new user")

    _ = userCreateCmd.MarkFlagRequired("password")
    _ = userCreateCmd.MarkFlagRequired("email")
}
