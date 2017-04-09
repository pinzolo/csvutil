package main

import "testing"

func Test_runTel(t *testing.T) {
	telOpt.Column = "名前"
	if c := runTel([]string{"../../testdata/utf8.csv"}); c != 0 {
		t.Errorf("Invalid success exit code: %d", c)
	}
	telOpt.Column = ""
}

func Test_runTelOnNoFile(t *testing.T) {
	if c := runTel([]string{"../../testdata/no-file.csv"}); c == 0 {
		t.Errorf("Invalid failed exit code: %d", c)
	}
}

func Test_runTelOnFail(t *testing.T) {
	if c := runTel([]string{"../../testdata/broken.csv"}); c == 0 {
		t.Errorf("Invalid failed exit code: %d", c)
	}
}
