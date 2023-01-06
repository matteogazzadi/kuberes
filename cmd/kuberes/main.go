package main

import (
	"context"
	"fmt"
	"sort"
	"time"

	calc "github.com/matteogazzadi/kuberes/pkg/calculator"
	domain "github.com/matteogazzadi/kuberes/pkg/domain"
	k8shelper "github.com/matteogazzadi/kuberes/pkg/k8shelper"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/rodaine/table"
)

// Main function - Application Entry Point
func main() {

	startTime := time.Now().UTC()

	// 1. Initialize Kubernets Clients
	ctx := context.Background()
	config := ctrl.GetConfigOrDie()
	clientset := kubernetes.NewForConfigOrDie(config)

	// Get All Pods in cluster
	pods, err := k8shelper.GetAllPods(clientset, ctx)

	if err != nil {
		panic(err)
	}

	var resources []domain.K8sStats

	groupByNamespace := true

	// Calculate Resources
	calc.CalculateResources(groupByNamespace, &pods, &resources)

	// Sort Resources by Namespace
	sort.Slice(resources, func(i, j int) bool {
		return resources[i].Namespace < resources[j].Namespace
	})

	// Table output
	var tbl table.Table

	// Generate the on screen Table
	if groupByNamespace {
		tbl = table.New("Namespace", "CPU-Request (mCore)", "CPU-Limit (mCore)", "Memory-Request (Mi)", "Memory-Limit (Mi)")
	} else {
		tbl = table.New("Namespace", "Pod Name", "CPU-Request (mCore)", "CPU-Limit (mCore)", "Memory-Request (Mi)", "Memory-Limit (Mi)")
	}

	for _, stats := range resources {
		if groupByNamespace {
			tbl.AddRow(stats.Namespace, stats.Cpu.Request, stats.Cpu.Limit, stats.Memory.Request, stats.Memory.Limit)
		} else {
			tbl.AddRow(stats.Namespace, stats.PodName, stats.Cpu.Request, stats.Cpu.Limit, stats.Memory.Request, stats.Memory.Limit)
		}

	}
	tbl.Print()

	fmt.Println()
	fmt.Println("Elapsed: ", time.Since(startTime).Milliseconds(), "ms")
}
