package csvutil

import (
	"bytes"
	"testing"
)

func TestCountWithBrokenCSV(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5
7,8,9
`
	r := bytes.NewBufferString(s)
	o := CountOption{}

	if _, err := Count(r, o); err == nil {
		t.Error("Count with broken csv should raise error.")
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
		t.Error(err)
	}
	if i != 3 {
		t.Errorf("Expected: %d, but got %d", 3, i)
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
		t.Error(err)
	}
	if i != 4 {
		t.Errorf("Expected: %d, but got %d", 4, i)
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
		t.Error(err)
	}
	if i != 3 {
		t.Errorf("Expected: %d, but got %d", 3, i)
	}
}
