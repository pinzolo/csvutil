package main

import "testing"

func Example_runCount() {
	runCount([]string{"../../testdata/utf8.csv"})
	// Output: 2
}

func Test_runCount(t *testing.T) {
	if c := runCount([]string{"../../testdata/utf8.csv"}); c != 0 {
		t.Errorf("Invalid success exit code: %d", c)
	}
}

func Test_runCountOnNoFile(t *testing.T) {
	if c := runCount([]string{"../../testdata/no-file.csv"}); c == 0 {
		t.Errorf("Invalid failed exit code: %d", c)
	}
}

func Test_runCountOnFail(t *testing.T) {
	if c := runCount([]string{"../../testdata/broken.csv"}); c == 0 {
		t.Errorf("Invalid failed exit code: %d", c)
	}
}
