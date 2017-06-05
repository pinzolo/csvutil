package csvutil

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func BenchmarkRemove(b *testing.B) {
	p, err := ioutil.ReadFile("testdata/bench.csv")
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		r := bytes.NewBuffer(p)
		w := &bytes.Buffer{}
		o := RemoveOption{
			ColumnSyms: []string{"1", "3", "5", "7", "9"},
		}
		Remove(r, w, o)
	}
}

func TestRemoveWithEmptyColumn(t *testing.T) {
	r := &bytes.Buffer{}
	w := &bytes.Buffer{}
	o := RemoveOption{}

	if err := Remove(r, w, o); err == nil {
		t.Fatal("Remove with empty column should raise error.")
	}
}

func TestRemoveWithNoHeaderAndNoDigitColumn(t *testing.T) {
	r := &bytes.Buffer{}
	w := &bytes.Buffer{}
	o := RemoveOption{
		NoHeader:   true,
		ColumnSyms: []string{"foo"},
	}

	if err := Remove(r, w, o); err == nil {
		t.Fatal("Remove with no header and no digit column symbol should raise error.")
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

	if err := Remove(r, w, o); err == nil {
		t.Fatal("When given header text as column symbol but CSV does not have header, Remove should raise error.")
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

	if err := Remove(r, w, o); err == nil {
		t.Fatal("Remove with unknown header symbol should raise error.")
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
		t.Fatal("Remove with broken csv should raise error.")
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

	if err := Remove(r, w, o); err != nil {
		t.Fatal(err)
	}
	expected := `aaa,ccc
1,3
4,6
7,9
`
	if actual := w.String(); actual != expected {
		t.Fatalf("Expectd: %s, but got %s", expected, actual)
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

	if err := Remove(r, w, o); err != nil {
		t.Fatal(err)
	}
	expected := `aaa,ccc
1,3
4,6
7,9
`
	if actual := w.String(); actual != expected {
		t.Fatalf("Expectd: %s, but got %s", expected, actual)
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

	if err := Remove(r, w, o); err != nil {
		t.Fatal(err)
	}
	expected := `1,3
4,6
7,9
`
	if actual := w.String(); actual != expected {
		t.Fatalf("Expectd: %s, but got %s", expected, actual)
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

	if err := Remove(r, w, o); err != nil {
		t.Fatal(err)
	}
	expected := `ccc
3
6
9
`
	if actual := w.String(); actual != expected {
		t.Fatalf("Expectd: %s, but got %s", expected, actual)
	}
}
