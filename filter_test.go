package csvutil

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func BenchmarkFilter(b *testing.B) {
	p, err := ioutil.ReadFile("testdata/bench.csv")
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		r := bytes.NewBuffer(p)
		w := &bytes.Buffer{}
		o := FilterOption{
			Pattern: "京都府",
		}
		Filter(r, w, o)
	}
}

func BenchmarkFilterWithRegexp(b *testing.B) {
	p, err := ioutil.ReadFile("testdata/bench.csv")
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		r := bytes.NewBuffer(p)
		w := &bytes.Buffer{}
		o := FilterOption{
			Pattern: "\\d+-\\d+",
			Regexp:  true,
		}
		Filter(r, w, o)
	}
}

func TestFilterWithNoHeaderButColumnNotNumber(t *testing.T) {
	s := `A1,B2,C3
D4,E5,F6
G7,H8,I9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := FilterOption{
		NoHeader:   true,
		ColumnSyms: []string{"aaa"},
		Pattern:    "4",
	}

	if err := Filter(r, w, o); err == nil {
		t.Error("Filter with not number column symbol for no header CSV should raise error.")
	}
}

func TestFilterWithoutPattern(t *testing.T) {
	s := `aaa,bbb,ccc
A1,B2,C3
D4,E5,F6
G7,H8,I9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := FilterOption{
		ColumnSyms: []string{"aaa"},
	}
	if err := Filter(r, w, o); err == nil {
		t.Error("Filter without pattern should raise error.")
	}
}

func TestFilterWithUnknownColumn(t *testing.T) {
	s := `aaa,bbb,ccc
A1,B2,C3
D4,E5,F6
G7,H8,I9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := FilterOption{
		ColumnSyms: []string{"ddd"},
		Pattern:    "4",
	}

	if err := Filter(r, w, o); err == nil {
		t.Error("Filter with unknown column should raise error.")
	}
}

func TestFilterWithBrokenCSV(t *testing.T) {
	s := `aaa,bbb,ccc
A1,B2,C3
D4,E5
G7,H8,I9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := FilterOption{
		ColumnSyms: []string{"aaa"},
		Pattern:    "4",
	}

	if err := Filter(r, w, o); err == nil {
		t.Error("Filter with broken csv should raise error.")
	}
}

func TestFilterWithBrokenRegexp(t *testing.T) {
	s := `aaa,bbb,ccc
A1,B2,C3
D4,E5
G7,H8,I9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := FilterOption{
		Pattern: "[A-E",
		Regexp:  true,
	}

	if err := Filter(r, w, o); err == nil {
		t.Error("Filter with broken regexp should raise error.")
	}
}

func TestFilterWithoutCoumnSyms(t *testing.T) {
	s := `aaa,bbb,ccc
A1,B2,C3
D4,E5,F6
G7,H8,I4
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := FilterOption{
		Pattern: "4",
	}

	if err := Filter(r, w, o); err != nil {
		t.Fatal(err)
	}

	expected := `aaa,bbb,ccc
D4,E5,F6
G7,H8,I4
`
	if actual := w.String(); actual != expected {
		t.Errorf("Expectd: %s, but got %s", expected, actual)
	}
}

func TestFilterWithNoHeader(t *testing.T) {
	s := `A1,B2,C3
D4,E5,F6
G7,H8,I4
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := FilterOption{
		Pattern:  "4",
		NoHeader: true,
	}

	if err := Filter(r, w, o); err != nil {
		t.Fatal(err)
	}

	expected := `D4,E5,F6
G7,H8,I4
`
	if actual := w.String(); actual != expected {
		t.Errorf("Expectd: %s, but got %s", expected, actual)
	}
}

func TestFilterWithCoumnSyms(t *testing.T) {
	s := `aaa,bbb,ccc
A1,B4,C3
D4,E5,F6
G7,H8,I4
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := FilterOption{
		Pattern:    "4",
		ColumnSyms: []string{"aaa", "bbb"},
	}

	if err := Filter(r, w, o); err != nil {
		t.Fatal(err)
	}

	expected := `aaa,bbb,ccc
A1,B4,C3
D4,E5,F6
`
	if actual := w.String(); actual != expected {
		t.Errorf("Expectd: %s, but got %s", expected, actual)
	}
}

func TestFilterWithRegexp(t *testing.T) {
	s := `aaa,bbb,ccc
A1,B2,C3
D4,E5,F6
G7,H8,I4
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := FilterOption{
		Pattern: "[A-E]",
		Regexp:  true,
	}

	if err := Filter(r, w, o); err != nil {
		t.Fatal(err)
	}

	expected := `aaa,bbb,ccc
A1,B2,C3
D4,E5,F6
`
	if actual := w.String(); actual != expected {
		t.Errorf("Expectd: %s, but got %s", expected, actual)
	}

}
