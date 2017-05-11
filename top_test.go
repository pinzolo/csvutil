package csvutil

import (
	"bytes"
	"testing"
)

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
