package csvutil

import (
	"bytes"
	"testing"
)

func TestAppendWithoutSize(t *testing.T) {
	r := &bytes.Buffer{}
	w := &bytes.Buffer{}
	o := AppendOption{}
	if err := Append(r, w, o); err == nil {
		t.Error("Append without size should raise error.")
	}
}

func TestAppendWithNegativeSize(t *testing.T) {
	r := &bytes.Buffer{}
	w := &bytes.Buffer{}
	o := AppendOption{
		Size: -1,
	}
	if err := Append(r, w, o); err == nil {
		t.Error("Append with negative size should raise error.")
	}
}

func TestAppendWithHeadersOnly(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := AppendOption{
		Headers: []string{"foo", "bar"},
	}
	if err := Append(r, w, o); err == nil {
		t.Error("Append without size should raise error.")
	}
}

func TestAppendWithBrokenCSV(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := AppendOption{
		Size: 1,
	}

	err := Append(r, w, o)
	if err == nil {
		t.Error("Append with broken csv should raise error.")
	}
}

func TestAppendWithSize(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := AppendOption{
		Size: 2,
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

func TestAppendWithGreaterSizeThanHeadersLength(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := AppendOption{
		Headers: []string{"foo", "bar"},
		Size:    3,
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

func TestAppendWithLessSizeThanHeadersLength(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := AppendOption{
		Headers: []string{"foo", "bar", "baz"},
		Size:    2,
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

func TestAppendWithLessSizeThanHeadersLengthButNoHeader(t *testing.T) {
	s := `1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := AppendOption{
		Headers:  []string{"foo", "bar", "baz"},
		Size:     2,
		NoHeader: true,
	}
	if err := Append(r, w, o); err != nil {
		t.Error(err)
	}

	expected := `1,2,3,,
4,5,6,,
7,8,9,,
`
	if actual := w.String(); actual != expected {
		t.Errorf("Expectd: %s, but got %s", expected, actual)
	}
}
