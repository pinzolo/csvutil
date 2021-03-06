package csvutil

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func BenchmarkHeader(b *testing.B) {
	p, err := ioutil.ReadFile("testdata/bench.csv")
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		r := bytes.NewBuffer(p)
		w := &bytes.Buffer{}
		o := HeaderOption{
			Index: true,
		}
		Header(r, w, o)
	}
}

func TestHeaderWithBrokenCSV(t *testing.T) {
	s := `a"aa,bbb,ccc
1,2,3
4,5
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := HeaderOption{}

	if err := Header(r, w, o); err == nil {
		t.Fatal("Header with broken csv should raise error.")
	}
}

func TestHeaderWithEmpty(t *testing.T) {
	s := ``
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := HeaderOption{}
	if err := Header(r, w, o); err != nil {
		t.Fatal(err)
	}

	expected := ""
	if actual := w.String(); actual != expected {
		t.Fatalf("Expectd: %s, but got %s", expected, actual)
	}
}

func TestHeader(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := HeaderOption{}
	if err := Header(r, w, o); err != nil {
		t.Fatal(err)
	}

	expected := `aaa
bbb
ccc
`
	if actual := w.String(); actual != expected {
		t.Fatalf("Expectd: %s, but got %s", expected, actual)
	}
}

func TestHeaderWithIndex(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := HeaderOption{
		Index: true,
	}
	if err := Header(r, w, o); err != nil {
		t.Fatal(err)
	}

	expected := `0	aaa
1	bbb
2	ccc
`
	if actual := w.String(); actual != expected {
		t.Fatalf("Expectd: %s, but got %s", expected, actual)
	}
}

func TestHeaderWithIndexOrigin(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := HeaderOption{
		Index:       true,
		IndexOrigin: 1,
	}
	if err := Header(r, w, o); err != nil {
		t.Fatal(err)
	}

	expected := `1	aaa
2	bbb
3	ccc
`
	if actual := w.String(); actual != expected {
		t.Fatalf("Expectd: %s, but got %s", expected, actual)
	}
}

func TestHeaderWithIndexOriginButWithoutIndex(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := HeaderOption{
		Index:       false,
		IndexOrigin: 1,
	}
	if err := Header(r, w, o); err != nil {
		t.Fatal(err)
	}

	expected := `aaa
bbb
ccc
`
	if actual := w.String(); actual != expected {
		t.Fatalf("Expectd: %s, but got %s", expected, actual)
	}
}
