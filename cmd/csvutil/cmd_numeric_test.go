package main

import (
	"regexp"
	"testing"
)

func Test_runNumeric(t *testing.T) {
	numericOpt.Column = "名前"
	if c := runNumeric([]string{testFilePath("utf8.csv")}); c != 0 {
		t.Errorf("Invalid success exit code: %d", c)
	}
	numericOpt.Column = ""
}

func Test_runNumericOnNoFile(t *testing.T) {
	if c := runNumeric([]string{testFilePath("no-file.csv")}); c == 0 {
		t.Errorf("Invalid failed exit code: %d", c)
	}
}

func Test_runNumericOnFail(t *testing.T) {
	if c := runNumeric([]string{testFilePath("broken.csv")}); c == 0 {
		t.Errorf("Invalid failed exit code: %d", c)
	}
}

func Test_runNumericOnBackup(t *testing.T) {
	f, err := prepareWritingTest()
	defer f()
	if err != nil {
		t.Error(err)
	}
	numericOpt.Column = "名前"
	numericOpt.Overwrite = true
	numericOpt.Backup = true
	runNumeric([]string{tempFilePath()})
	numericOpt.Backup = false
	numericOpt.Overwrite = false
	numericOpt.Column = ""
	if b, err := existsBackup(); err != nil || !b {
		t.Errorf("Failed backup")
	}
}

func Test_runNumericOnOverwrite(t *testing.T) {
	f, err := prepareWritingTest()
	defer f()
	if err != nil {
		t.Error(err)
	}
	numericOpt.Column = "名前"
	numericOpt.Overwrite = true
	runNumeric([]string{tempFilePath()})
	numericOpt.Overwrite = false
	numericOpt.Column = ""
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
	pat, err := regexp.Compile("^\\d+$")
	if err != nil {
		t.Fatal(err)
	}
	if !pat.MatchString(c[1][0]) || c[1][1] != "1" {
		t.Errorf("Overwrite failed. got %+v", c)
	}
	if !pat.MatchString(c[2][0]) || c[2][1] != "2" {
		t.Errorf("Overwrite failed. got %+v", c)
	}
}
