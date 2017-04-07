package csvutil

import (
	"bytes"
	"testing"
)

func TestSizeWithEmptyCSV(t *testing.T) {
	s := ""
	r := bytes.NewBufferString(s)
	o := SizeOption{}
	if _, err := Size(r, o); err == nil {
		t.Error("Size with empty CSV should raise error.")
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
