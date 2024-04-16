package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"text/template"
)

type Row struct {
	Address string
	Values  []string
}

type Sheet struct {
	SheetName       string
	ColumnAddresses []string
	RowAddresses    []string
	Cells           map[string]string
}

func (sh Sheet) CellsByRow() []Row {
	var cells []Row
	for _, rowAddr := range sh.RowAddresses {
		var rowCells []string
		for _, colAddr := range sh.ColumnAddresses {
			rowCells = append(rowCells, sh.Cells[colAddr+rowAddr])
		}
		newRow := Row{Address: rowAddr, Values: rowCells}
		cells = append(cells, newRow)
	}
	fmt.Printf("%v\n", cells)
	return cells
}

func get_sheets(w http.ResponseWriter, _ *http.Request) {
	tmpl, err := template.New("sheet.tmpl").ParseFiles("sheet.tmpl")
	if err != nil {
		panic(err)
	}
	testSheet := Sheet{
		SheetName: "HondaCivic",
		ColumnAddresses: []string{
			"A",
			"B",
			"C",
		},
		RowAddresses: []string{
			"1",
			"2",
			"3",
		},
	}
	testSheet.Cells = make(map[string]string)

	testSheet.Cells["A1"] = "Al1ce"
	testSheet.Cells["A2"] = "Ali2e"
	testSheet.Cells["A3"] = "Alic3"

	testSheet.Cells["B1"] = "1ob"
	testSheet.Cells["B2"] = "B2b"
	testSheet.Cells["B3"] = "Bo3"

	testSheet.Cells["C1"] = "Cr1ig"
	testSheet.Cells["C2"] = "Cra2g"
	testSheet.Cells["C3"] = "Crai3"

	err = tmpl.Execute(w, testSheet)
	if err != nil {
		panic(err)
	}
}

func get_sheet_by_name(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "<body><h1>"+strings.TrimPrefix(r.URL.Path, "/")+`</h1>`)
}

func main() {
	http.HandleFunc("/", get_sheets)
	http.HandleFunc("/{sheetName}", get_sheet_by_name)
	fmt.Fprintf(os.Stderr, "%v\n", http.ListenAndServe(":8080", nil))
}
