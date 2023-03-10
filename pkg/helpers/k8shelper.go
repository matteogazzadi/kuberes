package helpers

import (
	"context"
	"log"
	"regexp"
	"sync"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// Get all Pods for a given namespace
func GetPodsByNamespace(clientset *kubernetes.Clientset, ctx context.Context, namespace string, podChan chan<- []v1.Pod, wg *sync.WaitGroup) {

	defer wg.Done()

	// Get pods
	pods, err := clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})

	if err != nil {
		log.Fatal(err)
	} else {
		podChan <- pods.Items
	}
}

// Get All Namespaces in the cluster
func GetAllNamespace(clientset *kubernetes.Clientset, ctx context.Context) ([]v1.Namespace, error) {
	namespaces, err := clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})

	if err != nil {
		return nil, err
	}

	return namespaces.Items, nil
}

// Get all Pods
func GetAllPods(excludeNamespace *[]string, matchNamespace string, clientset *kubernetes.Clientset, ctx context.Context) ([]v1.Pod, error) {

	matchNsRegex, err := regexp.Compile(matchNamespace)

	namespaces, err := GetAllNamespace(clientset, ctx)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	// Pod Channel
	podChannel := make(chan []v1.Pod, len(namespaces))

	// Wait Group
	var wg sync.WaitGroup

	// Loop on Namespace
	for _, namespace := range namespaces {

		if contains(excludeNamespace, namespace.Name) {
			continue
		}

		if matchNsRegex != nil && !matchNsRegex.MatchString(namespace.Name) {
			continue
		}

		wg.Add(1)
		go GetPodsByNamespace(clientset, ctx, namespace.Name, podChannel, &wg)
	}

	wg.Wait()
	close(podChannel)

	var pods []v1.Pod

	for podList := range podChannel {
		for _, pod := range podList {
			pods = append(pods, pod)
		}
	}

	return pods, nil
}

func contains(s *[]string, str string) bool {

	if s != nil && len(*s) > 0 {
		for _, v := range *s {
			if v == str {
				return true
			}
		}
		return false
	} else {
		return false
	}
}
