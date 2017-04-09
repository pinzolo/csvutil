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
	NoHeader       bool
	Encoding       string
	OutputEncoding string
	Column         string
	MobileRate     int
	col            *column
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
			cw.Write(rec)
			continue
		}
		o, err = setupEmailColumn(o, hdr)
		if err != nil {
			return err
		}
		newRec := make([]string, len(rec))
		for i, s := range rec {
			if i == o.col.index {
				if lot(o.MobileRate) {
					newRec[i] = fakeMobileEmail()
				} else {
					newRec[i] = fakeEmail()
				}
			} else {
				newRec[i] = s
			}
		}
		cw.Write(newRec)
	}

	return nil
}

func setupEmailColumn(o EmailOption, hdr []string) (EmailOption, error) {
	if o.col == nil {
		col, err := newColumnWithIndex(o.Column, hdr)
		if err != nil {
			return o, errors.Wrap(err, "column not found")
		}
		o.col = col
	}
	return o, nil
}

func fakeEmail() string {
	return strings.ToLower(fake.EmailAddress())
}

func fakeMobileEmail() string {
	return strings.ToLower(fake.UserName()) + "@" + sampleString(mobileEmailDomains)
}