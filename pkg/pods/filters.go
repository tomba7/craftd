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

type FilterSpecifier interface {
	IsSatisfied(pod *v1.Pod) bool
}

type Filters []Filter

func (f Filters) filter(podList *v1.PodList) []*Pod {
	set := make(map[string]*Pod)
	key := func(pod *v1.Pod) string {
		return fmt.Sprintf("%s,%s", pod.Namespace, pod.Name)
	}
	for _, pod := range podList.Items {
		for _, filter := range f {
			if DefaultFilters[filter].IsSatisfied(&pod) {
				if _, exist := set[key(&pod)]; !exist {
					set[key(&pod)] = &Pod{
						Name:      pod.Name,
						Status:    status(&pod),
						Namespace: pod.Namespace,
					}
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

var DefaultFilters = map[Filter]FilterSpecifier{
	FilterTerminating: &FilterTerminatingPods{},
	FilterCompleted:   &FilterCompletedPods{},
	FilterFailed:      &FilterFailedPods{},
}

type FilterTerminatingPods struct{}

func (f *FilterTerminatingPods) IsSatisfied(pod *v1.Pod) bool {
	return status(pod) == PodStatusTerminating
}

type FilterCompletedPods struct{}

func (f *FilterCompletedPods) IsSatisfied(pod *v1.Pod) bool {
	return status(pod) == PodStatusCompleted
}

type FilterFailedPods struct{}

func (f *FilterFailedPods) IsSatisfied(pod *v1.Pod) bool {
	podStatus := status(pod)
	return podStatus == PodStatusFailed ||
		podStatus == PodStatusCrashLoopBackOff ||
		podStatus == PodStatusImagePullBackOff
}
