package main

import "testing"

func Example_runHeader() {
	runHeader([]string{"../../testdata/utf8.csv"})
	// Output: 名前
	// 個数
}

func Test_runHeader(t *testing.T) {
	if c := runHeader([]string{"../../testdata/utf8.csv"}); c != 0 {
		t.Errorf("Invalid success exit code: %d", c)
	}
}

func Test_runHeaderOnNoFile(t *testing.T) {
	if c := runHeader([]string{"../../testdata/no-file.csv"}); c == 0 {
		t.Errorf("Invalid failed exit code: %d", c)
	}
}

func Test_runHeaderOnFail(t *testing.T) {
	if c := runHeader([]string{"../../testdata/broken.csv"}); c == 0 {
		t.Errorf("Invalid failed exit code: %d", c)
	}
}
