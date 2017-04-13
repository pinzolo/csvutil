package csvutil

import (
	"io"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

// SubstituteOption is option holder for Substitute.
type SubstituteOption struct {
	// Source file does not have header line. (default false)
	NoHeader bool
	// Encoding of source file. (default utf8)
	Encoding string
	// Encoding for output
	OutputEncoding string
	// Target column symbol
	Column string
	// Target pattern
	Pattern string
	// Replacement value
	Replacement string
	// Use regexp
	Regexp  bool
	regex   *regexp.Regexp
	subFunc func(string) string
}

func (o *SubstituteOption) validate() error {
	if o.Column == "" {
		return errors.New("no column")
	}
	if o.Pattern == "" {
		return errors.New("no pattern")
	}
	if o.NoHeader {
		if !isDigit(o.Column) {
			return errors.New("not number column symbol")

		}
	}
	if o.Regexp {
		r, err := regexp.Compile(o.Pattern)
		if err != nil {
			return err
		}
		o.regex = r
		o.subFunc = func(s string) string {
			return o.regex.ReplaceAllString(s, o.Replacement)
		}
	} else {
		o.subFunc = func(s string) string {
			return strings.Replace(s, o.Pattern, o.Replacement, -1)
		}
	}
	return nil
}

func (o SubstituteOption) outputEncoding() string {
	if o.OutputEncoding != "" {
		return o.OutputEncoding
	}
	return o.Encoding
}

// Substitute value of given column.
func Substitute(r io.Reader, w io.Writer, o SubstituteOption) error {
	opt := &o
	if err := opt.validate(); err != nil {
		return errors.Wrap(err, "invalid option")
	}

	cr, bom := reader(r, opt.Encoding)
	cw := writer(w, bom, opt.outputEncoding())
	defer cw.Flush()

	var col *column
	csvp := NewCSVProcessor(cr, cw)
	if o.NoHeader {
		csvp.SetPreBodyRead(func() error {
			col = newColumnWithIndex(opt.Column, nil)
			return col.err
		})
	} else {
		csvp.SetHeaderHanlder(func(hdr []string) ([]string, error) {
			col = newColumnWithIndex(opt.Column, hdr)
			return hdr, col.err
		})
	}
	csvp.SetRecordHandler(func(rec []string) ([]string, error) {
		newRec := make([]string, len(rec))
		for i, s := range rec {
			if i == col.index {
				newRec[i] = opt.subFunc(s)
				continue
			}
			newRec[i] = s
		}
		return newRec, nil
	})

	return csvp.Process()
}
