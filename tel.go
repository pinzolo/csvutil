package csvutil

import (
	"fmt"
	"io"
	"math/rand"

	"github.com/pkg/errors"
)

var mobileTelAreaCodes = []string{
	"090",
	"080",
	"070",
	"050",
}

// TelOption is option holder for Tel.
type TelOption struct {
	// Source file does not have header line. (default false)
	NoHeader bool
	// Encoding of source file. (default utf8)
	Encoding string
	// Encoding for output.
	OutputEncoding string
	// Target column symbol.
	Column string
	// Rate of output mobile tel number.
	MobileRate int
}

func (o TelOption) validate() error {
	if o.Column == "" {
		return errors.New("no column")
	}
	if o.NoHeader {
		if !isDigit(o.Column) {
			return errors.New("not number column symbol")

		}
	}
	if o.MobileRate < 0 || 100 < o.MobileRate {
		return errors.New("invalid mobile rate (0 <= rate <= 100)")
	}
	return nil
}

func (o TelOption) outputEncoding() string {
	if o.OutputEncoding != "" {
		return o.OutputEncoding
	}
	return o.Encoding
}

// Tel overwrite value of given column by dummy tel number.
func Tel(r io.Reader, w io.Writer, o TelOption) error {
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
		if lot(o.MobileRate) {
			rec[col.index] = fakeMobileTel()
		} else {
			rec[col.index] = fakeTel()
		}
		return rec, nil
	})

	return csvp.Process()
}

func fakeTel() string {
	var ac string
	for ac == "" || containsString(mobileTelAreaCodes, ac) {
		ac = fmt.Sprintf("0%d", rand.Intn(99)+1)
	}
	return fmt.Sprintf("%s-%04d-%04d", ac, rand.Intn(10000), rand.Intn(10000))
}

func fakeMobileTel() string {
	return fmt.Sprintf("%s-%04d-%04d", sampleString(mobileTelAreaCodes), rand.Intn(10000), rand.Intn(10000))
}
