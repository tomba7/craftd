package cmd

import (
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a specific k8s resource",
}

func init() {
	rootCmd.AddCommand(getCmd)
}
