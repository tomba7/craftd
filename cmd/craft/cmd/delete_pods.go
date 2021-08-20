package cmd

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/tomba7/craftd/cmd/craft/filters"
	"github.com/tomba7/craftd/pkg/pods"

	"github.com/spf13/cobra"
)

// deletePodsCmd represents the pods command
var deletePodsCmd = &cobra.Command{
	Use:   "pods",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		filterList := filters.ParseFilters()
		minutes := filters.ParseDuration()
		for _, pod := range pods.Service.Delete(filterList, minutes) {
			fmt.Println(pod.Name, "deleted successfully")
		}
	},
}

func init() {
	deleteCmd.AddCommand(deletePodsCmd)

	deletePodsCmd.Flags().BoolP(pods.FilterCompletedString, "c", false, "Completed Pods")
	deletePodsCmd.Flags().BoolP(pods.FilterFailedString, "f", false, "Failed Pods")
	deletePodsCmd.Flags().BoolP(pods.FilterTerminatingString, "t", false, "Terminating Pods")
	deletePodsCmd.Flags().IntP("minutes", "m", 0, "Duration a Pod has been running (in minutes)")

	viper.BindPFlag(pods.FilterCompletedString, deletePodsCmd.Flag(pods.FilterCompletedString))
	viper.BindPFlag(pods.FilterFailedString, deletePodsCmd.Flag(pods.FilterFailedString))
	viper.BindPFlag(pods.FilterTerminatingString, deletePodsCmd.Flag(pods.FilterTerminatingString))
	viper.BindPFlag("minutes", deletePodsCmd.Flag("minutes"))
}
