package helpers

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"

	domain "github.com/matteogazzadi/kuberes/pkg/domain"
	"github.com/rodaine/table"
)

// Write data to Output STDOUT as Table
func WriteOutputAsTable(resources *[]domain.K8sStats, groupByNamespace bool) {

	// Table output
	var tbl table.Table

	// Generate the on screen Table
	if groupByNamespace {
		tbl = table.New("Namespace", "CPU-Request (mCore)", "CPU-Limit (mCore)", "Memory-Request (Mi)", "Memory-Limit (Mi)")
	} else {
		tbl = table.New("Namespace", "Pod Name", "CPU-Request (mCore)", "CPU-Limit (mCore)", "Memory-Request (Mi)", "Memory-Limit (Mi)")
	}

	var total domain.K8sStats

	for _, stats := range *resources {
		if groupByNamespace {
			tbl.AddRow(stats.Namespace, stats.Cpu.Request, stats.Cpu.Limit, stats.Memory.Request, stats.Memory.Limit)
		} else {
			tbl.AddRow(stats.Namespace, stats.PodName, stats.Cpu.Request, stats.Cpu.Limit, stats.Memory.Request, stats.Memory.Limit)
		}

		total.Cpu.Limit += stats.Cpu.Limit
		total.Cpu.Request += stats.Cpu.Request
		total.Memory.Limit += stats.Memory.Limit
		total.Memory.Request += stats.Memory.Request

	}

	// Add Total Row
	if groupByNamespace {
		tbl.AddRow("------", "------", "------", "------", "------")
		tbl.AddRow("Total", total.Cpu.Request, total.Cpu.Limit, total.Memory.Request, total.Memory.Limit)
	} else {
		tbl.AddRow("------", "------", "------", "------", "------", "------")
		tbl.AddRow("Total", "x", total.Cpu.Request, total.Cpu.Limit, total.Memory.Request, total.Memory.Limit)
	}

	tbl.Print()
}

// Write data to a .CSV file
func WriteOutputAsCsv(resources *[]domain.K8sStats, groupByNamespace bool, csvFilePath string) {

	csvFile, err := os.Create(csvFilePath)

	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	defer csvFile.Close()

	writer := csv.NewWriter(csvFile)

	// Write Header
	if groupByNamespace {
		row := []string{"Namespace", "CPU-Request (mCore)", "CPU-Limit (mCore)", "Memory-Request (Mi)", "Memory-Limit (Mi)"}
		writer.Write(row)
	} else {
		row := []string{"Namespace", "Pod Name", "CPU-Request (mCore)", "CPU-Limit (mCore)", "Memory-Request (Mi)", "Memory-Limit (Mi)"}
		writer.Write(row)
	}

	// Write Content
	var total domain.K8sStats

	for _, stats := range *resources {
		if groupByNamespace {
			row := []string{stats.Namespace, strconv.FormatInt(stats.Cpu.Request, 10), strconv.FormatInt(stats.Cpu.Limit, 10), strconv.FormatInt(stats.Memory.Request, 10), strconv.FormatInt(stats.Memory.Limit, 10)}
			writer.Write(row)
		} else {
			row := []string{stats.Namespace, stats.PodName, strconv.FormatInt(stats.Cpu.Request, 10), strconv.FormatInt(stats.Cpu.Limit, 10), strconv.FormatInt(stats.Memory.Request, 10), strconv.FormatInt(stats.Memory.Limit, 10)}
			writer.Write(row)
		}

		total.Cpu.Limit += stats.Cpu.Limit
		total.Cpu.Request += stats.Cpu.Request
		total.Memory.Limit += stats.Memory.Limit
		total.Memory.Request += stats.Memory.Request

	}

	writer.Flush()
}
