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
	// Count is appending column count.
	Count int
}

func (o AppendOption) validate() error {
	if o.Count < 0 {
		return errors.New("Count requires positive digit.")
	}
	if o.Count == 0 && len(o.Headers) == 0 {
		return errors.New("Required count or headers.")
	}
	return nil
}

func (o AppendOption) appendingCount() int {
	if len(o.Headers) > o.Count {
		return len(o.Headers)
	}
	return o.Count
}

func (o AppendOption) appendingHeaders() []string {
	if len(o.Headers) > o.Count {
		return o.Headers
	}

	hdr := make([]string, o.Count)
	copy(hdr, o.Headers)
	for i := 0; i < o.Count-len(o.Headers); i++ {
		hdr[len(o.Headers)+i] = "column" + strconv.Itoa(i+1)
	}
	return hdr
}

// Append empty values to end of each lines.
func Append(r io.Reader, w io.Writer, o AppendOption) error {
	if err := o.validate(); err != nil {
		return errors.Wrap(err, "Invalid option")
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
			return errors.Wrap(err, "Cannot read csv line")
		}
		if hdr == nil && !o.NoHeader {
			hdr = rec
			for _, h := range o.appendingHeaders() {
				hdr = append(hdr, h)
			}
			cw.Write(hdr)
			continue
		}
		for i := 0; i < o.appendingCount(); i++ {
			rec = append(rec, "")
		}
		cw.Write(rec)
	}

	return nil
}
