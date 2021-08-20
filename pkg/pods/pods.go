package pods

import (
	"context"
	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"strings"
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
}

type podLister struct {
	clientSet kubernetes.Interface
}

func newPodLister(clientSet kubernetes.Interface) *podLister {
	return &podLister{clientSet: clientSet}
}

func (p *podLister) get(namespace string, filters StatusFilters) []*Pod {
	podList, err := p.clientSet.
		CoreV1().
		Pods(namespace).
		List(context.TODO(), metaV1.ListOptions{LabelSelector: "app=intuit"})
	if err != nil {
		panic(err.Error())
	}
	if len(filters) != 0 {
		podFilter := &podFilter{}
		statusSpec := &statusSpecifier{filters: filters}
		return podFilter.filter(podList, statusSpec)
	}
	var result []*Pod
	for _, pod := range podList.Items {
		result = append(result, &Pod{
			Name: pod.Name,
			Status: status(&pod),
			Namespace: pod.Namespace,
		})
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