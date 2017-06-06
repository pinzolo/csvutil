package csvutil

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func BenchmarkBlank(b *testing.B) {
	p, err := ioutil.ReadFile("testdata/bench.csv")
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		r := bytes.NewBuffer(p)
		w := &bytes.Buffer{}
		o := BlankOption{
			Rate:       100,
			ColumnSyms: []string{"1", "3", "5", "7", "9"},
		}
		Blank(r, w, o)
	}
}

func TestBlankWithEmptyColumn(t *testing.T) {
	r := &bytes.Buffer{}
	w := &bytes.Buffer{}
	o := BlankOption{
		SpaceSize: 1,
		Rate:      100,
	}

	if err := Blank(r, w, o); err == nil {
		t.Fatal("Blank with empty column should raise error.")
	}
}

func TestBlankWithNoHeaderAndNoDigitColumn(t *testing.T) {
	r := &bytes.Buffer{}
	w := &bytes.Buffer{}
	o := BlankOption{
		NoHeader:   true,
		ColumnSyms: []string{"foo"},
		SpaceSize:  1,
		Rate:       100,
	}

	if err := Blank(r, w, o); err == nil {
		t.Fatal("Blank with no header and no digit column symbol should raise error.")
	}
}

func TestBlankWithNegativeSpaceWidth(t *testing.T) {
	r := &bytes.Buffer{}
	w := &bytes.Buffer{}
	o := BlankOption{
		ColumnSyms: []string{"1"},
		SpaceWidth: -1,
		SpaceSize:  1,
		Rate:       100,
	}

	if err := Blank(r, w, o); err == nil {
		t.Fatal("Blank with negative width should raise error.")
	}
}

func TestBlankWithOver2SpaceWidth(t *testing.T) {
	r := &bytes.Buffer{}
	w := &bytes.Buffer{}
	o := BlankOption{
		ColumnSyms: []string{"1"},
		SpaceWidth: 3,
		SpaceSize:  1,
		Rate:       100,
	}

	if err := Blank(r, w, o); err == nil {
		t.Fatal("Blank with over 2 width should raise error.")
	}
}

func TestBlankWithNegativeSpaceSize(t *testing.T) {
	r := &bytes.Buffer{}
	w := &bytes.Buffer{}
	o := BlankOption{
		ColumnSyms: []string{"1"},
		SpaceSize:  -1,
		Rate:       100,
	}

	if err := Blank(r, w, o); err == nil {
		t.Fatal("Blank with negative size should raise error.")
	}
}

func TestBlankWithZeroSpaceSize(t *testing.T) {
	r := &bytes.Buffer{}
	w := &bytes.Buffer{}
	o := BlankOption{
		ColumnSyms: []string{"1"},
		SpaceSize:  0,
		Rate:       100,
	}

	if err := Blank(r, w, o); err == nil {
		t.Fatal("Blank with zero size should raise error.")
	}
}

func TestBlankWithNegativeRate(t *testing.T) {
	r := &bytes.Buffer{}
	w := &bytes.Buffer{}
	o := BlankOption{
		ColumnSyms: []string{"1"},
		SpaceSize:  1,
		Rate:       -1,
	}

	if err := Blank(r, w, o); err == nil {
		t.Fatal("Blank with negative rate should raise error.")
	}
}

func TestBlankWithOver101Rate(t *testing.T) {
	r := &bytes.Buffer{}
	w := &bytes.Buffer{}
	o := BlankOption{
		ColumnSyms: []string{"1"},
		SpaceSize:  1,
		Rate:       101,
	}

	if err := Blank(r, w, o); err == nil {
		t.Fatal("Blank with over 101 rate should raise error.")
	}
}

func TestBlankWithSymbolColumnButNoHeader(t *testing.T) {
	s := `1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := BlankOption{
		NoHeader:   true,
		ColumnSyms: []string{"aaa"},
		SpaceSize:  1,
		Rate:       100,
	}

	if err := Blank(r, w, o); err == nil {
		t.Fatal("When given header text as column symbol but CSV does not have header, Blank should raise error.")
	}
}

func TestBlankWithUnknownHeaderSymbol(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := BlankOption{
		ColumnSyms: []string{"ddd"},
		SpaceSize:  1,
		Rate:       100,
	}

	if err := Blank(r, w, o); err == nil {
		t.Fatal("Blank with unknown header symbol should raise error.")
	}
}

func TestBlankWithBrokenCSV(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := BlankOption{
		ColumnSyms: []string{"ddd"},
		SpaceSize:  1,
		Rate:       100,
	}

	if err := Blank(r, w, o); err == nil {
		t.Fatal("Blank with broken csv should raise error.")
	}
}

func TestBlankWhenColumnIsHeaderText(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := BlankOption{
		ColumnSyms: []string{"bbb"},
		SpaceSize:  1,
		Rate:       100,
	}

	if err := Blank(r, w, o); err != nil {
		t.Fatal(err)
	}
	expected := `aaa,bbb,ccc
1,,3
4,,6
7,,9
`
	if actual := w.String(); actual != expected {
		t.Fatalf("Expectd: %s, but got %s", expected, actual)
	}
}

func TestBlankWhenColumnIsIndex(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := BlankOption{
		ColumnSyms: []string{"1"},
		SpaceSize:  1,
		Rate:       100,
	}

	if err := Blank(r, w, o); err != nil {
		t.Fatal(err)
	}
	expected := `aaa,bbb,ccc
1,,3
4,,6
7,,9
`
	if actual := w.String(); actual != expected {
		t.Fatalf("Expectd: %s, but got %s", expected, actual)
	}
}

func TestBlankWhenColumnIsIndexAndNoHeader(t *testing.T) {
	s := `1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := BlankOption{
		NoHeader:   true,
		ColumnSyms: []string{"1"},
		SpaceSize:  1,
		Rate:       100,
	}

	if err := Blank(r, w, o); err != nil {
		t.Fatal(err)
	}
	expected := `1,,3
4,,6
7,,9
`
	if actual := w.String(); actual != expected {
		t.Fatalf("Expectd: %s, but got %s", expected, actual)
	}
}

func TestBlankWhenSpaceWidthIs1AndSpaceSizeIs1(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := BlankOption{
		ColumnSyms: []string{"bbb"},
		SpaceWidth: 1,
		SpaceSize:  1,
		Rate:       100,
	}

	if err := Blank(r, w, o); err != nil {
		t.Fatal(err)
	}
	expected := `aaa,bbb,ccc
1," ",3
4," ",6
7," ",9
`
	if actual := w.String(); actual != expected {
		t.Fatalf("Expectd: %s, but got %s", expected, actual)
	}
}

func TestBlankWhenSpaceWidthIs2AndSpaceSizeIs3(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := BlankOption{
		ColumnSyms: []string{"bbb"},
		SpaceWidth: 2,
		SpaceSize:  3,
		Rate:       100,
	}

	if err := Blank(r, w, o); err != nil {
		t.Fatal(err)
	}
	expected := `aaa,bbb,ccc
1,"　　　",3
4,"　　　",6
7,"　　　",9
`
	if actual := w.String(); actual != expected {
		t.Fatalf("Expectd: %s, but got %s", expected, actual)
	}

}

func TestBlankWhenColumnIsMultiColumn(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := BlankOption{
		ColumnSyms: []string{"aaa", "bbb"},
		SpaceSize:  1,
		Rate:       100,
	}

	if err := Blank(r, w, o); err != nil {
		t.Fatal(err)
	}
	expected := `aaa,bbb,ccc
,,3
,,6
,,9
`
	if actual := w.String(); actual != expected {
		t.Fatalf("Expectd: %s, but got %s", expected, actual)
	}
}
