package csvutil

import (
	"bytes"
	"testing"
)

func TestAppendWithoutHeadersAndCount(t *testing.T) {
	r := &bytes.Buffer{}
	w := &bytes.Buffer{}
	o := AppendOption{}
	err := Append(r, w, o)
	if err == nil {
		t.Error("Append without headers and count should raise error.")
	}
}

func TestAppendWithNegativeCount(t *testing.T) {
	r := &bytes.Buffer{}
	w := &bytes.Buffer{}
	o := AppendOption{
		Count: -1,
	}
	err := Append(r, w, o)
	if err == nil {
		t.Error("Append with negative count should raise error.")
	}
}

func TestAppendWithCount(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBuffer([]byte(s))
	w := &bytes.Buffer{}
	o := AppendOption{
		Count: 2,
	}
	if err := Append(r, w, o); err != nil {
		t.Error(err)
	}

	expected := `aaa,bbb,ccc,column1,column2
1,2,3,,
4,5,6,,
7,8,9,,
`
	if actual := w.String(); actual != expected {
		t.Errorf("Expectd: %s, but got %s", expected, actual)
	}
}

func TestAppendWithHeaders(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBuffer([]byte(s))
	w := &bytes.Buffer{}
	o := AppendOption{
		Headers: []string{"foo", "bar"},
	}
	if err := Append(r, w, o); err != nil {
		t.Error(err)
	}

	expected := `aaa,bbb,ccc,foo,bar
1,2,3,,
4,5,6,,
7,8,9,,
`
	if actual := w.String(); actual != expected {
		t.Errorf("Expectd: %s, but got %s", expected, actual)
	}
}

func TestAppendWithGreaterCountThanHeadersLength(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBuffer([]byte(s))
	w := &bytes.Buffer{}
	o := AppendOption{
		Headers: []string{"foo", "bar"},
		Count:   3,
	}
	if err := Append(r, w, o); err != nil {
		t.Error(err)
	}

	expected := `aaa,bbb,ccc,foo,bar,column1
1,2,3,,,
4,5,6,,,
7,8,9,,,
`
	if actual := w.String(); actual != expected {
		t.Errorf("Expectd: %s, but got %s", expected, actual)
	}
}

func TestAppendWithLessCountThanHeadersLength(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBuffer([]byte(s))
	w := &bytes.Buffer{}
	o := AppendOption{
		Headers: []string{"foo", "bar", "baz"},
		Count:   2,
	}
	if err := Append(r, w, o); err != nil {
		t.Error(err)
	}

	expected := `aaa,bbb,ccc,foo,bar,baz
1,2,3,,,
4,5,6,,,
7,8,9,,,
`
	if actual := w.String(); actual != expected {
		t.Errorf("Expectd: %s, but got %s", expected, actual)
	}
}

func TestAppendWithLessCountThanHeadersLengthButNoHeader(t *testing.T) {
	s := `1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBuffer([]byte(s))
	w := &bytes.Buffer{}
	o := AppendOption{
		Headers:  []string{"foo", "bar", "baz"},
		Count:    2,
		NoHeader: true,
	}
	if err := Append(r, w, o); err != nil {
		t.Error(err)
	}

	expected := `1,2,3,,,
4,5,6,,,
7,8,9,,,
`
	if actual := w.String(); actual != expected {
		t.Errorf("Expectd: %s, but got %s", expected, actual)
	}
}
