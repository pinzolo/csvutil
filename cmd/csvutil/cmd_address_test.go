package main

import (
	"regexp"
	"testing"
)

func Test_runAddress(t *testing.T) {
	addressOpt.ZipCode = "0"
	if c := runAddress([]string{testFilePath("utf8.csv")}); c != 0 {
		t.Errorf("Invalid success exit code: %d", c)
	}
	addressOpt.ZipCode = ""
}

func Test_runAddressOnNoFile(t *testing.T) {
	if c := runAddress([]string{testFilePath("no-file.csv")}); c == 0 {
		t.Errorf("Invalid failed exit code: %d", c)
	}
}

func Test_runAddressOnFail(t *testing.T) {
	if c := runAddress([]string{testFilePath("broken.csv")}); c == 0 {
		t.Errorf("Invalid failed exit code: %d", c)
	}
}

func Test_runAddressOnBackup(t *testing.T) {
	f, err := prepareWritingTest()
	defer f()
	if err != nil {
		t.Error(err)
	}
	addressOpt.ZipCode = "0"
	addressOpt.Overwrite = true
	addressOpt.Backup = true
	runAddress([]string{tempFilePath()})
	addressOpt.Backup = false
	addressOpt.Overwrite = false
	addressOpt.ZipCode = ""
	if b, err := existsBackup(); err != nil || !b {
		t.Errorf("Failed backup")
	}
}

func Test_runAddressOnOverwrite(t *testing.T) {
	f, err := prepareWritingTest()
	defer f()
	if err != nil {
		t.Error(err)
	}
	addressOpt.ZipCode = "0"
	addressOpt.Overwrite = true
	runAddress([]string{tempFilePath()})
	addressOpt.Overwrite = false
	addressOpt.ZipCode = ""
	c, err := overwriteContent()
	if err != nil {
		t.Error(err)
	}
	r := regexp.MustCompile(`\d{3}-\d{4}`)
	if len(c[0]) != 2 {
		t.Errorf("Overwrite failed. got %+v", c)
	}
	if c[0][0] != "名前" || c[0][1] != "個数" {
		t.Errorf("Overwrite failed. got %+v", c)
	}
	if !r.MatchString(c[1][0]) || c[1][1] != "1" {
		t.Errorf("Overwrite failed. got %+v", c)
	}
	if !r.MatchString(c[2][0]) || c[2][1] != "2" {
		t.Errorf("Overwrite failed. got %+v", c)
	}
}
