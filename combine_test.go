package csvutil

import (
	"bytes"
	"testing"
)

func TestCombineWithEmptySource(t *testing.T) {
	s := `1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := CombineOption{
		Destination: "ccc",
	}

	if err := Combine(r, w, o); err == nil {
		t.Error("Combine with empty column should raise error.")
	}
}

func TestCombineWithNoHeaderAndNoDigitSource(t *testing.T) {
	r := &bytes.Buffer{}
	w := &bytes.Buffer{}
	o := CombineOption{
		NoHeader:    true,
		SourceSyms:  []string{"aaa", "bbb"},
		Destination: "ccc",
	}

	if err := Combine(r, w, o); err == nil {
		t.Error("Combine with no header and no digit column symbol should raise error.")
	}
}

func TestCombineWithSymbolSourceButNoHeader(t *testing.T) {
	s := `1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := CombineOption{
		NoHeader:    true,
		SourceSyms:  []string{"aaa", "bbb"},
		Destination: "2",
	}

	if err := Combine(r, w, o); err == nil {
		t.Error("When given header text as column symbol but CSV does not have header, Combine should raise error.")
	}
}

func TestCombineWithSymbolDestinationButNoHeader(t *testing.T) {
	s := `1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := CombineOption{
		NoHeader:    true,
		SourceSyms:  []string{"0", "1"},
		Destination: "ccc",
	}

	if err := Combine(r, w, o); err == nil {
		t.Error("When given header text as column symbol but CSV does not have header, Combine should raise error.")
	}
}

func TestCombineWithUnknownSourceSymbol(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := CombineOption{
		SourceSyms:  []string{"ddd"},
		Destination: "ccc",
	}

	if err := Combine(r, w, o); err == nil {
		t.Error("Combine with unknown source symbol should raise error.")
	}
}

func TestCombineWithUnknownDestinationSymbol(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := CombineOption{
		SourceSyms:  []string{"aaa", "bbb"},
		Destination: "ddd",
	}

	if err := Combine(r, w, o); err == nil {
		t.Error("Combine with unknown destination symbol should raise error.")
	}
}

func TestCombineWithBrokenCSV(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := CombineOption{
		SourceSyms:  []string{"aaa", "bbb"},
		Destination: "ddd",
	}

	if err := Combine(r, w, o); err == nil {
		t.Error("Combine with broken csv should raise error.")
	}
}

func TestCombine(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := CombineOption{
		SourceSyms:  []string{"aaa", "bbb"},
		Destination: "ccc",
	}

	if err := Combine(r, w, o); err != nil {
		t.Error(err)
	}
	expected := `aaa,bbb,ccc
1,2,12
4,5,45
7,8,78
`
	if actual := w.String(); actual != expected {
		t.Errorf("Expectd: %s, but got %s", expected, actual)
	}
}

func TestCombineWhenSourceIsIndex(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := CombineOption{
		SourceSyms:  []string{"0", "bbb"},
		Destination: "ccc",
	}

	if err := Combine(r, w, o); err != nil {
		t.Error(err)
	}
	expected := `aaa,bbb,ccc
1,2,12
4,5,45
7,8,78
`
	if actual := w.String(); actual != expected {
		t.Errorf("Expectd: %s, but got %s", expected, actual)
	}
}

func TestCombineWhenDestinationIsIndex(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := CombineOption{
		SourceSyms:  []string{"0", "bbb"},
		Destination: "2",
	}

	if err := Combine(r, w, o); err != nil {
		t.Error(err)
	}
	expected := `aaa,bbb,ccc
1,2,12
4,5,45
7,8,78
`
	if actual := w.String(); actual != expected {
		t.Errorf("Expectd: %s, but got %s", expected, actual)
	}
}

func TestCombineWhenSourceIsIndexAndNoHeader(t *testing.T) {
	s := `1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := CombineOption{
		NoHeader:    true,
		SourceSyms:  []string{"0", "1"},
		Destination: "2",
	}

	if err := Combine(r, w, o); err != nil {
		t.Error(err)
	}
	expected := `1,2,12
4,5,45
7,8,78
`
	if actual := w.String(); actual != expected {
		t.Errorf("Expectd: %s, but got %s", expected, actual)
	}
}

func TestCombineWhenSourceContainsDestination(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := CombineOption{
		SourceSyms:  []string{"0", "bbb"},
		Destination: "aaa",
	}

	if err := Combine(r, w, o); err != nil {
		t.Error(err)
	}
	expected := `aaa,bbb,ccc
12,2,3
45,5,6
78,8,9
`
	if actual := w.String(); actual != expected {
		t.Errorf("Expectd: %s, but got %s", expected, actual)
	}
}

func TestCombineWithDelimiter(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := CombineOption{
		SourceSyms:  []string{"0", "bbb"},
		Destination: "ccc",
		Delimiter:   "_",
	}

	if err := Combine(r, w, o); err != nil {
		t.Error(err)
	}
	expected := `aaa,bbb,ccc
1,2,1_2
4,5,4_5
7,8,7_8
`
	if actual := w.String(); actual != expected {
		t.Errorf("Expectd: %s, but got %s", expected, actual)
	}
}

func TestCombineWithSpecialCharDelimiter(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := CombineOption{
		SourceSyms:  []string{"bbb", "0"},
		Destination: "ccc",
		Delimiter:   ",",
	}

	if err := Combine(r, w, o); err != nil {
		t.Error(err)
	}
	expected := `aaa,bbb,ccc
1,2,"2,1"
4,5,"5,4"
7,8,"8,7"
`
	if actual := w.String(); actual != expected {
		t.Errorf("Expectd: %s, but got %s", expected, actual)
	}
}
