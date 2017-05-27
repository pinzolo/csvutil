package csvutil

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func BenchmarkTop(b *testing.B) {
	p, err := ioutil.ReadFile("testdata/bench.csv")
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		r := bytes.NewBuffer(p)
		w := &bytes.Buffer{}
		o := TopOption{
			Count: 100,
		}
		Top(r, w, o)
	}
}

func TestTopWithZeroCount(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := TopOption{}

	if err := Top(r, w, o); err == nil {
		t.Error("Top with zero count should raise error.")
	}
}

func TestTopWithNegativeCount(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := TopOption{
		Count: -1,
	}

	if err := Top(r, w, o); err == nil {
		t.Error("Top with negative count should raise error.")
	}
}

func TestTopWithBrokenCSV(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := TopOption{
		Count: 2,
	}

	if err := Top(r, w, o); err == nil {
		t.Error("Top with broken csv should raise error.")
	}
}

func TestTopWithBrokenheaderCSV(t *testing.T) {
	s := `a"aa",bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := TopOption{
		Count: 2,
	}

	if err := Top(r, w, o); err == nil {
		t.Error("Top with broken header csv should raise error.")
	}
}

func TestTopWithEmptyCSV(t *testing.T) {
	r := bytes.NewBufferString("")
	w := &bytes.Buffer{}
	o := TopOption{
		Count: 2,
	}

	expected := ""

	if err := Top(r, w, o); err != nil {
		t.Error(err)
	}

	if actual := w.String(); actual != expected {
		t.Errorf("Expectd: %s, but got %s", expected, actual)
	}
}

func TestTopWithLessCountThanLineCount(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := TopOption{
		Count: 2,
	}

	expected := `aaa,bbb,ccc
1,2,3
4,5,6
`
	if err := Top(r, w, o); err != nil {
		t.Error(err)
	}

	if actual := w.String(); actual != expected {
		t.Errorf("Expectd: %s, but got %s", expected, actual)
	}
}

func TestTopWithGreaterCountThanLineCount(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := TopOption{
		Count: 5,
	}

	expected := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	if err := Top(r, w, o); err != nil {
		t.Error(err)
	}

	if actual := w.String(); actual != expected {
		t.Errorf("Expectd: %s, but got %s", expected, actual)
	}
}

func TestTopWithNoHeader(t *testing.T) {
	s := `1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := TopOption{
		NoHeader: true,
		Count:    2,
	}

	expected := `1,2,3
4,5,6
`
	if err := Top(r, w, o); err != nil {
		t.Error(err)
	}

	if actual := w.String(); actual != expected {
		t.Errorf("Expectd: %s, but got %s", expected, actual)
	}
}
