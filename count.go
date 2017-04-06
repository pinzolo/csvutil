package csvutil

import (
	"io"

	"github.com/pkg/errors"
)

// CountOption is option holder for Count.
type CountOption struct {
	// Source file does not have header line. (default false)
	NoHeader bool
	// Encoding of source file. (default utf8)
	Encoding string
}

// Count CSV lines.
func Count(r io.Reader, o CountOption) (int, error) {
	cr, _ := reader(r, o.Encoding)

	i := 0
	var hdr bool
	for {
		_, err := cr.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return 0, errors.Wrap(err, "cannot read csv line")
		}
		if !hdr && !o.NoHeader {
			hdr = true
			continue
		}
		i++
	}

	return i, nil
}
