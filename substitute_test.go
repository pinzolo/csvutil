package csvutil

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func BenchmarkSubstitute(b *testing.B) {
	p, err := ioutil.ReadFile("testdata/bench.csv")
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		r := bytes.NewBuffer(p)
		w := &bytes.Buffer{}
		o := SubstituteOption{
			Column:      "住所",
			Pattern:     "-",
			Replacement: "@",
		}
		Substitute(r, w, o)
	}
}

func BenchmarkSubstituteWithRegexp(b *testing.B) {
	p, err := ioutil.ReadFile("testdata/bench.csv")
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		r := bytes.NewBuffer(p)
		w := &bytes.Buffer{}
		o := SubstituteOption{
			Column:      "住所",
			Pattern:     "\\d+-\\d+",
			Replacement: "@",
			Regexp:      true,
		}
		Substitute(r, w, o)
	}
}

func TestSubstituteWithoutColumn(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := SubstituteOption{
		Pattern: "4",
	}

	if err := Substitute(r, w, o); err == nil {
		t.Error("Substitute without column symbol should raise error.")
	}
}

func TestSubstituteWithNoHeaderButColumnNotNumber(t *testing.T) {
	s := `1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := SubstituteOption{
		NoHeader: true,
		Column:   "aaa",
		Pattern:  "4",
	}

	if err := Substitute(r, w, o); err == nil {
		t.Error("Substitute with not number column symbol for no header CSV should raise error.")
	}
}

func TestSubstituteWithoutPattern(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := SubstituteOption{
		Column: "aaa",
	}

	if err := Substitute(r, w, o); err == nil {
		t.Error("Substitute without pattern should raise error.")
	}
}

func TestSubstituteWithUnknownColumn(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := SubstituteOption{
		Column:  "ddd",
		Pattern: "4",
	}

	if err := Substitute(r, w, o); err == nil {
		t.Error("Substitute with unknown column should raise error.")
	}
}

func TestSubstituteWithBrokenCSV(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := SubstituteOption{
		Column:  "aaa",
		Pattern: "4",
	}

	if err := Substitute(r, w, o); err == nil {
		t.Error("Substitute with broken csv should raise error.")
	}
}

func TestSubstituteWithBrokenRegexp(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := SubstituteOption{
		Column:  "aaa",
		Pattern: "[1-4",
		Regexp:  true,
	}

	if err := Substitute(r, w, o); err == nil {
		t.Error("Substitute with broken regexp should raise error.")
	}
}

func TestSubstituteWithNoHeader(t *testing.T) {
	s := `x4x,2,3
4y4,5,6
z8z,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := SubstituteOption{
		Column:   "0",
		Pattern:  "4",
		NoHeader: true,
	}

	if err := Substitute(r, w, o); err != nil {
		t.Error(err)
	}

	data := readCSV(w.String())
	expecteds := []string{"xx", "y", "z8z"}
	for i, expected := range expecteds {
		if actual := data[i][0]; actual != expected {
			t.Errorf("Expected %s, but got %s", expected, actual)
		}
	}
}

func TestSubstituteWithNoReplacement(t *testing.T) {
	s := `aaa,bbb,ccc
x4x,2,3
4y4,5,6
z8z,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := SubstituteOption{
		Column:  "aaa",
		Pattern: "4",
	}

	if err := Substitute(r, w, o); err != nil {
		t.Error(err)
	}

	data := readCSV(w.String())
	expecteds := []string{"aaa", "xx", "y", "z8z"}
	for i, expected := range expecteds {
		if actual := data[i][0]; actual != expected {
			t.Errorf("Expected %s, but got %s", expected, actual)
		}
	}
}

func TestSubstituteWithNoReplacementUsingRegexp(t *testing.T) {
	s := `aaa,bbb,ccc
x4x,2,3
4y4,5,6
z8z,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := SubstituteOption{
		Column:  "aaa",
		Pattern: `\d`,
		Regexp:  true,
	}

	if err := Substitute(r, w, o); err != nil {
		t.Error(err)
	}

	data := readCSV(w.String())
	expecteds := []string{"aaa", "xx", "y", "zz"}
	for i, expected := range expecteds {
		if actual := data[i][0]; actual != expected {
			t.Errorf("Expected %s, but got %s", expected, actual)
		}
	}
}

func TestSubstituteWithReplacement(t *testing.T) {
	s := `aaa,bbb,ccc
x4x,2,3
4y4,5,6
z8z,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := SubstituteOption{
		Column:      "aaa",
		Pattern:     "[0-9]",
		Replacement: "FOO",
		Regexp:      true,
	}

	if err := Substitute(r, w, o); err != nil {
		t.Error(err)
	}

	data := readCSV(w.String())
	expecteds := []string{"aaa", "xFOOx", "FOOyFOO", "zFOOz"}
	for i, expected := range expecteds {
		if actual := data[i][0]; actual != expected {
			t.Errorf("Expected %s, but got %s", expected, actual)
		}
	}
}
