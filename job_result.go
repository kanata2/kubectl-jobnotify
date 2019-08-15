package jobnotify

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type jobResult struct {
	name        string
	startedAt   *metav1.Time
	completedAt *metav1.Time
	completed   int32
	failed      int32
}
