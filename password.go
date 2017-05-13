package csvutil

import (
	"io"

	"github.com/icrowley/fake"
	"github.com/pkg/errors"
)

// PasswordOption is option holder for Password.
type PasswordOption struct {
	// Source file does not have header line. (default false)
	NoHeader bool
	// Encoding of source file. (default utf8)
	Encoding string
	// Encoding for output.
	OutputEncoding string
	// Target column symbol.
	Column string
	// MinLength of password
	MinLength int
	// MaxLength of password
	MaxLength int
	// NoNumeric not using number flag
	NoNumeric bool
	// NoUpeer not using upper alphabets flag
	NoUpper bool
	// NoSpecial not using marks flag
	NoSpecial bool
}

func (o PasswordOption) validate() error {
	if o.Column == "" {
		return errors.New("no column")
	}
	if o.NoHeader {
		if !isDigit(o.Column) {
			return errors.New("not number column symbol")
		}
	}
	if o.MinLength <= 0 {
		return errors.New("min length not positive")
	}
	if o.MaxLength <= 0 {
		return errors.New("max length not positive")
	}
	if o.MinLength > o.MaxLength {
		return errors.New("max length less than min length")
	}
	return nil
}

func (o PasswordOption) outputEncoding() string {
	if o.OutputEncoding != "" {
		return o.OutputEncoding
	}
	return o.Encoding
}

// Password overwrite value of given column by dummy password address.
func Password(r io.Reader, w io.Writer, o PasswordOption) error {
	if err := o.validate(); err != nil {
		return errors.Wrap(err, "invalid option")
	}

	cr, bom := reader(r, o.Encoding)
	cw := writer(w, bom, o.outputEncoding())
	defer cw.Flush()

	var col *column
	csvp := NewCSVProcessor(cr, cw)
	if o.NoHeader {
		csvp.SetPreBodyRead(func() error {
			col = newColumnWithIndex(o.Column, nil)
			return col.err
		})
	} else {
		csvp.SetHeaderHanlder(func(hdr []string) ([]string, error) {
			col = newColumnWithIndex(o.Column, hdr)
			return hdr, col.err
		})
	}
	csvp.SetRecordHandler(func(rec []string) ([]string, error) {
		newRec := make([]string, len(rec))
		for i, s := range rec {
			if i == col.index {
				newRec[i] = fakePassword(o)
			} else {
				newRec[i] = s
			}
		}
		return newRec, nil
	})

	return csvp.Process()
}

func fakePassword(o PasswordOption) string {
	return fake.Password(o.MinLength, o.MaxLength, !o.NoUpper, !o.NoNumeric, !o.NoSpecial)
}
