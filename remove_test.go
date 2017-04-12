package csvutil

import (
	"bytes"
	"testing"
)

func TestRemoveWithEmptyColumn(t *testing.T) {
	r := &bytes.Buffer{}
	w := &bytes.Buffer{}
	o := RemoveOption{}

	err := Remove(r, w, o)
	if err == nil {
		t.Error("Remove with empty column should raise error.")
	}
}

func TestRemoveWithNoHeaderAndNoDigitColumn(t *testing.T) {
	r := &bytes.Buffer{}
	w := &bytes.Buffer{}
	o := RemoveOption{
		NoHeader:   true,
		ColumnSyms: []string{"foo"},
	}

	err := Remove(r, w, o)
	if err == nil {
		t.Error("Remove with no header and no digit column symbol should raise error.")
	}
}

func TestRemoveWithSymbolColumnButNoHeader(t *testing.T) {
	s := `1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := RemoveOption{
		NoHeader:   true,
		ColumnSyms: []string{"aaa"},
	}
	err := Remove(r, w, o)
	if err == nil {
		t.Error("When given header text as column symbol but CSV does not have header, Remove should raise error.")
	}
}

func TestRemoveWithUnknownHeaderSymbol(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := RemoveOption{
		ColumnSyms: []string{"ddd"},
	}
	err := Remove(r, w, o)
	if err == nil {
		t.Error("Remove with unknown header symbol should raise error.")
	}
}

func TestRemoveWithBrokenCSV(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := RemoveOption{
		ColumnSyms: []string{"aaa"},
	}

	if err := Remove(r, w, o); err == nil {
		t.Error("Remove with broken csv should raise error.")
	}
}

func TestRemoveWhenColumnIsHeaderText(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := RemoveOption{
		ColumnSyms: []string{"bbb"},
	}
	err := Remove(r, w, o)
	if err != nil {
		t.Error(err)
	}
	expected := `aaa,ccc
1,3
4,6
7,9
`
	if actual := w.String(); actual != expected {
		t.Errorf("Expectd: %s, but got %s", expected, actual)
	}
}

func TestRemoveWhenColumnIsIndex(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := RemoveOption{
		ColumnSyms: []string{"1"},
	}
	err := Remove(r, w, o)
	if err != nil {
		t.Error(err)
	}
	expected := `aaa,ccc
1,3
4,6
7,9
`
	if actual := w.String(); actual != expected {
		t.Errorf("Expectd: %s, but got %s", expected, actual)
	}
}

func TestRemoveWhenColumnIsIndexAndNoHeader(t *testing.T) {
	s := `1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := RemoveOption{
		NoHeader:   true,
		ColumnSyms: []string{"1"},
	}
	err := Remove(r, w, o)
	if err != nil {
		t.Error(err)
	}
	expected := `1,3
4,6
7,9
`
	if actual := w.String(); actual != expected {
		t.Errorf("Expectd: %s, but got %s", expected, actual)
	}
}

func TestRemoveWhenColumnIsMultiColumn(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := RemoveOption{
		ColumnSyms: []string{"aaa", "bbb"},
	}
	err := Remove(r, w, o)
	if err != nil {
		t.Error(err)
	}
	expected := `ccc
3
6
9
`
	if actual := w.String(); actual != expected {
		t.Errorf("Expectd: %s, but got %s", expected, actual)
	}
}
