package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tomba7/craftd/pkg/pods"
)

// getPodsCmd represents the pods command
var getPodsCmd = &cobra.Command{
	Use:   "pods",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		filters := parseFilters()
		fmt.Printf("%-20s %-20s %-20s\n", "NAMESPACE", "STATUS", "NAME")
		for _, pod := range pods.Service.Get(filters) {
			fmt.Printf("%-20s %-20s %-20s\n", pod.Namespace, pod.Status, pod.Name)
		}
	},
}

func parseFilters() pods.Filters {
	var filters []pods.Filter
	filterFlags := []string{
		pods.FilterCompletedString,
		pods.FilterFailedString,
		pods.FilterTerminatingString,
	}
	for _, filter := range filterFlags {
		if viper.GetBool(filter) {
			filters = append(filters, pods.ToFilter(filter))
		}
	}
	return filters
}

func init() {
	getCmd.AddCommand(getPodsCmd)

	getPodsCmd.Flags().BoolP(pods.FilterCompletedString, "c", false, "Completed Pods")
	getPodsCmd.Flags().BoolP(pods.FilterFailedString, "f", false, "Failed Pods")
	getPodsCmd.Flags().BoolP(pods.FilterTerminatingString, "t", false, "Terminating Pods")

	viper.BindPFlag(pods.FilterCompletedString, getPodsCmd.Flag(pods.FilterCompletedString))
	viper.BindPFlag(pods.FilterFailedString, getPodsCmd.Flag(pods.FilterFailedString))
	viper.BindPFlag(pods.FilterTerminatingString, getPodsCmd.Flag(pods.FilterTerminatingString))
}
