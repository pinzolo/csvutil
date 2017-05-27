package csvutil

import (
	"bytes"
	"testing"
)

func BenchmarkGenerate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := &bytes.Buffer{}
		o := GenerateOption{
			Size:  100,
			Count: 1000,
		}
		Generate(w, o)
	}
}

func TestGenerateWithoutSize(t *testing.T) {
	w := &bytes.Buffer{}
	o := GenerateOption{
		Count: 3,
	}
	if err := Generate(w, o); err == nil {
		t.Error("Generate without size should raise error.")
	}
}

func TestGenerateWithNegativeSize(t *testing.T) {
	w := &bytes.Buffer{}
	o := GenerateOption{
		Count: 3,
		Size:  -1,
	}
	if err := Generate(w, o); err == nil {
		t.Error("Generate with negative size should raise error.")
	}
}

func TestGenerateWithoutCount(t *testing.T) {
	w := &bytes.Buffer{}
	o := GenerateOption{
		Size: 3,
	}
	if err := Generate(w, o); err == nil {
		t.Error("Generate without size should raise error.")
	}
}

func TestGenerateWithNegativeCount(t *testing.T) {
	w := &bytes.Buffer{}
	o := GenerateOption{
		Count: -1,
		Size:  3,
	}
	if err := Generate(w, o); err == nil {
		t.Error("Generate with negative size should raise error.")
	}
}

func TestGenerate(t *testing.T) {
	w := &bytes.Buffer{}
	o := GenerateOption{
		Count: 3,
		Size:  3,
	}
	expected := `column1,column2,column3
,,
,,
,,
`
	if err := Generate(w, o); err != nil {
		t.Error(err)
	}

	if actual := w.String(); actual != expected {
		t.Errorf("Expectd: %s, but got %s", expected, actual)
	}
}

func TestGenerateWithHeaders(t *testing.T) {
	w := &bytes.Buffer{}
	o := GenerateOption{
		Count:   3,
		Size:    3,
		Headers: []string{"foo", "bar", "baz"},
	}
	expected := `foo,bar,baz
,,
,,
,,
`
	if err := Generate(w, o); err != nil {
		t.Error(err)
	}

	if actual := w.String(); actual != expected {
		t.Errorf("Expectd: %s, but got %s", expected, actual)
	}
}

func TestGenerateWithNoHeader(t *testing.T) {
	w := &bytes.Buffer{}
	o := GenerateOption{
		Count:    3,
		Size:     3,
		NoHeader: true,
	}
	expected := `,,
,,
,,
`
	if err := Generate(w, o); err != nil {
		t.Error(err)
	}

	if actual := w.String(); actual != expected {
		t.Errorf("Expectd: %s, but got %s", expected, actual)
	}
}

func TestGenerateWithGreaterSizeThanHeadersLength(t *testing.T) {
	w := &bytes.Buffer{}
	o := GenerateOption{
		Count:   3,
		Size:    3,
		Headers: []string{"foo", "bar"},
	}
	expected := `foo,bar,column1
,,
,,
,,
`
	if err := Generate(w, o); err != nil {
		t.Error(err)
	}

	if actual := w.String(); actual != expected {
		t.Errorf("Expectd: %s, but got %s", expected, actual)
	}
}

func TestGenerateWithLessSizeThanHeadersLength(t *testing.T) {
	w := &bytes.Buffer{}
	o := GenerateOption{
		Count:   3,
		Size:    3,
		Headers: []string{"foo", "bar", "baz", "xxx"},
	}
	expected := `foo,bar,baz
,,
,,
,,
`
	if err := Generate(w, o); err != nil {
		t.Error(err)
	}

	if actual := w.String(); actual != expected {
		t.Errorf("Expectd: %s, but got %s", expected, actual)
	}
}

func TestGenerateWithNoHeaderButGivenHeaders(t *testing.T) {
	w := &bytes.Buffer{}
	o := GenerateOption{
		Count:    3,
		Size:     3,
		NoHeader: true,
		Headers:  []string{"foo", "bar", "baz"},
	}
	expected := `,,
,,
,,
`
	if err := Generate(w, o); err != nil {
		t.Error(err)
	}

	if actual := w.String(); actual != expected {
		t.Errorf("Expectd: %s, but got %s", expected, actual)
	}
}
