package main

import (
	"context"
	"flag"
	"fmt"
	"sort"
	"strings"
	"time"

	calc "github.com/matteogazzadi/kuberes/pkg/calculator"
	domain "github.com/matteogazzadi/kuberes/pkg/domain"
	helper "github.com/matteogazzadi/kuberes/pkg/helpers"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
)

// Arguments Parameter
var groupByNamespace bool
var outputAsTable bool
var outputAsCsv bool
var csvOutputFilePath string

// Init Function - Arguments parsing
func init() {
	var outputFormat string

	flag.BoolVar(&groupByNamespace, "group-by-ns", true, "Should group statistics by namespace ?")
	flag.StringVar(&outputFormat, "output", "table", "Output type. Valid values are: table,csv")
	flag.StringVar(&csvOutputFilePath, "csv-path", "", "Full Path to the .CSV File to produce")
	flag.Parse()

	// Check if output parameter is valid
	switch strings.ToLower(outputFormat) {
	case "table":
		outputAsTable = true
		outputAsCsv = false
	case "csv":
		outputAsTable = false
		outputAsCsv = true
	default:
		panic("Unrecognized 'output' value '" + outputFormat + "'")
	}

	// Check if CSV output Path is valid (only if output is CSV)
	if outputAsCsv && csvOutputFilePath == "" {
		panic("CSV Output path is not set")
	}
}

// Main function - Application Entry Point
func main() {

	startTime := time.Now().UTC()

	// 1. Initialize Kubernets Clients
	ctx := context.Background()
	config := ctrl.GetConfigOrDie()
	clientset := kubernetes.NewForConfigOrDie(config)

	// Get All Pods in cluster
	pods, err := helper.GetAllPods(clientset, ctx)

	if err != nil {
		panic(err)
	}

	var resources []domain.K8sStats

	// Calculate Resources
	calc.CalculateResources(groupByNamespace, &pods, &resources)

	// Sort Resources by Namespace
	sort.Slice(resources, func(i, j int) bool {
		return resources[i].Namespace < resources[j].Namespace
	})

	// =============== //
	// Generate Output //
	// =============== //

	if outputAsTable {
		// Table Output
		helper.WriteOutputAsTable(&resources, groupByNamespace)
	} else {
		if outputAsCsv && csvOutputFilePath != "" {
			helper.WriteOutputAsCsv(&resources, groupByNamespace, csvOutputFilePath)
		}
	}

	// Report Elapsed Time
	fmt.Println()
	fmt.Println("Elapsed: ", time.Since(startTime).Milliseconds(), "ms")
}
