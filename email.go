package csvutil

import (
	"io"
	"strings"

	"github.com/icrowley/fake"
	"github.com/pkg/errors"
)

var mobileEmailDomains = []string{
	"docomo.ne.jp",
	"ezweb.ne.jp",
	"softbank.ne.jp",
	"i.softbank.ne.jp",
	"ymobile.ne.jp",
	"emobile.ne.jp",
}

// EmailOption is option holder for Email.
type EmailOption struct {
	// Source file does not have header line. (default false)
	NoHeader bool
	// Encoding of source file. (default utf8)
	Encoding string
	// Encoding for output.
	OutputEncoding string
	// Target column symbol.
	Column string
	// Rate of output mobile email address.
	MobileRate int
}

func (o EmailOption) validate() error {
	if o.Column == "" {
		return errors.New("no column")
	}
	if o.NoHeader {
		if !isDigit(o.Column) {
			return errors.New("not number column symbol")

		}
	}
	if o.MobileRate < 0 || 100 < o.MobileRate {
		return errors.New("invalid mobile rate (0 <= rate <= 100)")
	}
	return nil
}

func (o EmailOption) outputEncoding() string {
	if o.OutputEncoding != "" {
		return o.OutputEncoding
	}
	return o.Encoding
}

// Email overwrite value of given column by dummy email address.
func Email(r io.Reader, w io.Writer, o EmailOption) error {
	if err := o.validate(); err != nil {
		return errors.Wrap(err, "invalid option")
	}

	cr, bom := reader(r, o.Encoding)
	cw := writer(w, bom, o.outputEncoding())
	defer cw.Flush()

	var col *column
	csvp := NewCSVProcessor(cr, cw)
	if o.NoHeader {
		csvp.SetPreBodyRead(func() error {
			col = newColumnWithIndex(o.Column, nil)
			return col.err
		})
	} else {
		csvp.SetHeaderHanlder(func(hdr []string) ([]string, error) {
			col = newColumnWithIndex(o.Column, hdr)
			return hdr, col.err
		})
	}
	csvp.SetRecordHandler(func(rec []string) ([]string, error) {
		if lot(o.MobileRate) {
			rec[col.index] = fakeMobileEmail()
		} else {
			rec[col.index] = fakeEmail()
		}
		return rec, nil
	})

	return csvp.Process()
}

func fakeEmail() string {
	return strings.ToLower(fake.EmailAddress())
}

func fakeMobileEmail() string {
	return strings.ToLower(fake.UserName()) + "@" + sampleString(mobileEmailDomains)
}
