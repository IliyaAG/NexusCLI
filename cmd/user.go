// cmd/user.go
package cmd

import (
    "fmt"
    "os"
    "strings"
    "text/tabwriter"
    "github.com/spf13/cobra"
)

var (
    userPassword string
    userFirstName string
    userLastName  string
    userEmail     string
    userRoles     []string
)

// userCmd represents the user command
var userCmd = &cobra.Command{
    Use:   "user",
    Short: "Manage Nexus users",
    Long:  `Allows creating, deleting, and listing Nexus users.`,
}

// userCreateCmd represents the create subcommand for user
var userCreateCmd = &cobra.Command{
    Use:   "create <username>",
    Short: "Create a new Nexus user",
    Args:  cobra.ExactArgs(1), // Requires username
    Run: func(cmd *cobra.Command, args []string) {
        username := args[0]

        if userPassword == "" {
            fmt.Println("Error: --password is required.")
            _ = cmd.Help()
            os.Exit(1)
        }
        if userEmail == "" {
            fmt.Println("Error: --email is required.")
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

// userDeleteCmd represents the delete subcommand for user
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

// userListCmd represents the list subcommand for user
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

        w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.Debug)
        fmt.Fprintln(w, "USER ID\tFIRST NAME\tLAST NAME\tEMAIL\tSTATUS\tROLES")
        for _, user := range users {
            userID := user["userId"].(string)
            firstName := user["firstName"].(string)
            lastName := user["lastName"].(string)
            email := user["emailAddress"].(string)
            status := user["status"].(string)
            roles := []string{}
        if r, ok := user["roles"].([]interface{}); ok {
            for _, role := range r {
                if s, ok := role.(string); ok {
                    roles = append(roles, s)
                }
            }
        }
        fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n", userID, firstName, lastName, email, status, strings.Join(roles, ", "))
        }
        w.Flush()
    },
}

func init() {
    rootCmd.AddCommand(userCmd)
    userCmd.AddCommand(userCreateCmd)
    userCmd.AddCommand(userDeleteCmd)
    userCmd.AddCommand(userListCmd)

 // Flags for user create command
    userCreateCmd.Flags().StringVarP(&userPassword, "password", "p", "", "Password for the new user (required)")
    userCreateCmd.Flags().StringVarP(&userFirstName, "first-name", "f", "", "First name of the new user")
    userCreateCmd.Flags().StringVarP(&userLastName, "last-name", "l", "", "Last name of the new user")
    userCreateCmd.Flags().StringVarP(&userEmail, "email", "e", "", "Email address of the new user (required)")
    userCreateCmd.Flags().StringSliceVarP(&userRoles, "roles", "r", []string{"nx-anonymous"}, "Comma-separated list of roles for the new user")

 // Mark required flags
    _ = userCreateCmd.MarkFlagRequired("password")
    _ = userCreateCmd.MarkFlagRequired("email")
}
