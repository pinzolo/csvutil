package csvutil

import (
	"io"
	"strconv"

	"github.com/pkg/errors"
)

// AppendOption is option holder for Append.
type AppendOption struct {
	// Source file does not have header line. (default false)
	NoHeader bool
	// Encoding of source file. (default utf8)
	Encoding string
	// Headers is appending header list.
	Headers []string
	// Size is appending column size.
	Size int
}

func (o AppendOption) validate() error {
	if o.Size <= 0 {
		return errors.New("negative or zero size")
	}
	return nil
}

func (o AppendOption) appendingHeaders() []string {
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

// Append empty values to end of each lines.
func Append(r io.Reader, w io.Writer, o AppendOption) error {
	if err := o.validate(); err != nil {
		return errors.Wrap(err, "invalid option")
	}

	cr, bom := reader(r, o.Encoding)
	cw := writer(w, bom, o.Encoding)
	defer cw.Flush()

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
			for _, h := range o.appendingHeaders() {
				hdr = append(hdr, h)
			}
			cw.Write(hdr)
			continue
		}
		for i := 0; i < o.Size; i++ {
			rec = append(rec, "")
		}
		cw.Write(rec)
	}

	return nil
}
