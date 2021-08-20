package pods

import (
	v1 "k8s.io/api/core/v1"
	"time"
)

type Filter int

const (
	FilterAll Filter = iota
	FilterTerminating
	FilterFailed
	FilterCompleted
)

const (
	FilterAllString         = "all"
	FilterTerminatingString = "terminating"
	FilterFailedString      = "failed"
	FilterCompletedString   = "completed"
)

func (f Filter) String() string {
	switch f {
	case FilterTerminating:
		return FilterTerminatingString
	case FilterFailed:
		return FilterFailedString
	case FilterCompleted:
		return FilterCompletedString
	default:
		return FilterAllString
	}
}

func ToFilter(filterString string) Filter {
	switch filterString {
	case FilterTerminatingString:
		return FilterTerminating
	case FilterFailedString:
		return FilterFailed
	case FilterCompletedString:
		return FilterCompleted
	default:
		return FilterAll
	}
}

type specifier interface {
	satisfies(pod *v1.Pod) bool
}

type StatusFilters []Filter

var defaultStatusFilters = map[Filter]specifier{
	FilterTerminating: &filterTerminatingPods{},
	FilterCompleted:   &filterCompletedPods{},
	FilterFailed:      &filterFailedPods{},
}

// filterTerminatingPods filters pods in terminating state
type filterTerminatingPods struct{}

func (f *filterTerminatingPods) satisfies(pod *v1.Pod) bool {
	return status(pod) == PodStatusTerminating
}

// filterCompletedPods filters pods in completed state
type filterCompletedPods struct{}

func (f *filterCompletedPods) satisfies(pod *v1.Pod) bool {
	return status(pod) == PodStatusCompleted
}

// filterFailedPods filters failed pods
type filterFailedPods struct{}

func (f *filterFailedPods) satisfies(pod *v1.Pod) bool {
	podStatus := status(pod)
	return podStatus == PodStatusFailed || podStatus == PodStatusCrashLoopBackOff ||
		podStatus == PodStatusImagePullBackOff
}

type statusSpecifier struct {
	filters StatusFilters
}

func (s *statusSpecifier) satisfies(pod *v1.Pod) bool {
	for _, filter := range s.filters {
		if defaultStatusFilters[filter].satisfies(pod) {
			return true
		}
	}
	return false
}

type timeSpecifier struct {
	duration time.Duration
}

func (s *timeSpecifier) satisfies(pod *v1.Pod) bool {
	return time.Since(pod.Status.StartTime.Time) >= s.duration
}

type statusAndTimeSpecifier struct {
	status, time specifier
}

func (s *statusAndTimeSpecifier) satisfies(pod *v1.Pod) bool {
	return s.status.satisfies(pod) && s.time.satisfies(pod)
}