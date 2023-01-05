package main

import (
	"context"
	"log"
	"sort"

	domain "github.com/matteogazzadi/kuberes/pkg/domain"
	k8shelper "github.com/matteogazzadi/kuberes/pkg/k8shelper"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/fatih/color"
	"github.com/rodaine/table"
)

// Main function - Application Entry Point
func main() {

	// 1. Initialize Kubernets Clients
	ctx := context.Background()
	config := ctrl.GetConfigOrDie()
	clientset := kubernetes.NewForConfigOrDie(config)

	// Resources Map object
	resources := make(map[string]*domain.PodStats)

	// Retrieve the list of namespaces in cluster
	namespaces, err := k8shelper.GetAllNamespace(clientset, ctx)

	if err != nil {
		panic(err)
	}

	// Loop all namespace and calculate POD resources
	for _, namespace := range namespaces {
		pods, err := k8shelper.GetPodsByNamespace(clientset, ctx, namespace.Name)

		if err != nil {
			log.Fatal(err)
			continue
		}

		// Check if entry for the given namespace exist in the local map.
		_, ok := resources[namespace.Name]

		// If not present, add it.
		if !ok {

			var newStats domain.PodStats

			newStats.Namespace = namespace.Name
			resources[namespace.Name] = &newStats
		}

		// Loop on pods
		for _, pod := range pods {

			stats, _ := resources[namespace.Name]

			// Loop On Containers
			for _, container := range pod.Spec.Containers {

				// CPU
				cpuRequest := container.Resources.Requests.Cpu().MilliValue()
				cpuLimit := container.Resources.Limits.Cpu().MilliValue()

				// Memory
				memRequest, _ := container.Resources.Requests.Memory().AsInt64()
				memLimit, _ := container.Resources.Limits.Memory().AsInt64()

				// Convert MB to Mib
				memRequest = memRequest / 1048576
				memLimit = memLimit / 1048576

				stats.Cpu.Limit += cpuLimit
				stats.Cpu.Request += cpuRequest

				stats.Memory.Limit += memLimit
				stats.Memory.Request += memRequest
			}
		}

	}

	// Sort the resources key alphabetically
	keys := make([]string, 0, len(resources))
	for k := range resources {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Generate the on screen Table
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	tbl := table.New("Namespace", "CPU (Request)", "CPU (Limit)", "Memory (Request)", "Memory (Limit)")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for _, k := range keys {
		stats := resources[k]
		tbl.AddRow(stats.Namespace, stats.Cpu.Request, stats.Cpu.Limit, stats.Memory.Request, stats.Memory.Limit)
	}

	tbl.Print()
}
