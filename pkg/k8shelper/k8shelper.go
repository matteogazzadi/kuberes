package k8shelper

import (
	"context"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// Get all Pods for a given namespace
func GetPodsByNamespace(clientset *kubernetes.Clientset, ctx context.Context, namespace string) ([]v1.Pod, error) {

	// Get pods
	pods, err := clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})

	if err != nil {
		return nil, err
	}

	return pods.Items, nil
}

// Get All Namespaces in the cluster
func GetAllNamespace(clientset *kubernetes.Clientset, ctx context.Context) ([]v1.Namespace, error) {
	namespaces, err := clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})

	if err != nil {
		return nil, err
	}

	return namespaces.Items, nil
}
