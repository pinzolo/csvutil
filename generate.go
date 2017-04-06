package csvutil

import (
	"io"
	"strconv"

	"github.com/pkg/errors"
)

// GenerateOption is option holder for Generate.
type GenerateOption struct {
	// Source file does not have header line. (default false)
	NoHeader bool
	// Encoding for output.
	OutputEncoding string
	// Headers is appending header list.
	Headers []string
	// Size is generating column size.
	Size int
	// Count is generating line count.
	Count int
}

func (o GenerateOption) validate() error {
	if o.Size <= 0 {
		return errors.New("negative or zero size")
	}
	if o.Count <= 0 {
		return errors.New("negative or zero count")
	}
	return nil
}

func (o GenerateOption) actualHeaders() []string {
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

// Generate empty values CSV.
func Generate(w io.Writer, o GenerateOption) error {
	if err := o.validate(); err != nil {
		return errors.Wrap(err, "invalid option")
	}

	cw := writer(w, false, o.OutputEncoding)
	defer cw.Flush()
	if !o.NoHeader {
		cw.Write(o.actualHeaders())
	}
	for i := 0; i < o.Count; i++ {
		rec := make([]string, o.Size)
		cw.Write(rec)
	}

	return nil
}
