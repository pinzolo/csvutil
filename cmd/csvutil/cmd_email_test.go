package main

import (
	"strings"
	"testing"
)

func Test_runEmail(t *testing.T) {
	emailOpt.Column = "名前"
	if c := runEmail([]string{testFilePath("utf8.csv")}); c != 0 {
		t.Errorf("Invalid success exit code: %d", c)
	}
	emailOpt.Column = ""
}

func Test_runEmailOnNoFile(t *testing.T) {
	if c := runEmail([]string{testFilePath("no-file.csv")}); c == 0 {
		t.Errorf("Invalid failed exit code: %d", c)
	}
}

func Test_runEmailOnFail(t *testing.T) {
	if c := runEmail([]string{testFilePath("broken.csv")}); c == 0 {
		t.Errorf("Invalid failed exit code: %d", c)
	}
}

func Test_runEmailOnBackup(t *testing.T) {
	f, err := prepareWritingTest()
	defer f()
	if err != nil {
		t.Error(err)
	}
	emailOpt.Column = "名前"
	emailOpt.Overwrite = true
	emailOpt.Backup = true
	runEmail([]string{tempFilePath()})
	emailOpt.Backup = false
	emailOpt.Overwrite = false
	emailOpt.Column = ""
	if b, err := existsBackup(); err != nil || !b {
		t.Errorf("Failed backup")
	}
}

func Test_runEmailOnOverwrite(t *testing.T) {
	f, err := prepareWritingTest()
	defer f()
	if err != nil {
		t.Error(err)
	}
	emailOpt.Column = "名前"
	emailOpt.Overwrite = true
	runEmail([]string{tempFilePath()})
	emailOpt.Overwrite = false
	emailOpt.Column = ""
	c, err := overwriteContent()
	if err != nil {
		t.Error(err)
	}
	if len(c[0]) != 2 {
		t.Errorf("Overwrite failed. got %+v", c)
	}
	if c[0][0] != "名前" || c[0][1] != "個数" {
		t.Errorf("Overwrite failed. got %+v", c)
	}
	if !strings.Contains(c[1][0], "@") || c[1][1] != "1" {
		t.Errorf("Overwrite failed. got %+v", c)
	}
	if !strings.Contains(c[2][0], "@") || c[2][1] != "2" {
		t.Errorf("Overwrite failed. got %+v", c)
	}
}
