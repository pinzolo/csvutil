package csvutil

import (
	"io"
	"strings"

	"github.com/pkg/errors"
)

// BlankOption is option holder for Blank.
type BlankOption struct {
	// Source file does not have header line. (default false)
	NoHeader bool
	// Encoding of source file. (default utf8)
	Encoding string
	// Encoding for output.
	OutputEncoding string
	// ColumnSyms header or column index list.
	ColumnSyms []string
	// Rate of fill
	Rate int
	// Space character width.
	//   0: no space(empaty)
	//   1: half space
	//   2: full width space
	SpaceWidth int
	// Space character count.
	SpaceSize int
}

func (o BlankOption) validate() error {
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
	if o.SpaceWidth < 0 || 2 < o.SpaceWidth {
		return errors.New("invalid space width")
	}
	if o.SpaceSize < 0 {
		return errors.New("invalid space size")
	}
	if o.Rate < 0 || 100 < o.Rate {
		return errors.New("invalid rate")
	}
	return nil
}

func (o BlankOption) outputEncoding() string {
	if o.OutputEncoding != "" {
		return o.OutputEncoding
	}
	return o.Encoding
}

// Blank overwrite value of given column by empty or spaces.
func Blank(r io.Reader, w io.Writer, o BlankOption) error {
	if err := o.validate(); err != nil {
		return errors.Wrap(err, "invalid option")
	}

	cr, bom := reader(r, o.Encoding)
	cw := writer(w, bom, o.outputEncoding())
	defer cw.Flush()

	var cols []*column
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
		if cols == nil {
			cols, err = newUniqueColumns(o.ColumnSyms, hdr)
			if err != nil {
				return errors.Wrap(err, "column not found")
			}
		}
		newRec := make([]string, len(rec))
		for i, s := range rec {
			newRec[i] = s
			for _, col := range cols {
				if i == col.index && lot(o.Rate) {
					newRec[i] = getSpace(o)
				}
			}
		}
		cw.Write(newRec)
	}

	return nil
}

func getSpace(o BlankOption) string {
	sp := ""
	if o.SpaceWidth == 1 {
		sp = " "
	} else if o.SpaceWidth == 2 {
		sp = "　"
	}
	if o.SpaceSize == 0 {
		sp = ""
	} else {
		sp = strings.Repeat(sp, o.SpaceSize)
	}
	return sp
}
