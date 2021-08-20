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

func TestPodLister_GetTerminatingPods(t *testing.T) {
	mockLister := &podLister{
		clientSet: fake.NewSimpleClientset(),
	}
	podName := "pod-terminating"
	pod := v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      podName,
			Namespace: "default",
			Labels: map[string]string{
				"app": "intuit",
			},
			DeletionTimestamp: &metav1.Time{Time: time.Now()},
		},
		Status: v1.PodStatus{
			Phase:     "Running",
			Reason:    "",
			StartTime: &metav1.Time{Time: time.Now()},
		},
	}
	_, err := mockLister.clientSet.
		CoreV1().
		Pods("default").
		Create(context.TODO(), &pod, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	result := mockLister.get("default", nil, nil)
	assert.True(t, len(result) == 1, "Only one pod expected in terminating state")
	assert.Equal(t, result[0].Name, podName, "Unexpected pod name")
}
