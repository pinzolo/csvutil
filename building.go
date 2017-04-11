package csvutil

import (
	"io"

	"github.com/pkg/errors"
)

// BuildingOption is option holder for Building.
type BuildingOption struct {
	// Source file does not have header line. (default false)
	NoHeader bool
	// Encoding of source file. (default utf8)
	Encoding string
	// Encoding for output
	OutputEncoding string
	// Target column symbol
	Column string
	// Rate of office output
	OfficeRate int
	// BlockNumber width(1 or 2)
	NumberWidth int
	// Append to source value
	Append bool
}

func (o BuildingOption) validate() error {
	if o.Column == "" {
		return errors.New("no column")
	}
	if o.NoHeader {
		if !isDigit(o.Column) {
			return errors.New("not number column symbol")
		}
	}
	if o.OfficeRate < 0 || 100 < o.OfficeRate {
		return errors.New("invalid office rate (0 <= rate <= 100)")
	}
	if o.NumberWidth != 1 && o.NumberWidth != 2 {
		return errors.New("invalid number width (1 or 2)")
	}

	return nil
}

func (o BuildingOption) isFullWidth() bool {
	return o.NumberWidth == 2
}

func (o BuildingOption) outputEncoding() string {
	if o.OutputEncoding != "" {
		return o.OutputEncoding
	}
	return o.Encoding
}

// Building overwrite value of given column by dummy office or apartment.
func Building(r io.Reader, w io.Writer, o BuildingOption) error {
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
			col, err = newColumnWithIndex(o.Column, hdr)
			if err != nil {
				return errors.Wrap(err, "column not found")
			}
		}
		newRec := make([]string, len(rec))
		for i, s := range rec {
			if i == col.index {
				if o.Append {
					newRec[i] = s
				}
				if lot(o.OfficeRate) {
					newRec[i] += fakeOffice(o.isFullWidth())
				} else {
					newRec[i] += fakeApartment(o.isFullWidth())
				}
			} else {
				newRec[i] = s
			}
		}
		cw.Write(newRec)
	}

	return nil
}
