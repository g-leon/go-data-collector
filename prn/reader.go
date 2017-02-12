package prn

import (
	"bufio"
	"io"
	"strings"
)

// A reader reads records from a PRN file, using the
// width of the given headers as an indicator of
// maximum width for each field.
type reader struct {
	*bufio.Scanner
	position  []int
	header    string
	firstCall bool
}

func NewReader(r io.Reader, rows []string) *reader {
	cpr := &reader{
		Scanner:  bufio.NewScanner(r),
		position: make([]int, 0),
	}

	cpr.loadHeader()
	for i := 1; i < len(rows); i++ {
		cpr.position = append(cpr.position, strings.Index(cpr.header, rows[i]))
	}
	cpr.position = append(cpr.position, cpr.position[len(cpr.position)-1]+len(rows[len(rows)-1]))

	return cpr
}

func (r *reader) loadHeader() {
	r.Scanner.Scan()
	r.header = r.Scanner.Text()
}

// Read reads one record from r. The record is a slice of strings with each
// string representing one field.
func (r *reader) Read() (record []string, err error) {
	var row string
	if r.firstCall {
		row = r.header
		r.firstCall = false
	} else if r.Scanner.Scan() {
		row = r.Scanner.Text()
	}
	if err = r.Scanner.Err(); err != nil {
		return
	}
	if row == "" {
		return nil, io.EOF
	}

	record = make([]string, 0, len(r.position))
	record = append(record, strings.TrimSpace(row[0:r.position[0]]))
	record = append(record, strings.TrimSpace(row[r.position[0]:r.position[1]]))
	record = append(record, strings.TrimSpace(row[r.position[1]:r.position[2]]))
	record = append(record, strings.TrimSpace(row[r.position[2]:r.position[3]]))
	record = append(record, strings.TrimSpace(row[r.position[3]:r.position[4]]))
	record = append(record, strings.TrimSpace(row[r.position[4]:]))

	return
}
