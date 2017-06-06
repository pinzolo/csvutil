package csvutil

import (
	"bytes"
	"io/ioutil"
	"regexp"
	"strings"
	"testing"
)

var telNumRegex = regexp.MustCompile(`\d+-\d+-\d`)

func BenchmarkTel(b *testing.B) {
	p, err := ioutil.ReadFile("testdata/bench.csv")
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		r := bytes.NewBuffer(p)
		w := &bytes.Buffer{}
		o := TelOption{
			Column:     "電話番号",
			MobileRate: 20,
		}
		Tel(r, w, o)
	}
}

func TestTelWithoutColumn(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := TelOption{}

	if err := Tel(r, w, o); err == nil {
		t.Fatal("Email without column symbol should raise error.")
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

	if err := Tel(r, w, o); err == nil {
		t.Fatal("Email with not number column symbol for no header CSV should raise error.")
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

	if err := Tel(r, w, o); err == nil {
		t.Fatal("Tel with negative mobile rate should raise error.")
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

	if err := Tel(r, w, o); err == nil {
		t.Fatal("Tel with over 100 rate should raise error.")
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

	if err := Tel(r, w, o); err == nil {
		t.Fatal("Tel with unknown column should raise error.")
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
		t.Fatal("Tel with broken csv should raise error.")
	}
}

func TestTelWithNoHeader(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := TelOption{
		Column:   "0",
		NoHeader: true,
	}

	if err := Tel(r, w, o); err != nil {
		t.Fatal(err)
	}

	data := readCSV(w.String())
	if ok := allOKNoHeader(data, 0, isTelNumber); !ok {
		t.Fatalf("Tel failed updating on tel number. %+v", data)
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
		t.Fatal(err)
	}

	data := readCSV(w.String())
	if ok := allOK(data, 0, isTelNumber); !ok {
		t.Fatalf("Tel failed updating on tel number. %+v", data)
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
		t.Fatal(err)
	}

	data := readCSV(w.String())
	if ok := allOK(data, 0, isMobileTelNumber); !ok {
		t.Fatalf("Tel failed updating on mobile tel number. %+v", data)
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
