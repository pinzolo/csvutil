package csvutil

import (
	"io"

	"github.com/pkg/errors"
)

// TopOption is option holder for Top.
type TopOption struct {
	// Source file does not have header line. (default false)
	NoHeader bool
	// Encoding of source file. (default utf8)
	Encoding string
	// Encoding for output.
	OutputEncoding string
	// Headers is appending header list.
	Headers []string
	// Count is reading line count.
	Count int
}

func (o TopOption) validate() error {
	if o.Count <= 0 {
		return errors.New("negative or zero count")
	}
	return nil
}

func (o TopOption) outputEncoding() string {
	if o.OutputEncoding != "" {
		return o.OutputEncoding
	}
	return o.Encoding
}

// Top reads lines from top.
func Top(r io.Reader, w io.Writer, o TopOption) error {
	if err := o.validate(); err != nil {
		return errors.Wrap(err, "invalid option")
	}

	cr, bom := reader(r, o.Encoding)
	cw := writer(w, bom, o.outputEncoding())
	defer cw.Flush()

	if !o.NoHeader {
		hdr, err := cr.Read()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		cw.Write(hdr)
	}

	for i := 0; i < o.Count; i++ {
		rec, err := cr.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		cw.Write(rec)
	}
	return nil
}
