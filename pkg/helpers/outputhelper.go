package helpers

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"unicode/utf8"

	domain "github.com/matteogazzadi/kuberes/pkg/domain"
	"github.com/rodaine/table"
	excelize "github.com/xuri/excelize/v2"
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

	// Add Total Row
	if groupByNamespace {
		row := []string{"Total", strconv.FormatInt(total.Cpu.Request, 10), strconv.FormatInt(total.Cpu.Limit, 10), strconv.FormatInt(total.Memory.Request, 10), strconv.FormatInt(total.Memory.Limit, 10)}
		writer.Write(row)
	} else {
		row := []string{"Total", "x", strconv.FormatInt(total.Cpu.Request, 10), strconv.FormatInt(total.Cpu.Limit, 10), strconv.FormatInt(total.Memory.Request, 10), strconv.FormatInt(total.Memory.Limit, 10)}
		writer.Write(row)
	}

	writer.Flush()

	log.Println("Written .CSV file: '" + csvFilePath + "'")
}

// Write data to a .XLSX file
func WriteOutputAsXlsx(resources *[]domain.K8sStats, groupByNamespace bool, xlsxFilePath string) {

	const SHEET_NAME string = "ResourcesReport"

	f := excelize.NewFile()
	defer f.Close()

	// Set Sheet Name
	index := f.NewSheet("Sheet1")
	f.SetActiveSheet(index)
	f.SetSheetName("Sheet1", SHEET_NAME)

	// Bold style for cell
	boldStyle, err := f.NewStyle(`{"font":{"bold":true,"italic":false}}`)

	// Write Header
	if groupByNamespace {
		row := []string{"Namespace", "CPU-Request (mCore)", "CPU-Limit (mCore)", "Memory-Request (Mi)", "Memory-Limit (Mi)"}
		f.SetSheetRow(SHEET_NAME, "A1", &row)
	} else {
		row := []string{"Namespace", "Pod Name", "CPU-Request (mCore)", "CPU-Limit (mCore)", "Memory-Request (Mi)", "Memory-Limit (Mi)"}
		f.SetSheetRow(SHEET_NAME, "A1", &row)
	}

	f.SetRowStyle(SHEET_NAME, 1, 1, boldStyle)

	// Write Content
	row := 1

	for _, stats := range *resources {
		row++

		if groupByNamespace {

			// Namespace Name
			axis, _ := excelize.CoordinatesToCellName(1, row)
			f.SetCellStr(SHEET_NAME, axis, stats.Namespace)

			// CPU Request
			axis, _ = excelize.CoordinatesToCellName(2, row)
			f.SetCellInt(SHEET_NAME, axis, int(stats.Cpu.Request))

			// CPU Limit
			axis, _ = excelize.CoordinatesToCellName(3, row)
			f.SetCellInt(SHEET_NAME, axis, int(stats.Cpu.Limit))

			// Memory Request
			axis, _ = excelize.CoordinatesToCellName(4, row)
			f.SetCellInt(SHEET_NAME, axis, int(stats.Memory.Request))

			// CPU Limit
			axis, _ = excelize.CoordinatesToCellName(5, row)
			f.SetCellInt(SHEET_NAME, axis, int(stats.Memory.Limit))

		} else {

			// Namespace Name
			axis, _ := excelize.CoordinatesToCellName(1, row)
			f.SetCellStr(SHEET_NAME, axis, stats.Namespace)

			// Pod Name
			axis, _ = excelize.CoordinatesToCellName(2, row)
			f.SetCellStr(SHEET_NAME, axis, stats.PodName)

			// CPU Request
			axis, _ = excelize.CoordinatesToCellName(3, row)
			f.SetCellInt(SHEET_NAME, axis, int(stats.Cpu.Request))

			// CPU Limit
			axis, _ = excelize.CoordinatesToCellName(4, row)
			f.SetCellInt(SHEET_NAME, axis, int(stats.Cpu.Limit))

			// Memory Request
			axis, _ = excelize.CoordinatesToCellName(5, row)
			f.SetCellInt(SHEET_NAME, axis, int(stats.Memory.Request))

			// CPU Limit
			axis, _ = excelize.CoordinatesToCellName(6, row)
			f.SetCellInt(SHEET_NAME, axis, int(stats.Memory.Limit))
		}
	}

	// === Write TOTAL Row ===
	lastDataRow := row
	row++

	if groupByNamespace {

		// ==== ADD PATTERN ====
		patternStyle, _ := f.NewStyle(&excelize.Style{
			Fill: excelize.Fill{Type: "pattern", Pattern: 3, Color: []string{"000000"}},
		})

		for i := 1; i < 6; i++ {
			axis, _ := excelize.CoordinatesToCellName(i, row)
			f.SetCellStyle(SHEET_NAME, axis, axis, patternStyle)
		}

		// ==== ADD Data ====
		row++

		// Namespace Name
		axis, _ := excelize.CoordinatesToCellName(1, row)
		f.SetCellStr(SHEET_NAME, axis, "Total")
		f.SetCellStyle(SHEET_NAME, axis, axis, boldStyle)

		// CPU Request
		axis, _ = excelize.CoordinatesToCellName(2, row)
		f.SetCellFormula(SHEET_NAME, axis, "=SUM(B2:B"+strconv.Itoa(lastDataRow)+")")
		f.SetCellStyle(SHEET_NAME, axis, axis, boldStyle)

		// CPU Limit
		axis, _ = excelize.CoordinatesToCellName(3, row)
		f.SetCellFormula(SHEET_NAME, axis, "=SUM(C2:C"+strconv.Itoa(lastDataRow)+")")
		f.SetCellStyle(SHEET_NAME, axis, axis, boldStyle)

		// Memory Request
		axis, _ = excelize.CoordinatesToCellName(4, row)
		f.SetCellFormula(SHEET_NAME, axis, "=SUM(D2:D"+strconv.Itoa(lastDataRow)+")")
		f.SetCellStyle(SHEET_NAME, axis, axis, boldStyle)

		// CPU Limit
		axis, _ = excelize.CoordinatesToCellName(5, row)
		f.SetCellFormula(SHEET_NAME, axis, "=SUM(E2:E"+strconv.Itoa(lastDataRow)+")")
		f.SetCellStyle(SHEET_NAME, axis, axis, boldStyle)
	} else {

		// ==== ADD PATTERN ====
		patternStyle, _ := f.NewStyle(&excelize.Style{
			Fill: excelize.Fill{Type: "pattern", Pattern: 3, Color: []string{"000000"}},
		})

		for i := 1; i < 7; i++ {
			axis, _ := excelize.CoordinatesToCellName(i, row)
			f.SetCellStyle(SHEET_NAME, axis, axis, patternStyle)
		}

		// ==== ADD Data ====
		row++

		// Namespace Name
		axis, _ := excelize.CoordinatesToCellName(1, row)
		f.SetCellStr(SHEET_NAME, axis, "Total")
		f.SetCellStyle(SHEET_NAME, axis, axis, boldStyle)

		// Pod Name
		axis, _ = excelize.CoordinatesToCellName(2, row)
		f.SetCellStr(SHEET_NAME, axis, "x")
		f.SetCellStyle(SHEET_NAME, axis, axis, boldStyle)

		// CPU Request
		axis, _ = excelize.CoordinatesToCellName(2, row)
		f.SetCellFormula(SHEET_NAME, axis, "=SUM(C2:C"+strconv.Itoa(lastDataRow)+")")
		f.SetCellStyle(SHEET_NAME, axis, axis, boldStyle)

		// CPU Limit
		axis, _ = excelize.CoordinatesToCellName(3, row)
		f.SetCellFormula(SHEET_NAME, axis, "=SUM(D2:D"+strconv.Itoa(lastDataRow)+")")
		f.SetCellStyle(SHEET_NAME, axis, axis, boldStyle)

		// Memory Request
		axis, _ = excelize.CoordinatesToCellName(4, row)
		f.SetCellFormula(SHEET_NAME, axis, "=SUM(E2:E"+strconv.Itoa(lastDataRow)+")")
		f.SetCellStyle(SHEET_NAME, axis, axis, boldStyle)

		// CPU Limit
		axis, _ = excelize.CoordinatesToCellName(5, row)
		f.SetCellFormula(SHEET_NAME, axis, "=SUM(F2:F"+strconv.Itoa(lastDataRow)+")")
		f.SetCellStyle(SHEET_NAME, axis, axis, boldStyle)
	}

	// Auto-Fit Columnd Width
	// Ref to: https://github.com/qax-os/excelize/issues/92#issuecomment-821578446
	cols, err := f.GetCols(SHEET_NAME)
	if err != nil {
		log.Fatal(err)
	}
	for idx, col := range cols {
		largestWidth := 0
		for _, rowCell := range col {
			cellWidth := utf8.RuneCountInString(rowCell) + 2 // + 2 for margin
			if cellWidth > largestWidth {
				largestWidth = cellWidth
			}
		}
		name, err := excelize.ColumnNumberToName(idx + 1)
		if err != nil {
			log.Fatal(err)
		}
		f.SetColWidth(SHEET_NAME, name, name, float64(largestWidth))
	}

	// Save file
	f.SaveAs(xlsxFilePath)

	log.Println("Written .XLSX file: '" + xlsxFilePath + "'")
}
