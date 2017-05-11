package csvutil

import (
	"io"

	"github.com/pkg/errors"
)

// TailOption is option holder for Tail.
type TailOption struct {
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

func (o TailOption) validate() error {
	if o.Count <= 0 {
		return errors.New("negative or zero count")
	}
	return nil
}

func (o TailOption) outputEncoding() string {
	if o.OutputEncoding != "" {
		return o.OutputEncoding
	}
	return o.Encoding
}

// Tail reads lines from tail.
func Tail(r io.Reader, w io.Writer, o TailOption) error {
	if err := o.validate(); err != nil {
		return errors.Wrap(err, "invalid option")
	}

	cr, bom := reader(r, o.Encoding)
	cw := writer(w, bom, o.outputEncoding())

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

	var recs [][]string
	for {
		rec, err := cr.Read()
		if err != nil {
			if err == io.EOF {
				if len(recs) != 0 {
					cw.WriteAll(recs)
				}
				break
			}
			return err
		}
		if len(recs) < o.Count {
			recs = append(recs, rec)
			continue
		}
		recs = append(recs[1:], rec)
	}
	return nil
}
