package jobnotify

import (
	"golang.org/x/xerrors"
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
)

func watch(resource, namespace string) (*jobResult, error) {
	configFlags := &genericclioptions.ConfigFlags{}
	restConfig, err := configFlags.ToRESTConfig()
	if err != nil {
		return nil, err
	}
	client := kubernetes.NewForConfigOrDie(restConfig)
	_, err = client.BatchV1().Jobs(namespace).Get(resource, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	watcher, err := client.BatchV1().Jobs(namespace).Watch(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	now := metav1.Now()
	for event := range watcher.ResultChan() {
		job := event.Object.(*batchv1.Job)
		jobName := job.Name
		if resource != jobName {
			continue
		}
		if job.Status.CompletionTime.Before(&now) {
			return nil, xerrors.New("specified job has already finished.")
		}
		if *job.Spec.Completions == job.Status.Succeeded+job.Status.Failed {
			return &jobResult{
				name:        resource,
				startedAt:   job.Status.StartTime,
				completedAt: job.Status.CompletionTime,
				completed:   job.Status.Succeeded,
				failed:      job.Status.Failed,
			}, nil
		}
	}
	return nil, xerrors.New("unexpected error")
}
