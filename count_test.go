package csvutil

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func BenchmarkCount(b *testing.B) {
	p, err := ioutil.ReadFile("testdata/bench.csv")
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		r := bytes.NewBuffer(p)
		o := CountOption{}
		Count(r, o)
	}
}

func TestCountWithBrokenCSV(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5
7,8,9
`
	r := bytes.NewBufferString(s)
	o := CountOption{}

	if _, err := Count(r, o); err == nil {
		t.Fatal("Count with broken csv should raise error.")
	}
}

func TestCountWithHeader(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`

	r := bytes.NewBufferString(s)
	o := CountOption{}

	i, err := Count(r, o)
	if err != nil {
		t.Fatal(err)
	}
	if i != 3 {
		t.Fatalf("Expected: %d, but got %d", 3, i)
	}
}

func TestCountWithNoHeader(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`

	r := bytes.NewBufferString(s)
	o := CountOption{
		NoHeader: true,
	}

	i, err := Count(r, o)
	if err != nil {
		t.Fatal(err)
	}
	if i != 4 {
		t.Fatalf("Expected: %d, but got %d", 4, i)
	}
}

func TestCountWhenValueHaveNewLines(t *testing.T) {
	s := `aaa,bbb,"a
b
c"
1,2,3
4,5,"1
2
3"
7,8,9
`

	r := bytes.NewBufferString(s)
	o := CountOption{}

	i, err := Count(r, o)
	if err != nil {
		t.Fatal(err)
	}
	if i != 3 {
		t.Fatalf("Expected: %d, but got %d", 3, i)
	}
}
