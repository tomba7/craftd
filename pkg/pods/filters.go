package pods

import (
	"fmt"
	v1 "k8s.io/api/core/v1"
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

type Specifier interface {
	Satisfies(pod *v1.Pod) bool
}

type StatusFilters []Filter

var DefaultStatusFilters = map[Filter]Specifier{
	FilterTerminating: &FilterTerminatingPods{},
	FilterCompleted:   &FilterCompletedPods{},
	FilterFailed:      &FilterFailedPods{},
}

// FilterTerminatingPods filters pods in terminating state
type FilterTerminatingPods struct{}

func (f *FilterTerminatingPods) Satisfies(pod *v1.Pod) bool {
	return status(pod) == PodStatusTerminating
}

// FilterCompletedPods filters pods in completed state
type FilterCompletedPods struct{}

func (f *FilterCompletedPods) Satisfies(pod *v1.Pod) bool {
	return status(pod) == PodStatusCompleted
}

// FilterFailedPods filters failed pods
type FilterFailedPods struct{}

func (f *FilterFailedPods) Satisfies(pod *v1.Pod) bool {
	podStatus := status(pod)
	return podStatus == PodStatusFailed || podStatus == PodStatusCrashLoopBackOff ||
		podStatus == PodStatusImagePullBackOff
}

type podFilter struct{}

type statusSpecifier struct {
	filters StatusFilters
}

func (s *statusSpecifier) Satisfies(pod *v1.Pod) bool {
	for _, filter := range s.filters {
		if DefaultStatusFilters[filter].Satisfies(pod) {
			return true
		}
	}
	return false
}

func (f *podFilter) filter(podList *v1.PodList, spec Specifier) []*Pod {
	set := make(map[string]*Pod)
	key := func(pod *v1.Pod) string {
		return fmt.Sprintf("%s,%s", pod.Namespace, pod.Name)
	}
	for _, pod := range podList.Items {
		if spec.Satisfies(&pod) {
			if _, exist := set[key(&pod)]; !exist {
				set[key(&pod)] = &Pod{
					Name:      pod.Name,
					Status:    status(&pod),
					Namespace: pod.Namespace,
				}
			}
		}
	}
	var result []*Pod
	for _, pod := range set {
		result = append(result, pod)
	}
	return result
}
