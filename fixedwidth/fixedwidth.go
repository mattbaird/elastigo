package fixedwidth

import (
	"bytes"
	"fmt"
	"log"
	"regexp"
	"strings"
)

type Field struct {
	Header string
	Width  int
}

func (f *Field) String() string {
	return fmt.Sprintf("Field{Header:\"%s\", Width:%d}", f.Header, f.Width)
}

type Column struct {
	FieldInfo Field
	Data      []string
}

func (c *Column) String() string {
	d := ""
	if len(c.Data) > 0 {
		d = fmt.Sprint("'", strings.Join(c.Data, "', '"), "'")
	}
	return fmt.Sprintf("Column{FieldInfo:%s, Data:[%s]}", c.FieldInfo.String(), d)
}

type FixedWidthTable []Column

var tokenRegexp *regexp.Regexp

// getTokenRegexp returns the compiled Regexp object corresponding
// to a single field header with trailing whitespace. This method
// caches the compiled Regexp instance in the field tokenRegexp.
func getTokenRegexp() *regexp.Regexp {

	if tokenRegexp == nil {
		t, err := regexp.Compile("(\\S+(?:\\s+)?)")
		if err != nil {
			log.Fatal(err)
		}
		tokenRegexp = t
	}
	return tokenRegexp
}

// parseHeaderLine assumes that the input is the header line of a fixed
// width table, whose header names do not contain whitespace. It
// returns the parsed header structure with the width of the last
// field being -1.
func parseHeaderLine(line string) (fields []Field) {

	fields = []Field{}

	for _, token := range getTokenRegexp().FindAllString(line, -1) {
		fields = append(fields, Field{strings.Trim(token, " "), len(token)})
	}
	if len(fields) > 0 {
		fields[len(fields)-1].Width = -1
	}

	return
}

// parseDataLines treats the specified lines as the data and the specified
// fields as the headers for the data. One column is built and returned for
// each field whose data vector is populated with width-parsed values from
// the corresponding position in each line.
func parseDataLines(lines []string, fields []Field) (t FixedWidthTable, err error) {

	// Initialize a new Column slice
	t = []Column{}
	for _, f := range fields {
		t = append(t, Column{FieldInfo: f, Data: []string{}})
	}

	// Parse columns, removing the values from each line
	for c, f := range fields {

		w := f.Width
		for r := 0; r < len(lines); r++ {
			
			if len(lines[r]) < 1 {
				continue
			}
			
			line := lines[r]
			if w > 0 && len(line) > w {
				t[c].Data = append(t[c].Data, strings.Trim(line[0:w], " "))
				lines[r] = line[w:]
			} else {
				t[c].Data = append(t[c].Data, line)
			}
		}
	}

	return t, nil
}

func (t *FixedWidthTable) Width() int {

	w := 0
	if t != nil {
		w = len(*t)
	}
	return w
}

func (tp *FixedWidthTable) Height() int {

	t := *tp
	h := 0
	if t != nil && len(t) > 0 {
		h = len(t[0].Data)
	}

	return h
}

// Item returns the data item at the specified row and column from the
// given FixedWidthTable using zero-based indexing. If the specified
// indices are out of bounds, an empty string is returned.
func (tp *FixedWidthTable) Item(row int, col int) string {

	t := *tp

	switch {
	case t == nil || len(t) < 1 || len(t[0].Data) < 1:
		//log.Fatalf("Input c is nil, empty or has no Data")
		return ""
	case col < 0 || col >= len(t):
		//log.Fatalf("Index col (%d) is out of range [0,%d]", col, len(t)-1)
		return ""
	case row < 0 || row >= len(t[0].Data):
		//log.Fatalf("Index row (%d) is out of range [0,%d]", row, len(t[0].Data)-1)
		return ""
	}

	return t[col].Data[row]
}

// RowMap returns a map containing all of the rows data, indexed by
// the header name for each column. If the table is empty, or the
// row is out of bounds, an empty map is returned.
func (tp *FixedWidthTable) RowMap(row int) map[string]string {

	t := *tp
	m := map[string]string{}
	if row < 0 || row >= t.Height() || t.Width() < 1 {
		return m
	}

	for col, c := range t {
		m[c.FieldInfo.Header] = t.Item(row, col)
	}

	return m
}

// String returns the string representation of the fixed width table,
// basically the inverse of NewFixedWidthTable().
func (tp *FixedWidthTable) String() string {

	t := *tp
	var buffer bytes.Buffer
	var pat = make(map[int]string)

	for col, c := range t {
		pat[col] = fmt.Sprint("%-", c.FieldInfo.Width, "s")
		buffer.WriteString(fmt.Sprintf(pat[col], c.FieldInfo.Header))
	}
	buffer.WriteString("\n")

	for r := 0; r < t.Height(); r++ {

		for c, col := range t {
			buffer.WriteString(fmt.Sprintf(pat[c], col.Data[r]))
		}
		if r < t.Height()-1 {
			buffer.WriteString("\n")
		}
	}

	return buffer.String()
}

// NewFixedWidthTable assumes that the input is a newline-delimited
// set of lines, the first of which is a space-delimited fixed width
// header line. The trailing spaces in each header field are assumed to
// be included in the width of that field. The remaining lines should 
// contain data for each field in columns of the same width. For
// example:
//  
//   field1    f2 field3  field4
//   data1     d2 data3   data four
//   data fivesd6 datasevsdata 8888 888 88
//
func NewFixedWidthTable(data []byte) (*FixedWidthTable, error) {

	lines := strings.Split(string(data[:]), "\n")
	if len(lines) < 2 {
		return nil, fmt.Errorf("At least two input lines are required, recieved %d", len(lines))
	}
	
	f := parseHeaderLine(lines[0])
	t, err := parseDataLines(lines[1:], f)
	if err != nil {
		return nil, err
	}

	return &t, nil
}
