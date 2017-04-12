package csvutil

import (
	"bytes"
	"regexp"
	"strings"
	"testing"
)

var telNumRegex = regexp.MustCompile(`\d+-\d+-\d`)

func TestTelWithoutColumn(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := TelOption{}

	err := Tel(r, w, o)
	if err == nil {
		t.Error("Email without column symbol should raise error.")
	}
}

func TestTelWithNoHeaderButColumnNotNumber(t *testing.T) {
	s := `1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := TelOption{
		NoHeader: true,
		Column:   "foo",
	}

	err := Tel(r, w, o)
	if err == nil {
		t.Error("Email with not number column symbol for no header CSV should raise error.")
	}

}

func TestTelWithNegativeMobileRate(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := TelOption{
		Column:     "aaa",
		MobileRate: -1,
	}

	err := Tel(r, w, o)
	if err == nil {
		t.Error("Tel with negative mobile rate should raise error.")
	}
}

func TestTelWithOver100MobileRate(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := TelOption{
		Column:     "aaa",
		MobileRate: 101,
	}

	err := Tel(r, w, o)
	if err == nil {
		t.Error("Tel with over 100 rate should raise error.")
	}
}

func TestTelWithUnknownColumn(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := TelOption{
		Column: "ddd",
	}

	err := Tel(r, w, o)
	if err == nil {
		t.Error("Tel with unknown column should raise error.")
	}
}

func TestTelWithBrokenCSV(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := TelOption{
		Column: "aaa",
	}

	if err := Tel(r, w, o); err == nil {
		t.Error("Tel with broken csv should raise error.")
	}
}

func TestTelWitColumn(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := TelOption{
		Column: "aaa",
	}

	if err := Tel(r, w, o); err != nil {
		t.Error(err)
	}

	data := readCSV(w.String())
	if ok := allOK(data, 0, isTelNumber); !ok {
		t.Errorf("Tel failed updating on tel number. %+v", data)
	}
}

func TestTelMobile(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := TelOption{
		Column:     "aaa",
		MobileRate: 100,
	}

	if err := Tel(r, w, o); err != nil {
		t.Error(err)
	}

	data := readCSV(w.String())
	if ok := allOK(data, 0, isMobileTelNumber); !ok {
		t.Errorf("Tel failed updating on mobile tel number. %+v", data)
	}
}

func isTelNumber(s string) bool {
	return telNumRegex.MatchString(s)
}

func isMobileTelNumber(s string) bool {
	if !isTelNumber(s) {
		return false
	}

	for _, c := range mobileTelAreaCodes {
		if strings.HasPrefix(s, c) {
			return true
		}
	}
	return false
}

func readCSV(csv string) [][]string {
	b := bytes.NewBufferString(csv)
	r, _ := NewReader(b)
	ss, _ := r.ReadAll()
	return ss
}

func allOK(data [][]string, i int, f func(string) bool) bool {
	return allOKNoHeader(data[1:], i, f)
}

func allOKNoHeader(data [][]string, i int, f func(string) bool) bool {
	for _, rec := range data {
		for j, s := range rec {
			if j == i {
				if !f(s) {
					return false
				}
			}
		}
	}
	return true
}
