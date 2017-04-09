package main

import "testing"

func Test_runEmail(t *testing.T) {
	emailOpt.Column = "名前"
	if c := runEmail([]string{"../../testdata/utf8.csv"}); c != 0 {
		t.Errorf("Invalid success exit code: %d", c)
	}
	emailOpt.Column = ""
}

func Test_runEmailOnNoFile(t *testing.T) {
	if c := runEmail([]string{"../../testdata/no-file.csv"}); c == 0 {
		t.Errorf("Invalid failed exit code: %d", c)
	}
}

func Test_runEmailOnFail(t *testing.T) {
	if c := runEmail([]string{"../../testdata/broken.csv"}); c == 0 {
		t.Errorf("Invalid failed exit code: %d", c)
	}
}
