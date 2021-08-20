package filters

import (
	"github.com/spf13/viper"
	"github.com/tomba7/craftd/pkg/pods"
)

func ParseFilters() pods.StatusFilters {
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

func ParseDuration() *int {
	var minutes *int
	if duration := viper.GetInt("minutes"); duration != 0 {
		minutes = &duration
	}
	return minutes
}