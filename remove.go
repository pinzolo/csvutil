package csvutil

import (
	"io"

	"github.com/pkg/errors"
)

// RemoveOption is option holder for Remove.
type RemoveOption struct {
	// Source file does not have header line. (default false)
	NoHeader bool
	// Encoding of source file. (default utf8)
	Encoding string
	// Encoding for output.
	OutputEncoding string
	// ColumnSyms header or column index list.
	ColumnSyms []string
}

func (o RemoveOption) validate() error {
	if len(o.ColumnSyms) == 0 {
		return errors.New("no column")
	}
	if o.NoHeader {
		for _, c := range o.ColumnSyms {
			if !isDigit(c) {
				return errors.New("not number column symbol")
			}
		}
	}
	return nil
}

func (o RemoveOption) outputEncoding() string {
	if o.OutputEncoding != "" {
		return o.OutputEncoding
	}
	return o.Encoding
}

// Remove column(s) in CSV.
func Remove(r io.Reader, w io.Writer, o RemoveOption) error {
	if err := o.validate(); err != nil {
		return errors.Wrap(err, "invalid option")
	}

	cr, bom := reader(r, o.Encoding)
	cw := writer(w, bom, o.outputEncoding())
	defer cw.Flush()

	var cols columns
	csvp := NewCSVProcessor(cr, cw)
	if o.NoHeader {
		csvp.SetPreBodyRead(func() error {
			cols = newUniqueColumns(o.ColumnSyms, nil)
			return cols.err()
		})
	} else {
		csvp.SetHeaderHanlder(func(hdr []string) ([]string, error) {
			cols = newUniqueColumns(o.ColumnSyms, hdr)
			return removeFromRecord(hdr, cols), cols.err()
		})
	}
	csvp.SetRecordHandler(func(rec []string) ([]string, error) {
		return removeFromRecord(rec, cols), nil
	})

	return csvp.Process()
}

func removeFromRecord(rec []string, cols columns) []string {
	var newRec []string
	for i, s := range rec {
		rm := false
		for _, col := range cols {
			if i == col.index {
				rm = true
			}
		}
		if !rm {
			newRec = append(newRec, s)
		}
	}
	return newRec

}
