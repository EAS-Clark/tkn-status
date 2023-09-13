package helpers

import (
	"context"
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func RunTaskStatusChecker(namespace, pipelineRunName string, taskNames []string, pollingFrequency string) error {
	// Parse polling frequency duration
	pollingDuration, err := time.ParseDuration(pollingFrequency)
	if err != nil {
		return err
	}

	// Initialize the Kubernetes client
	config, err := rest.InClusterConfig()
	if err != nil {
		return err
	}
	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	// Poll the status of the tasks
	for _, taskName := range taskNames {
		for {
			pod, err := kubeClient.CoreV1().Pods(namespace).Get(context.Background(), taskName, metav1.GetOptions{})
			if err != nil {
				return err
			}

			taskRunStatus := ""
			if pod.Status.Phase == v1.PodSucceeded {
				taskRunStatus = "Succeeded"
			} else if pod.Status.Phase == v1.PodFailed {
				taskRunStatus = "Failed"
			} else {
				// Add more logic here to handle other possible states
				taskRunStatus = "Running"
			}

			if taskRunStatus != "" {
				fmt.Printf("Task '%s' status: %s\n", taskName, taskRunStatus)
				if taskRunStatus == "Succeeded" || taskRunStatus == "Failed" {
					break // Task has completed, exit loop
				}
			} else {
				fmt.Printf("Task '%s' not found in PipelineRun\n", taskName)
				break
			}

			time.Sleep(pollingDuration)
		}
	}

	return nil
}
