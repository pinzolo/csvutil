package csvutil

import (
	"io"
	"sort"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

const (
	// SortDataTypeText is used when want use value as text.
	SortDataTypeText = "text"
	// SortDataTypeNumber is used when want use value as number.
	SortDataTypeNumber = "number"
	// EmptyNatural is used when empty string is ordered in natural.
	EmptyNatural = "natural"
	// EmptyFirst is used when empty string is ordered in first.
	EmptyFirst = "first"
	// EmptyLast is used when empty string is ordered in last.
	EmptyLast = "last"
)

var (
	supportedSortDataTypes  = []string{SortDataTypeText, SortDataTypeNumber}
	supportedEmptyHandlings = []string{EmptyNatural, EmptyFirst, EmptyLast}
)

// SortOption is option holder for Sort.
type SortOption struct {
	// Source file does not have header line. (default false)
	NoHeader bool
	// Encoding of source file. (default utf8)
	Encoding string
	// Encoding for output.
	OutputEncoding string
	// Column symbol of target column
	Column string
	// DataType is sort key's data type
	DataType string
	// Sort in descending order
	Descending bool
	// Handling method of empty string
	EmptyHandling string
}

func (o SortOption) outputEncoding() string {
	if o.OutputEncoding != "" {
		return o.OutputEncoding
	}
	return o.Encoding
}

func (o SortOption) validate() error {
	if o.Column == "" {
		return errors.New("no column")
	}
	if o.NoHeader {
		if !isDigit(o.Column) {
			return errors.New("not number column symbol")

		}
	}
	if o.DataType != "" && !containsString(supportedSortDataTypes, o.DataType) {
		return errors.Errorf("unsupported sort data type: %s", o.DataType)
	}
	if o.EmptyHandling != "" && !containsString(supportedEmptyHandlings, o.EmptyHandling) {
		return errors.Errorf("unsupported empty handling: %s", o.EmptyHandling)
	}
	return nil
}

// Sort CSV.
func Sort(r io.Reader, w io.Writer, o SortOption) error {
	if err := o.validate(); err != nil {
		return errors.Wrap(err, "invalid option")
	}

	cr, bom := reader(r, o.Encoding)
	cw := writer(w, bom, o.outputEncoding())
	defer cw.Flush()

	recs, err := cr.ReadAll()
	if err != nil {
		return err
	}
	if len(recs) == 0 {
		return errors.New("empty CSV")
	}

	var data [][]string
	var col *column
	if o.NoHeader {
		col = newColumnWithIndex(o.Column, nil)
		data = recs
	} else {
		col = newColumnWithIndex(o.Column, recs[0])
		cw.Write(recs[0])
		data = recs[1:]
	}

	if o.DataType == SortDataTypeNumber {
		for _, rec := range data {
			if rec[col.index] == "" {
				continue
			}
			_, err := strconv.ParseFloat(rec[col.index], 64)
			if err != nil {
				return err
			}
		}
	}

	data = sortCSVData(data, col, o)
	return cw.WriteAll(data)
}

func compareStringsAsNumber(s1 string, s2 string) float64 {
	var f1, f2 float64
	if s1 != "" {
		f1, _ = strconv.ParseFloat(s1, 64)
	}
	if s2 != "" {
		f2, _ = strconv.ParseFloat(s2, 64)
	}

	return f1 - f2
}

func sortCSVData(data [][]string, col *column, o SortOption) [][]string {
	sort.Slice(data, func(i, j int) bool {
		si := data[i][col.index]
		sj := data[j][col.index]
		if o.EmptyHandling != EmptyNatural {
			if si == "" {
				return o.EmptyHandling == EmptyFirst
			} else if sj == "" {
				return o.EmptyHandling == EmptyLast
			}
		}
		if o.Descending {
			if o.DataType == SortDataTypeText {
				return strings.Compare(si, sj) > 0
			}
			return compareStringsAsNumber(si, sj) > 0.0
		}
		if o.DataType == SortDataTypeText {
			return strings.Compare(si, sj) < 0
		}
		return compareStringsAsNumber(si, sj) < 0.0
	})
	return data
}
