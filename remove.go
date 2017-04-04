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

// Remove column(s) in CSV.
func Remove(r io.Reader, w io.Writer, o RemoveOption) error {
	if err := o.validate(); err != nil {
		return errors.Wrap(err, "invalid option")
	}

	cr, bom := reader(r, o.Encoding)
	cw := writer(w, bom, o.Encoding)
	defer cw.Flush()

	var cols []*column
	var hdr []string
	for {
		rec, err := cr.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return errors.Wrap(err, "cannot read line")
		}
		if hdr == nil && !o.NoHeader {
			hdr = rec
		}
		if cols == nil {
			cols, err = newUniqueColumns(o.ColumnSyms, hdr)
			if err != nil {
				return errors.Wrap(err, "cannot find index")
			}
		}
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
		cw.Write(newRec)
	}

	return nil
}
