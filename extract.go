package csvutil

import (
	"io"

	"github.com/pkg/errors"
)

// ExtractOption is option holder for Extract.
type ExtractOption struct {
	// Source file does not have header line. (default false)
	NoHeader bool
	// Encoding of source file. (default utf8)
	Encoding string
	// Encoding for output.
	OutputEncoding string
	// ColumnSyms header or column index list.
	ColumnSyms []string
}

func (o ExtractOption) validate() error {
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

func (o ExtractOption) outputEncoding() string {
	if o.OutputEncoding != "" {
		return o.OutputEncoding
	}
	return o.Encoding
}

// Extract column(s) from CSV.
func Extract(r io.Reader, w io.Writer, o ExtractOption) error {
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
			if err := cols.err(); err != nil {
				return nil, err
			}
			return extractFromRecord(hdr, cols), nil
		})
	}
	csvp.SetRecordHandler(func(rec []string) ([]string, error) {
		return extractFromRecord(rec, cols), nil
	})

	return csvp.Process()
}

func extractFromRecord(rec []string, cols columns) []string {
	newRec := make([]string, len(cols))
	for n, col := range cols {
		for i, s := range rec {
			if i == col.index {
				newRec[n] = s
				break
			}
		}
	}
	return newRec
}
