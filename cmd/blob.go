package cmd

import (
    "fmt"
    "nexuscli/internal/output"

    "github.com/spf13/cobra"
)

var (
    blobPath string
)

var blobCmd = &cobra.Command{
    Use:   "blob",
    Short: "Manage Nexus Blob Stores",
}

var blobListCmd = &cobra.Command{
    Use:   "list",
    Short: "List all blob stores",
    RunE: func(cmd *cobra.Command, args []string) error {
        blobs, err := nexusClient.ListBlobStores()
        if err != nil {
            return err
        }
        return output.Print(blobs, outputFormat)
    },
}

var blobCreateCmd = &cobra.Command{
    Use:   "create [name]",
    Short: "Create a new blob store",
    Args:  cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        name := args[0]
        if blobPath == "" {
            return fmt.Errorf("you must provide --path for the blob store")
        }
        if err := nexusClient.CreateBlobStore(name, blobPath); err != nil {
            return err
        }
        fmt.Printf("Blob store '%s' created successfully\n", name)
        return nil
    },
}

var blobDeleteCmd = &cobra.Command{
    Use:   "delete [name]",
    Short: "Delete a blob store",
    Args:  cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        name := args[0]
        if err := nexusClient.DeleteBlobStore(name); err != nil {
            return err
        }
        fmt.Printf("Blob store '%s' deleted successfully\n", name)
        return nil
    },
}

func init() {
    rootCmd.AddCommand(blobCmd)

    blobCmd.AddCommand(blobListCmd)
    blobCmd.AddCommand(blobCreateCmd)
    blobCmd.AddCommand(blobDeleteCmd)

    blobCreateCmd.Flags().StringVar(&blobPath, "path", "", "Path for the blob store (required)")
}
