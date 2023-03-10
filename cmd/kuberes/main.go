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
var outputAsXlsx bool
var outputFilePath string
var excludeNamespace []string
var matchNamespace string

// Init Function - Arguments parsing
func init() {
	var outputFormat string
	var excludeNs string

	flag.BoolVar(&groupByNamespace, "group-by-ns", true, "Should group statistics by namespace ?")
	flag.StringVar(&outputFormat, "output", "table", "Output type. Valid values are: table,csv,xlsx")
	flag.StringVar(&outputFilePath, "file-path", "", "Full Path to the .CSV/.XLSX File to produce")
	flag.StringVar(&excludeNs, "exclude-ns", "", "Namespaces Names to be excluded separated by ,")
	flag.StringVar(&matchNamespace, "match-ns-regex", "", "Namespaces Names to be matched on the given RegEx")
	flag.Parse()

	// Check if output parameter is valid
	switch strings.ToLower(outputFormat) {
	case "table":
		outputAsTable = true
		outputAsCsv = false
		outputAsXlsx = false
	case "csv":
		outputAsTable = false
		outputAsCsv = true
		outputAsXlsx = false
	case "xlsx":
		outputAsTable = false
		outputAsCsv = false
		outputAsXlsx = true
	default:
		panic("Unrecognized 'output' value '" + outputFormat + "'")
	}

	// Check if output file Path is valid (only if output is CSV or XLSX)
	if (outputAsCsv || outputAsXlsx) && outputFilePath == "" {
		panic("Output file path is not set")
	}

	if excludeNs != "" {
		excludeNamespace = strings.Split(excludeNs, ",")
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
	pods, err := helper.GetAllPods(&excludeNamespace, matchNamespace, clientset, ctx)

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
		if outputAsCsv && outputFilePath != "" {
			helper.WriteOutputAsCsv(&resources, groupByNamespace, outputFilePath)
		} else {
			if outputAsXlsx && outputFilePath != "" {
				helper.WriteOutputAsXlsx(&resources, groupByNamespace, outputFilePath)
			}
		}
	}

	// Report Elapsed Time
	fmt.Println()
	fmt.Println("Elapsed: ", time.Since(startTime).Milliseconds(), "ms")
}
