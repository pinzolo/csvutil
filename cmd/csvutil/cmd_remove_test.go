package main

import "testing"

func Example_runRemove() {
	removeOpt.Column = "名前"
	runRemove([]string{"../../testdata/utf8.csv"})
	removeOpt.Column = ""
	// Output: 個数
	// 1
	// 2
}

func Test_runRemove(t *testing.T) {
	removeOpt.Column = "名前"
	if c := runRemove([]string{"../../testdata/utf8.csv"}); c != 0 {
		t.Errorf("Invalid success exit code: %d", c)
	}
	removeOpt.Column = ""
}

func Test_runRemoveOnNoFile(t *testing.T) {
	if c := runRemove([]string{"../../testdata/no-file.csv"}); c == 0 {
		t.Errorf("Invalid failed exit code: %d", c)
	}
}

func Test_runRemoveOnFail(t *testing.T) {
	if c := runRemove([]string{"../../testdata/broken.csv"}); c == 0 {
		t.Errorf("Invalid failed exit code: %d", c)
	}
}
