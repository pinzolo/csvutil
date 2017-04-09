package csvutil

import (
	"bytes"
	"strings"
	"testing"
)

func TestEmailWithoutColumn(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := EmailOption{}

	err := Email(r, w, o)
	if err == nil {
		t.Error("Email without column symbol should raise error.")
	}
}

func TestEmailWithNoHeaderButColumnNotNumber(t *testing.T) {
	s := `1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := EmailOption{
		NoHeader: true,
		Column:   "foo",
	}

	err := Email(r, w, o)
	if err == nil {
		t.Error("Email with not number column symbol for no header CSV should raise error.")
	}

}

func TestEmailWithNegativeMobileRate(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := EmailOption{
		Column:     "aaa",
		MobileRate: -1,
	}

	err := Email(r, w, o)
	if err == nil {
		t.Error("Email with negative mobile rate should raise error.")
	}
}

func TestEmailWithOver100MobileRate(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := EmailOption{
		Column:     "aaa",
		MobileRate: 101,
	}

	err := Email(r, w, o)
	if err == nil {
		t.Error("Email with over 100 mobile rate should raise error.")
	}
}

func TestEmailWithUnknownColumn(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := EmailOption{
		Column: "ddd",
	}

	err := Email(r, w, o)
	if err == nil {
		t.Error("Email with unknown column should raise error.")
	}
}

func TestEmailWitColumn(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := EmailOption{
		Column: "aaa",
	}

	if err := Email(r, w, o); err != nil {
		t.Error(err)
	}

	data := readCSV(w.String())
	if ok := allOK(data, 0, isEmail); !ok {
		t.Errorf("Email failed updating on email address. %+v", data)
	}
}

func TestEmailMobile(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := EmailOption{
		Column:     "aaa",
		MobileRate: 100,
	}

	if err := Email(r, w, o); err != nil {
		t.Error(err)
	}

	data := readCSV(w.String())
	if ok := allOK(data, 0, isMobileEmailNumber); !ok {
		t.Errorf("Email failed updating on mobile email address. %+v", data)
	}
}

func isEmail(s string) bool {
	return strings.Contains(s, "@")
}

func isMobileEmailNumber(s string) bool {
	if !isEmail(s) {
		return false
	}

	for _, c := range mobileEmailDomains {
		if strings.HasSuffix(s, c) {
			return true
		}
	}
	return false
}
