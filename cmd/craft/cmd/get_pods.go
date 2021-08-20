package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	filters "github.com/tomba7/craftd/cmd/craft/filters"
	"github.com/tomba7/craftd/pkg/pods"
)

// getPodsCmd represents the pods command
var getPodsCmd = &cobra.Command{
	Use:   "pods",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		filterList := filters.ParseFilters()
		minutes := filters.ParseDuration()
		fmt.Printf("%-20s %-20s %-10s %-20s\n", "NAMESPACE", "STATUS", "DURATION", "NAME")
		for _, pod := range pods.Service.Get(filterList, minutes) {
			fmt.Printf("%-20s %-20s %-10s %-20s\n", pod.Namespace, pod.Status, pod.Duration, pod.Name)
		}
	},
}

func init() {
	getCmd.AddCommand(getPodsCmd)

	getPodsCmd.Flags().BoolP(pods.FilterCompletedString, "c", false, "Completed Pods")
	getPodsCmd.Flags().BoolP(pods.FilterFailedString, "f", false, "Failed Pods")
	getPodsCmd.Flags().BoolP(pods.FilterTerminatingString, "t", false, "Terminating Pods")
	getPodsCmd.Flags().IntP("minutes", "m", 0, "Duration a Pod has been running (in minutes)")

	viper.BindPFlag(pods.FilterCompletedString, getPodsCmd.Flag(pods.FilterCompletedString))
	viper.BindPFlag(pods.FilterFailedString, getPodsCmd.Flag(pods.FilterFailedString))
	viper.BindPFlag(pods.FilterTerminatingString, getPodsCmd.Flag(pods.FilterTerminatingString))
	viper.BindPFlag("minutes", getPodsCmd.Flag("minutes"))
}
