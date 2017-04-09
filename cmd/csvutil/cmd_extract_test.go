package main

import "testing"

func Example_runExtract() {
	extractOpt.Column = "名前"
	runExtract([]string{"../../testdata/utf8.csv"})
	extractOpt.Column = ""
	// Output: 名前
	// りんご
	// みかん
}

func Test_runExtract(t *testing.T) {
	extractOpt.Column = "名前"
	if c := runExtract([]string{"../../testdata/utf8.csv"}); c != 0 {
		t.Errorf("Invalid success exit code: %d", c)
	}
	extractOpt.Column = ""
}

func Test_runExtractOnNoFile(t *testing.T) {
	if c := runExtract([]string{"../../testdata/no-file.csv"}); c == 0 {
		t.Errorf("Invalid failed exit code: %d", c)
	}
}

func Test_runExtractOnFail(t *testing.T) {
	if c := runExtract([]string{"../../testdata/broken.csv"}); c == 0 {
		t.Errorf("Invalid failed exit code: %d", c)
	}
}
