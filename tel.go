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
	var hdr []string
	for {
		rec, err := cr.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return errors.Wrap(err, "cannot read csv line")
		}
		if hdr == nil && !o.NoHeader {
			hdr = rec
			cw.Write(rec)
			continue
		}
		if col == nil {
			col = newColumnWithIndex(o.Column, hdr)
		}
		if col.err != nil {
			return col.err
		}
		newRec := make([]string, len(rec))
		for i, s := range rec {
			if i == col.index {
				if lot(o.MobileRate) {
					newRec[i] = fakeMobileTel()
				} else {
					newRec[i] = fakeTel()
				}
			} else {
				newRec[i] = s
			}
		}
		cw.Write(newRec)
	}

	return nil
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
