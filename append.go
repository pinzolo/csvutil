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
	// Encoding for output.
	OutputEncoding string
	// Headers is appending header list.
	Headers []string
	// Size is appending column size.
	Size int
}

func (o AppendOption) outputEncoding() string {
	if o.OutputEncoding != "" {
		return o.OutputEncoding
	}
	return o.Encoding
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
	cw := writer(w, bom, o.outputEncoding())
	defer cw.Flush()

	csvp := NewCSVProcessor(cr, cw)
	if !o.NoHeader {
		csvp.SetHeaderHanlder(func(hdr []string) ([]string, error) {
			for _, h := range o.appendingHeaders() {
				hdr = append(hdr, h)
			}
			return hdr, nil
		})
	}
	csvp.SetRecordHandler(func(rec []string) ([]string, error) {
		newRec := make([]string, len(rec)+o.Size)
		copy(newRec, rec)
		return newRec, nil
	})

	return csvp.Process()
}
