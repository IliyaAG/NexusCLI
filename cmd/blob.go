package cmd

import (
    "fmt"
    "os"
    "nexuscli/internal/output"
    "github.com/spf13/cobra"
)

var (
    blobPath string
)

var blobCmd = &cobra.Command{
    Use:   "blob",
    Short: "Manage Nexus blob stores",
    Long:  `Create, list, and delete Nexus blob stores.`,
}

var blobListCmd = &cobra.Command{
    Use:   "list",
    Short: "List all blob stores",
    Run: func(cmd *cobra.Command, args []string) {
        blobs, err := nexusClient.ListBlobStores()
        if err != nil {
            fmt.Printf("Error listing blob stores: %v\n", err)
            os.Exit(1)
        }

        if len(blobs) == 0 {
            fmt.Println("No blob stores found.")
            return
        }

        items := []map[string]interface{}{}
        for _, b := range blobs {
            items = append(items, map[string]interface{}{
                "NAME": b["name"],
                "TYPE": b["type"],
                "PATH": b["path"],
            })
        }

        headers := []string{"NAME", "TYPE", "PATH"}
        output.Render(items, outputFormat, headers, func(r map[string]interface{}) {
            fmt.Printf("\033[33m%s\033[0m\t%s\t%s\n",
                r["NAME"], r["TYPE"], r["PATH"])
        })
    },
}

var blobCreateCmd = &cobra.Command{
    Use:   "create <name>",
    Short: "Create a new file blob store",
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        name := args[0]
        if blobPath == "" {
            fmt.Println("Error: --path is required for creating a file blob store.")
            _ = cmd.Help()
            os.Exit(1)
        }

        if err := nexusClient.CreateBlobStore(name, blobPath); err != nil {
            fmt.Printf("Error creating blob store '%s': %v\n", name, err)
            os.Exit(1)
        }
        fmt.Printf("Blob store '%s' created successfully (path: %s).\n", name, blobPath)
    },
}

var blobDeleteCmd = &cobra.Command{
    Use:   "delete <name>",
    Short: "Delete a blob store",
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        name := args[0]

        if err := nexusClient.DeleteBlobStore(name); err != nil {
            fmt.Printf("Error deleting blob store '%s': %v\n", name, err)
            os.Exit(1)
        }
        fmt.Printf("Blob store '%s' deleted successfully.\n", name)
    },
}

func init() {
    rootCmd.AddCommand(blobCmd)

    blobCmd.AddCommand(blobListCmd)
    blobCmd.AddCommand(blobCreateCmd)
    blobCmd.AddCommand(blobDeleteCmd)

    blobCreateCmd.Flags().StringVar(&blobPath, "path", "", "Filesystem path for the file blob store (required)")
}
