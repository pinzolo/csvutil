package main

import "testing"

func Example_runBlank() {
	blankOpt.Column = "0"
	runBlank([]string{"../../testdata/utf8.csv"})
	blankOpt.Column = ""
	// Output: 名前,個数
	// ,1
	// ,2
}

func Test_runBlank(t *testing.T) {
	blankOpt.Column = "0"
	if c := runBlank([]string{"../../testdata/utf8.csv"}); c != 0 {
		t.Errorf("Invalid success exit code: %d", c)
	}
	blankOpt.Column = ""
}

func Test_runBlankOnNoFile(t *testing.T) {
	if c := runBlank([]string{"../../testdata/no-file.csv"}); c == 0 {
		t.Errorf("Invalid failed exit code: %d", c)
	}
}

func Test_runBlankOnFail(t *testing.T) {
	if c := runBlank([]string{"../../testdata/broken.csv"}); c == 0 {
		t.Errorf("Invalid failed exit code: %d", c)
	}
}
