package csvutil

import (
	"io"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

// FilterOption is option holder for Filter.
type FilterOption struct {
	// Source file does not have header line. (default false)
	NoHeader bool
	// Encoding of source file. (default utf8)
	Encoding string
	// Encoding for output
	OutputEncoding string
	// ColumnSyms header or column index list.
	ColumnSyms []string
	// Target pattern
	Pattern string
	// Use regexp
	Regexp  bool
	regex   *regexp.Regexp
	matches func(string) bool
}

func (o *FilterOption) validate() error {
	if o.NoHeader && len(o.ColumnSyms) != 0 {
		for _, c := range o.ColumnSyms {
			if !isDigit(c) {
				return errors.New("not number column symbol")
			}
		}
	}
	if o.Pattern == "" {
		return errors.New("no pattern")
	}
	if o.Regexp {
		r, err := regexp.Compile(o.Pattern)
		if err != nil {
			return err
		}
		o.regex = r
		o.matches = func(s string) bool {
			return o.regex.MatchString(s)
		}
	} else {
		o.matches = func(s string) bool {
			return strings.Contains(s, o.Pattern)
		}
	}
	return nil
}

func (o FilterOption) outputEncoding() string {
	if o.OutputEncoding != "" {
		return o.OutputEncoding
	}
	return o.Encoding
}

// Filter value of given column.
func Filter(r io.Reader, w io.Writer, o FilterOption) error {
	opt := &o
	if err := opt.validate(); err != nil {
		return errors.Wrap(err, "invalid option")
	}

	cr, bom := reader(r, opt.Encoding)
	cw := writer(w, bom, opt.outputEncoding())
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
			return hdr, cols.err()
		})
	}
	csvp.SetRecordHandler(func(rec []string) ([]string, error) {
		for i, s := range rec {
			if len(cols) == 0 && o.matches(s) {
				return rec, nil
			}
			for _, col := range cols {
				if i == col.index && o.matches(s) {
					return rec, nil
				}
			}
		}
		return nil, nil
	})

	return csvp.Process()
}
