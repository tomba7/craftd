package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// deletePodsCmd represents the pods command
var deletePodsCmd = &cobra.Command{
	Use:   "pods",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("delete pods called")
	},
}

func init() {
	deleteCmd.AddCommand(deletePodsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deletePodsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deletePodsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
