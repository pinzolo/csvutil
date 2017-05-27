package csvutil

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func BenchmarkTail(b *testing.B) {
	p, err := ioutil.ReadFile("testdata/bench.csv")
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		r := bytes.NewBuffer(p)
		w := &bytes.Buffer{}
		o := TailOption{
			Count: 100,
		}
		Tail(r, w, o)
	}
}

func TestTailWithZeroCount(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := TailOption{}

	if err := Tail(r, w, o); err == nil {
		t.Error("Tail with zero count should raise error.")
	}
}

func TestTailWithNegativeCount(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := TailOption{
		Count: -1,
	}

	if err := Tail(r, w, o); err == nil {
		t.Error("Tail with negative count should raise error.")
	}
}

func TestTailWithBrokenCSV(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := TailOption{
		Count: 2,
	}

	if err := Tail(r, w, o); err == nil {
		t.Error("Tail with broken csv should raise error.")
	}
}

func TestTailWithBrokenheaderCSV(t *testing.T) {
	s := `a"aa",bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := TailOption{
		Count: 2,
	}

	if err := Tail(r, w, o); err == nil {
		t.Error("Tail with broken header csv should raise error.")
	}
}

func TestTailWithEmptyCSV(t *testing.T) {
	r := bytes.NewBufferString("")
	w := &bytes.Buffer{}
	o := TailOption{
		Count: 2,
	}

	expected := ""

	if err := Tail(r, w, o); err != nil {
		t.Error(err)
	}

	if actual := w.String(); actual != expected {
		t.Errorf("Expectd: %s, but got %s", expected, actual)
	}
}

func TestTailWithLessCountThanLineCount(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := TailOption{
		Count: 2,
	}

	expected := `aaa,bbb,ccc
4,5,6
7,8,9
`
	if err := Tail(r, w, o); err != nil {
		t.Error(err)
	}

	if actual := w.String(); actual != expected {
		t.Errorf("Expectd: %s, but got %s", expected, actual)
	}
}

func TestTailWithGreaterCountThanLineCount(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := TailOption{
		Count: 5,
	}

	expected := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	if err := Tail(r, w, o); err != nil {
		t.Error(err)
	}

	if actual := w.String(); actual != expected {
		t.Errorf("Expectd: %s, but got %s", expected, actual)
	}
}

func TestTailWithNoHeader(t *testing.T) {
	s := `1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := TailOption{
		NoHeader: true,
		Count:    2,
	}

	expected := `4,5,6
7,8,9
`
	if err := Tail(r, w, o); err != nil {
		t.Error(err)
	}

	if actual := w.String(); actual != expected {
		t.Errorf("Expectd: %s, but got %s", expected, actual)
	}
}
