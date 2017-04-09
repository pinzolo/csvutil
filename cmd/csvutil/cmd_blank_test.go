package main

import "testing"

func Example_runBlank() {
	blankOpt.Column = "0"
	runBlank([]string{testFilePath("utf8.csv")})
	blankOpt.Column = ""
	// Output: 名前,個数
	// ,1
	// ,2
}

func Test_runBlank(t *testing.T) {
	blankOpt.Column = "0"
	if c := runBlank([]string{testFilePath("utf8.csv")}); c != 0 {
		t.Errorf("Invalid success exit code: %d", c)
	}
	blankOpt.Column = ""
}

func Test_runBlankOnNoFile(t *testing.T) {
	if c := runBlank([]string{testFilePath("no-file.csv")}); c == 0 {
		t.Errorf("Invalid failed exit code: %d", c)
	}
}

func Test_runBlankOnFail(t *testing.T) {
	if c := runBlank([]string{testFilePath("broken.csv")}); c == 0 {
		t.Errorf("Invalid failed exit code: %d", c)
	}
}

func Test_runBlankOnBackup(t *testing.T) {
	f, err := prepareWritingTest()
	defer f()
	if err != nil {
		t.Error(err)
	}
	blankOpt.Column = "名前"
	blankOpt.Overwrite = true
	blankOpt.Backup = true
	runBlank([]string{tempFilePath()})
	blankOpt.Backup = false
	blankOpt.Overwrite = false
	blankOpt.Column = ""
	if b, err := existsBackup(); err != nil || !b {
		t.Errorf("Failed backup")
	}
}

func Test_runBlankOnOverwrite(t *testing.T) {
	f, err := prepareWritingTest()
	defer f()
	if err != nil {
		t.Error(err)
	}
	blankOpt.Column = "名前"
	blankOpt.Overwrite = true
	runBlank([]string{tempFilePath()})
	blankOpt.Overwrite = false
	blankOpt.Column = ""
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
	if c[1][0] != "" || c[1][1] != "1" {
		t.Errorf("Overwrite failed. got %+v", c)
	}
	if c[2][0] != "" || c[2][1] != "2" {
		t.Errorf("Overwrite failed. got %+v", c)
	}
}
