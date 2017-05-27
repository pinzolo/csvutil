package csvutil

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func BenchmarkSize(b *testing.B) {
	p, err := ioutil.ReadFile("testdata/bench.csv")
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		r := bytes.NewBuffer(p)
		o := SizeOption{}
		Size(r, o)
	}
}

func TestSizeWithEmptyCSV(t *testing.T) {
	s := ""
	r := bytes.NewBufferString(s)
	o := SizeOption{}
	actual, err := Size(r, o)
	if err != nil {
		t.Error(err)
	}

	expected := 0
	if actual != expected {
		t.Errorf("Expectd: %d, but got %d", expected, actual)
	}
}

func TestSizeWithBrokenCSV(t *testing.T) {
	s := `a"aa,bbb,ccc
1,2,3
4,5
7,8,9
`
	r := bytes.NewBufferString(s)
	o := SizeOption{}

	if _, err := Size(r, o); err == nil {
		t.Error("Size with broken csv should raise error.")
	}
}

func TestSize(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	o := SizeOption{}
	actual, err := Size(r, o)
	if err != nil {
		t.Error(err)
	}

	expected := 3
	if actual != expected {
		t.Errorf("Expectd: %d, but got %d", expected, actual)
	}
}
