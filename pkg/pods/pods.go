package pods

import (
	"context"
	"fmt"
	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"strings"
	"time"
)

type Status int

const (
	PodStatusUnknown Status = iota
	PodStatusFailed
	PodStatusContainerCreating
	PodStatusRunning
	PodStatusCompleted
	PodStatusTerminating
	PodStatusCrashLoopBackOff
	PodStatusImagePullBackOff
)

func (s Status) String() string {
	switch s {
	case PodStatusFailed:
		return "Failed"
	case PodStatusContainerCreating:
		return "ContainerCreating"
	case PodStatusRunning:
		return "Running"
	case PodStatusCompleted:
		return "Completed"
	case PodStatusTerminating:
		return "Terminating"
	case PodStatusCrashLoopBackOff:
		return "CrashLoopBackOff"
	case PodStatusImagePullBackOff:
		return "ImagePullBackoff"
	default:
		return "Unknown"
	}
}

type Pod struct {
	Name      string
	Namespace string
	Status    Status
	Duration  string
}

type podLister struct {
	clientSet kubernetes.Interface
}

func newPodLister(clientSet kubernetes.Interface) *podLister {
	return &podLister{clientSet: clientSet}
}

func (p *podLister) get(namespace string, filters StatusFilters, minutes *int) []*Pod {
	podList, err := p.clientSet.
		CoreV1().
		Pods(namespace).
		List(context.TODO(), metaV1.ListOptions{LabelSelector: "app=intuit"})
	if err != nil {
		panic(err.Error())
	}

	var timeSpec *timeSpecifier
	if minutes != nil {
		timeSpec = &timeSpecifier{
			duration: time.Minute * time.Duration(*minutes),
		}
	}

	// Do we need to filter data?
	if len(filters) != 0 {
		statusSpec := &statusSpecifier{filters: filters}
		if minutes != nil {
			statusAndTimeSpec := &statusAndTimeSpecifier{
				status: statusSpec, time: timeSpec,
			}
			// Filter both status and duration
			return filter(podList, statusAndTimeSpec)
		} else {
			// Filter status only
			return filter(podList, statusSpec)
		}
	} else {
		if minutes != nil {
			// Filter by duration only
			return filter(podList, timeSpec)
		}
		var result []*Pod
		for _, pod := range podList.Items {
			result = append(result, &Pod{
				Name:      pod.Name,
				Status:    status(&pod),
				Namespace: pod.Namespace,
				Duration:  timeSince(time.Since(pod.Status.StartTime.Time)),
			})
		}
		return result
	}
}

func filter(podList *v1.PodList, spec specifier) []*Pod {
	set := make(map[string]*Pod)
	key := func(pod *v1.Pod) string {
		return fmt.Sprintf("%s,%s", pod.Namespace, pod.Name)
	}
	for _, pod := range podList.Items {
		if spec.satisfies(&pod) {
			if _, exist := set[key(&pod)]; !exist {
				set[key(&pod)] = &Pod{
					Name:      pod.Name,
					Status:    status(&pod),
					Namespace: pod.Namespace,
					Duration: timeSince(time.Since(pod.Status.StartTime.Time)),
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

func status(pod *v1.Pod) Status {
	// Takes precedence over all states
	if pod.ObjectMeta.DeletionTimestamp != nil {
		return PodStatusTerminating
	} else if pod.Status.Phase == v1.PodSucceeded {
		return PodStatusCompleted
	} else if pod.Status.Phase == v1.PodFailed || pod.Status.Phase == v1.PodUnknown {
		return PodStatusFailed
	} else {
		for _, container := range pod.Status.ContainerStatuses {
			if container.State.Waiting != nil {
				reason := container.State.Waiting.Reason
				if strings.Contains(reason, "CrashLoopBackOff") {
					return PodStatusCrashLoopBackOff
				} else if strings.Contains(reason, "ImagePullBackOff") {
					return PodStatusImagePullBackOff
				}
			}
		}
	}
	return PodStatusRunning
}

func timeSince(duration time.Duration) string {
	duration = duration.Round(time.Minute)
	hours := duration / time.Hour
	duration -= hours * time.Hour
	minutes := duration / time.Minute
	return fmt.Sprintf("%02d:%02d", hours, minutes)
}
