package csvutil

import (
	"fmt"
	"io"
	"math"
	"math/rand"
	"strconv"

	"github.com/pkg/errors"
)

// NumericOption is option holder for Numeric.
type NumericOption struct {
	// Source file does not have header line. (default false)
	NoHeader bool
	// Encoding of source file. (default utf8)
	Encoding string
	// Encoding for output.
	OutputEncoding string
	// Target column symbol.
	Column string
	// Max value
	Max int
	// Min value
	Min int
	// Output decimal instead of integer
	Decimal bool
	// Digit of decimal
	DecimalDigit int
}

func (o NumericOption) validate() error {
	if o.Column == "" {
		return errors.New("no column")
	}
	if o.NoHeader {
		if !isDigit(o.Column) {
			return errors.New("not number column symbol")

		}
	}
	if o.Max <= o.Min {
		return errors.New("max should be greater than min")
	}
	if o.Decimal && o.DecimalDigit <= 0 {
		return errors.New("decimal digit is not positive")
	}
	return nil
}

func (o NumericOption) outputEncoding() string {
	if o.OutputEncoding != "" {
		return o.OutputEncoding
	}
	return o.Encoding
}

// Numeric overwrite value of given column by random numbers.
func Numeric(r io.Reader, w io.Writer, o NumericOption) error {
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
		rec[col.index] = fakeNumeric(o)
		return rec, nil
	})

	return csvp.Process()
}

func fakeNumeric(o NumericOption) string {
	if o.Decimal {
		return fakeDecimal(o)
	}
	return fakeInteger(o)
}

func fakeDecimal(o NumericOption) string {
	coefficient := int(math.Pow10(o.DecimalDigit))
	lim := (o.Max - o.Min) * coefficient
	n := rand.Intn(lim) + (o.Min * coefficient)
	nega := false
	if n < 0 {
		nega = true
		n = n * -1
	}
	s := fmt.Sprintf("%0"+strconv.Itoa((o.DecimalDigit+1))+"d", n)
	l := len(s)
	s = s[:l-o.DecimalDigit] + "." + s[l-o.DecimalDigit:]
	if nega {
		return "-" + s
	}
	return s
}

func fakeInteger(o NumericOption) string {
	lim := o.Max - o.Min
	n := rand.Intn(lim) + o.Min
	return strconv.Itoa(n)
}
