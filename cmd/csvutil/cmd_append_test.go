package main

import "testing"

func Example_runAppend() {
	runAppend([]string{"../../testdata/utf8.csv"})
	// Output: 名前,個数,column1
	// りんご,1,
	// みかん,2,
}

func Test_runAppend(t *testing.T) {
	if c := runAppend([]string{"../../testdata/utf8.csv"}); c != 0 {
		t.Errorf("Invalid success exit code: %d", c)
	}
}

func Test_runAppendOnNoFile(t *testing.T) {
	if c := runAppend([]string{"../../testdata/no-file.csv"}); c == 0 {
		t.Errorf("Invalid failed exit code: %d", c)
	}
}

func Test_runAppendOnFail(t *testing.T) {
	if c := runAppend([]string{"../../testdata/broken.csv"}); c == 0 {
		t.Errorf("Invalid failed exit code: %d", c)
	}
}
