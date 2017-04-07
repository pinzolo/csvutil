package csvutil

import (
	"io"

	"github.com/pkg/errors"
)

// SizeOption is option holder for Size.
type SizeOption struct {
	// Encoding of source file. (default utf8)
	Encoding string
}

// Size CSV lines.
func Size(r io.Reader, o SizeOption) (int, error) {
	cr, _ := reader(r, o.Encoding)

	rec, err := cr.Read()
	if err != nil {
		if err == io.EOF {
			return 0, nil
		}
		return 0, errors.Wrap(err, "cannot read csv line")
	}

	return len(rec), nil
}
