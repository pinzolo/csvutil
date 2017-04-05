package csvutil

import (
	"io"
	"strconv"

	"github.com/pkg/errors"
)

// HeaderOption is option holder for Header.
type HeaderOption struct {
	// Encoding of source file. (default utf8)
	Encoding string
	// Print index
	Index bool
	// Index origin number
	IndexOrigin int
}

// Header print headers of CSV.
func Header(r io.Reader, w io.Writer, o HeaderOption) error {
	cr, bom := reader(r, o.Encoding)
	cw := writer(w, bom, o.Encoding)
	cw.Comma = '\t'
	defer cw.Flush()

	hdr, err := cr.Read()
	if err != nil {
		if err == io.EOF {
			return errors.Wrap(err, "empty csv")
		}
		return errors.Wrap(err, "cannot read csv header")
	}
	for i, h := range hdr {
		if o.Index {
			idx := strconv.Itoa(o.IndexOrigin + i)
			cw.Write([]string{idx, h})
			continue
		}
		cw.Write([]string{h})
	}

	return nil
}
