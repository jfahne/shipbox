package main

import (
	//"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"text/template"

	_ "github.com/mattn/go-sqlite3"
)

const CREATE_TEST_TABLE = `
CREATE TABLE box
('columnAddress', 'rowAddress', 'rawValue')
`

const POPULATE_TEST_TABLE = `
INSERT INTO box 
VALUES 
('A', '1', 'SCOOPITY'), 
('A', '2', 'WOOPITY'), 
('A', '3', 'Poopity')
`

const DROP_TEST_TABLE = `
DROP TABLE box
`

type Addressable interface {
	Address() string
}

type Assessable interface {
	Raw() string
	Assess() interface{}
}

type BoxLike interface {
	Addressable
	Assessable
}

type Addressed struct {
	address string
}

func (adr Addressed) Address() string { return adr.address }

type Assessed struct {
	rawValue      string
	assessedValue interface{}
}

func (ass Assessed) Raw() string         { return ass.rawValue }
func (ass Assessed) Assess() interface{} { return ass.assessedValue }

type Box struct {
	Addressed
	Assessed
}

type BoxQueryRow struct {
	columnAddress string
	rowAddress    string
	rawValue      string
}

type ScoopCell struct {
	Box
}

func (scp ScoopCell) Address() string {
	return scp.Address()
}

func (scp ScoopCell) Assess() interface{} {
	return scp.Raw() + " POOP"
}

type ShipBox struct {
	Box
	columnAddress string
	rowAddress    string
}

func (shbx ShipBox) ColumnAddress() string { return shbx.columnAddress }
func (shbx ShipBox) RowAddress() string    { return shbx.rowAddress }
func (shbx ShipBox) Address() string       { return shbx.columnAddress + shbx.rowAddress }
func (shbx ShipBox) Assess() interface{}   { return 25 }

func simpleHandler(w http.ResponseWriter, req *http.Request) {
	conn, err := sql.Open("sqlite3", "box.db")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	_, err = conn.Exec(CREATE_TEST_TABLE)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Create Table failed: %v\n", err)
		os.Exit(1)
	}
	_, err = conn.Exec(POPULATE_TEST_TABLE)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Populate Table failed: %v\n", err)
		os.Exit(1)
	}

	var queryRow BoxQueryRow
	var boxes []BoxLike
	rows, err := conn.Query("select * from box")
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&queryRow.columnAddress, &queryRow.rowAddress, &queryRow.rawValue)
		fmt.Println("Query Row")
		fmt.Printf("%v\n", queryRow)

		var box ShipBox
		box.columnAddress = queryRow.columnAddress
		box.rowAddress = queryRow.rowAddress
		box.address = box.Address()
		box.rawValue = queryRow.rawValue
		box.assessedValue = box.Assess()
		boxes = append(boxes, box)
	}
	fmt.Println("Box")
	fmt.Printf("%v\n", boxes)

	tmpl, err := template.New("box.tmpl").ParseFiles("box.tmpl")
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(w, boxes[0])
	if err != nil {
		panic(err)
	}
	_, err = conn.Exec(DROP_TEST_TABLE)
}

//func main() {
//	http.HandleFunc("/", simpleHandler)
//	fmt.Fprintf(os.Stderr, "%v\n", http.ListenAndServe(":8080", nil))
//}
