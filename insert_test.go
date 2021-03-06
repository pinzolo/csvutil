package csvutil

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func BenchmarkInsert(b *testing.B) {
	p, err := ioutil.ReadFile("testdata/bench.csv")
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		r := bytes.NewBuffer(p)
		w := &bytes.Buffer{}
		o := InsertOption{
			Size:   100,
			Before: "都道府県",
		}
		Insert(r, w, o)
	}
}

func TestInsertWithoutSize(t *testing.T) {
	r := &bytes.Buffer{}
	w := &bytes.Buffer{}
	o := InsertOption{}

	if err := Insert(r, w, o); err == nil {
		t.Fatal("Insert without size should raise error.")
	}
}

func TestInsertWithNegativeSize(t *testing.T) {
	r := &bytes.Buffer{}
	w := &bytes.Buffer{}
	o := InsertOption{
		Size: -1,
	}

	if err := Insert(r, w, o); err == nil {
		t.Fatal("Insert with negative size should raise error.")
	}
}

func TestInsertWithHeadersOnly(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := InsertOption{
		Headers: []string{"foo", "bar"},
	}

	if err := Insert(r, w, o); err == nil {
		t.Fatal("Insert without size should raise error.")
	}
}

func TestInsertWithInvalidBefore(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := InsertOption{
		Size:    2,
		Before:  "ddd",
		Headers: []string{"foo", "bar"},
	}
	if err := Insert(r, w, o); err == nil {
		t.Fatal("Insert with invalid before should raise error.")
	}
}

func TestInsertWithBrokenCSV(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := InsertOption{
		Size: 1,
	}

	if err := Insert(r, w, o); err == nil {
		t.Fatal("Insert with broken csv should raise error.")
	}
}

func TestInsertWithSizeAndNoBefore(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := InsertOption{
		Size: 2,
	}
	if err := Insert(r, w, o); err != nil {
		t.Fatal(err)
	}

	expected := `column1,column2,aaa,bbb,ccc
,,1,2,3
,,4,5,6
,,7,8,9
`
	if actual := w.String(); actual != expected {
		t.Fatalf("Expectd: %s, but got %s", expected, actual)
	}
}

func TestInsertWithGreaterSizeThanHeadersLength(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := InsertOption{
		Headers: []string{"foo", "bar"},
		Size:    3,
	}
	if err := Insert(r, w, o); err != nil {
		t.Fatal(err)
	}

	expected := `foo,bar,column1,aaa,bbb,ccc
,,,1,2,3
,,,4,5,6
,,,7,8,9
`
	if actual := w.String(); actual != expected {
		t.Fatalf("Expectd: %s, but got %s", expected, actual)
	}
}

func TestInsertWithLessSizeThanHeadersLength(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := InsertOption{
		Headers: []string{"foo", "bar", "baz"},
		Size:    2,
	}
	if err := Insert(r, w, o); err != nil {
		t.Fatal(err)
	}

	expected := `foo,bar,aaa,bbb,ccc
,,1,2,3
,,4,5,6
,,7,8,9
`
	if actual := w.String(); actual != expected {
		t.Fatalf("Expectd: %s, but got %s", expected, actual)
	}
}

func TestInsertWithLessSizeThanHeadersLengthButNoHeader(t *testing.T) {
	s := `1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := InsertOption{
		Headers:  []string{"foo", "bar", "baz"},
		Size:     2,
		NoHeader: true,
	}
	if err := Insert(r, w, o); err != nil {
		t.Fatal(err)
	}

	expected := `,,1,2,3
,,4,5,6
,,7,8,9
`
	if actual := w.String(); actual != expected {
		t.Fatalf("Expectd: %s, but got %s", expected, actual)
	}
}

func TestInsertWithSizeAndBefore(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := InsertOption{
		Size:   2,
		Before: "bbb",
	}
	if err := Insert(r, w, o); err != nil {
		t.Fatal(err)
	}

	expected := `aaa,column1,column2,bbb,ccc
1,,,2,3
4,,,5,6
7,,,8,9
`
	if actual := w.String(); actual != expected {
		t.Fatalf("Expectd: %s, but got %s", expected, actual)
	}
}

func TestInsertWithSizeAndBeforeAndHeaders(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := InsertOption{
		Size:    2,
		Before:  "bbb",
		Headers: []string{"foo", "bar"},
	}
	if err := Insert(r, w, o); err != nil {
		t.Fatal(err)
	}

	expected := `aaa,foo,bar,bbb,ccc
1,,,2,3
4,,,5,6
7,,,8,9
`
	if actual := w.String(); actual != expected {
		t.Fatalf("Expectd: %s, but got %s", expected, actual)
	}
}
