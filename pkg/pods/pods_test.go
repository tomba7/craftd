package pods

import (
	"context"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	"testing"
	"time"
)

func newPod(name, namespace string, phase v1.PodPhase) *v1.Pod {
	pod := v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: "default",
			Labels: map[string]string{
				"app": "intuit",
			},
		},
		Status: v1.PodStatus{
			Phase:     v1.PodRunning,
			Reason:    "",
			StartTime: &metav1.Time{Time: time.Now()},
		},
	}
	return &pod
}

func TestPodLister_GetTerminatingPods(t *testing.T) {
	mockLister := &podLister{
		clientSet: fake.NewSimpleClientset(),
	}
	podName := "pod-terminating"
	pod := newPod(podName, "default", v1.PodRunning)
	pod.ObjectMeta.DeletionTimestamp = &metav1.Time{Time: time.Now()}
	_, err := mockLister.clientSet.
		CoreV1().
		Pods("default").
		Create(context.TODO(), pod, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	result := mockLister.get("default", nil, nil)
	assert.True(t, len(result) == 1, "Only one pod expected in terminating state")
	assert.Equal(t, result[0].Name, podName, "Unexpected pod name")
}

func TestPodLister_GetCompletedPods(t *testing.T) {
	mockLister := &podLister{
		clientSet: fake.NewSimpleClientset(),
	}
	podName := "pod-completed"
	pod := newPod(podName, "default", v1.PodSucceeded)
	_, err := mockLister.clientSet.
		CoreV1().
		Pods("default").
		Create(context.TODO(), pod, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	result := mockLister.get("default", nil, nil)
	assert.True(t, len(result) == 1, "Only one pod expected in terminating state")
	assert.Equal(t, result[0].Name, podName, "Unexpected pod name")
}

/*
func TestPodLister_GetTimeElapsedPods(t *testing.T) {
	mockLister := &podLister{
		clientSet: fake.NewSimpleClientset(),
	}
	podName := "pod-elapsed"
	minutes := 10
	elapsedDuration := time.Minute * time.Duration(minutes)
	pod := newPod(podName, "default", v1.PodSucceeded)
	pod.ObjectMeta.SetCreationTimestamp(metav1.Time{
		Time: time.Now().Truncate(elapsedDuration),
	})
	fmt.Println(time.Now())
	fmt.Println(pod.ObjectMeta.CreationTimestamp.Time)
	_, err := mockLister.clientSet.
		CoreV1().
		Pods("default").
		Create(context.TODO(), pod, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	minutes++
	result := mockLister.get("default", nil, &minutes)
	assert.True(t, len(result) == 1, "Only one pod expected in terminating state")
	assert.Equal(t, result[0].Name, podName, "Unexpected pod name")
}
*/