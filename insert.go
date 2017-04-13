package csvutil

import (
	"io"
	"strconv"

	"github.com/pkg/errors"
)

// InsertOption is option holder for Insert.
type InsertOption struct {
	// Source file does not have header line. (default false)
	NoHeader bool
	// Encoding of source file. (default utf8)
	Encoding string
	// Encoding for output.
	OutputEncoding string
	// Headers is header list for insert.
	Headers []string
	// Size is appending column size.
	Size int
	// Before is insert start column symbol.
	Before string
}

func (o InsertOption) before() string {
	if o.Before == "" {
		return "0"
	}
	return o.Before
}

func (o InsertOption) outputEncoding() string {
	if o.OutputEncoding != "" {
		return o.OutputEncoding
	}
	return o.Encoding
}

func (o InsertOption) validate() error {
	if o.Size <= 0 {
		return errors.New("negative or zero size")
	}
	return nil
}

func (o InsertOption) headers() []string {
	hdr := make([]string, o.Size)
	hl := len(o.Headers)
	for i := 0; i < o.Size; i++ {
		if hl > i {
			hdr[i] = o.Headers[i]
			continue
		}
		hdr[i] = "column" + strconv.Itoa(i-hl+1)
	}
	return hdr
}

// Insert empty values to CSV
func Insert(r io.Reader, w io.Writer, o InsertOption) error {
	if err := o.validate(); err != nil {
		return errors.Wrap(err, "invalid option")
	}

	cr, bom := reader(r, o.Encoding)
	cw := writer(w, bom, o.outputEncoding())
	defer cw.Flush()

	var col *column
	vals := make([]string, o.Size)
	csvp := NewCSVProcessor(cr, cw)
	if o.NoHeader {
		csvp.SetPreBodyRead(func() error {
			col = newColumnWithIndex(o.before(), nil)
			return col.err
		})
	} else {
		csvp.SetHeaderHanlder(func(hdr []string) ([]string, error) {
			col = newColumnWithIndex(o.before(), hdr)
			if col.err != nil {
				return nil, col.err
			}
			return insertTo(hdr, col, o.Size, o.headers()), nil
		})
	}
	csvp.SetRecordHandler(func(rec []string) ([]string, error) {
		return insertTo(rec, col, o.Size, vals), nil
	})

	return csvp.Process()
}

func insertTo(rec []string, col *column, size int, ss []string) []string {
	hdr := make([]string, len(rec)+size)
	for i, h := range rec {
		if i < col.index {
			hdr[i] = h
		} else if i == col.index {
			for j, s := range ss {
				hdr[i+j] = s
			}
			hdr[col.index+size] = h
		} else {
			hdr[i+size] = h
		}
	}

	return hdr
}
