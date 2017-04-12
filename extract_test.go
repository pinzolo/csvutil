package csvutil

import (
	"bytes"
	"testing"
)

func TestExtractWithEmptyColumn(t *testing.T) {
	r := &bytes.Buffer{}
	w := &bytes.Buffer{}
	o := ExtractOption{}

	if err := Extract(r, w, o); err == nil {
		t.Error("Extract with empty column should raise error.")
	}
}

func TestExtractWithNoHeaderAndNoDigitColumn(t *testing.T) {
	r := &bytes.Buffer{}
	w := &bytes.Buffer{}
	o := ExtractOption{
		NoHeader:   true,
		ColumnSyms: []string{"foo"},
	}

	if err := Extract(r, w, o); err == nil {
		t.Error("Extract with no header and no digit column symbol should raise error.")
	}
}

func TestExtractWithSymbolColumnButNoHeader(t *testing.T) {
	s := `1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := ExtractOption{
		NoHeader:   true,
		ColumnSyms: []string{"aaa"},
	}

	if err := Extract(r, w, o); err == nil {
		t.Error("When given header text as column symbol but CSV does not have header, Extract should raise error.")
	}
}

func TestExtractWithUnknownHeaderSymbol(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := ExtractOption{
		ColumnSyms: []string{"ddd"},
	}

	if err := Extract(r, w, o); err == nil {
		t.Error("Extract with unknown header symbol should raise error.")
	}
}

func TestExtractWithBrokenCSV(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := ExtractOption{
		ColumnSyms: []string{"aaa"},
	}

	if err := Extract(r, w, o); err == nil {
		t.Error("Extract with broken csv should raise error.")
	}
}

func TestExtractWhenColumnIsHeaderText(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := ExtractOption{
		ColumnSyms: []string{"bbb"},
	}

	if err := Extract(r, w, o); err != nil {
		t.Error(err)
	}
	expected := `bbb
2
5
8
`
	if actual := w.String(); actual != expected {
		t.Errorf("Expectd: %s, but got %s", expected, actual)
	}
}

func TestExtractWhenColumnIsIndex(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := ExtractOption{
		ColumnSyms: []string{"1"},
	}

	if err := Extract(r, w, o); err != nil {
		t.Error(err)
	}
	expected := `bbb
2
5
8
`
	if actual := w.String(); actual != expected {
		t.Errorf("Expectd: %s, but got %s", expected, actual)
	}
}

func TestExtractWhenColumnIsIndexAndNoHeader(t *testing.T) {
	s := `1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := ExtractOption{
		NoHeader:   true,
		ColumnSyms: []string{"1"},
	}

	if err := Extract(r, w, o); err != nil {
		t.Error(err)
	}
	expected := `2
5
8
`
	if actual := w.String(); actual != expected {
		t.Errorf("Expectd: %s, but got %s", expected, actual)
	}
}

func TestExtractWhenColumnIsMultiColumn(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := ExtractOption{
		ColumnSyms: []string{"aaa", "bbb"},
	}

	if err := Extract(r, w, o); err != nil {
		t.Error(err)
	}
	expected := `aaa,bbb
1,2
4,5
7,8
`
	if actual := w.String(); actual != expected {
		t.Errorf("Expectd: %s, but got %s", expected, actual)
	}
}
