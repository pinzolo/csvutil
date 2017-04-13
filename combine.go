package csvutil

import (
	"io"
	"strings"

	"github.com/pkg/errors"
)

// CombineOption is option holder for Combine.
type CombineOption struct {
	// Source file does not have header line. (default false)
	NoHeader bool
	// Encoding of source file. (default utf8)
	Encoding string
	// Encoding for output.
	OutputEncoding string
	// SourceSyms header or column index list
	SourceSyms []string
	// Destination column symbol
	Destination string
	// Delimiter
	Delimiter string
}

func (o CombineOption) validate() error {
	if len(o.SourceSyms) == 0 {
		return errors.New("no column")
	}
	if o.NoHeader {
		for _, c := range o.SourceSyms {
			if !isDigit(c) {
				return errors.Errorf("not number column symbol: %s", c)
			}
		}
	}
	return nil
}

func (o CombineOption) outputEncoding() string {
	if o.OutputEncoding != "" {
		return o.OutputEncoding
	}
	return o.Encoding
}

// Combine column(s) from CSV.
func Combine(r io.Reader, w io.Writer, o CombineOption) error {
	if err := o.validate(); err != nil {
		return errors.Wrap(err, "invalid option")
	}

	cr, bom := reader(r, o.Encoding)
	cw := writer(w, bom, o.outputEncoding())
	defer cw.Flush()

	var srcs columns
	var dst *column
	csvp := NewCSVProcessor(cr, cw)
	if o.NoHeader {
		csvp.SetPreBodyRead(func() error {
			srcs = newUniqueColumns(o.SourceSyms, nil)
			dst = newColumnWithIndex(o.Destination, nil)
			if err := srcs.err(); err != nil {
				return err
			}
			return dst.err
		})
	} else {
		csvp.SetHeaderHanlder(func(hdr []string) ([]string, error) {
			srcs = newUniqueColumns(o.SourceSyms, hdr)
			dst = newColumnWithIndex(o.Destination, hdr)
			if err := srcs.err(); err != nil {
				return hdr, err
			}
			return hdr, dst.err
		})
	}
	csvp.SetRecordHandler(func(rec []string) ([]string, error) {
		newRec := make([]string, len(rec))
		vals := make([]string, len(srcs))
		for i, src := range srcs {
			vals[i] = rec[src.index]
		}
		for i, s := range rec {
			if i == dst.index {
				newRec[i] = strings.Join(vals, o.Delimiter)
			} else {
				newRec[i] = s
			}
		}
		return newRec, nil
	})

	return csvp.Process()
}
