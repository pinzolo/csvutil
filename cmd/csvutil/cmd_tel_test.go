package main

import (
	"regexp"
	"testing"
)

func Test_runTel(t *testing.T) {
	telOpt.Column = "名前"
	if c := runTel([]string{testFilePath("utf8.csv")}); c != 0 {
		t.Fatalf("Invalid success exit code: %d", c)
	}
	telOpt.Column = ""
}

func Test_runTelOnNoFile(t *testing.T) {
	if c := runTel([]string{testFilePath("no-file.csv")}); c == 0 {
		t.Fatalf("Invalid failed exit code: %d", c)
	}
}

func Test_runTelOnFail(t *testing.T) {
	if c := runTel([]string{testFilePath("broken.csv")}); c == 0 {
		t.Fatalf("Invalid failed exit code: %d", c)
	}
}

func Test_runTelOnBackup(t *testing.T) {
	f, err := prepareWritingTest()
	defer f()
	if err != nil {
		t.Fatal(err)
	}
	telOpt.Column = "名前"
	telOpt.Overwrite = true
	telOpt.Backup = true
	runTel([]string{tempFilePath()})
	telOpt.Backup = false
	telOpt.Overwrite = false
	telOpt.Column = ""
	if b, err := existsBackup(); err != nil || !b {
		t.Fatalf("Failed backup")
	}
}

func Test_runTelOnOverwrite(t *testing.T) {
	f, err := prepareWritingTest()
	defer f()
	if err != nil {
		t.Fatal(err)
	}
	telOpt.Column = "名前"
	telOpt.Overwrite = true
	runTel([]string{tempFilePath()})
	telOpt.Overwrite = false
	telOpt.Column = ""
	c, err := overwriteContent()
	if err != nil {
		t.Fatal(err)
	}
	r := regexp.MustCompile(`\d+-\d+-\d`)
	if len(c[0]) != 2 {
		t.Fatalf("Overwrite failed. got %+v", c)
	}
	if c[0][0] != "名前" || c[0][1] != "個数" {
		t.Fatalf("Overwrite failed. got %+v", c)
	}
	if !r.MatchString(c[1][0]) || c[1][1] != "1" {
		t.Fatalf("Overwrite failed. got %+v", c)
	}
	if !r.MatchString(c[2][0]) || c[2][1] != "2" {
		t.Fatalf("Overwrite failed. got %+v", c)
	}
}
