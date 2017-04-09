package main

import (
	"testing"
)

func Example_runRemove() {
	removeOpt.Column = "名前"
	runRemove([]string{testFilePath("utf8.csv")})
	removeOpt.Column = ""
	// Output: 個数
	// 1
	// 2
}

func Test_runRemove(t *testing.T) {
	removeOpt.Column = "名前"
	if c := runRemove([]string{testFilePath("utf8.csv")}); c != 0 {
		t.Errorf("Invalid success exit code: %d", c)
	}
	removeOpt.Column = ""
}

func Test_runRemoveOnNoFile(t *testing.T) {
	if c := runRemove([]string{testFilePath("no-file.csv")}); c == 0 {
		t.Errorf("Invalid failed exit code: %d", c)
	}
}

func Test_runRemoveOnFail(t *testing.T) {
	if c := runRemove([]string{testFilePath("broken.csv")}); c == 0 {
		t.Errorf("Invalid failed exit code: %d", c)
	}
}

func Test_runRemoveOnBackup(t *testing.T) {
	f, err := prepareWritingTest()
	defer f()
	if err != nil {
		t.Error(err)
	}
	removeOpt.Column = "名前"
	removeOpt.Overwrite = true
	removeOpt.Backup = true
	runRemove([]string{tempFilePath()})
	removeOpt.Backup = false
	removeOpt.Overwrite = false
	removeOpt.Column = ""
	if b, err := existsBackup(); err != nil || !b {
		t.Errorf("Failed backup")
	}
}

func Test_runRemoveOnOverwrite(t *testing.T) {
	f, err := prepareWritingTest()
	defer f()
	if err != nil {
		t.Error(err)
	}
	removeOpt.Column = "名前"
	removeOpt.Overwrite = true
	runRemove([]string{tempFilePath()})
	removeOpt.Overwrite = false
	removeOpt.Column = ""
	c, err := overwriteContent()
	if err != nil {
		t.Error(err)
	}
	if len(c[0]) != 1 || c[0][0] != "個数" || c[1][0] != "1" || c[2][0] != "2" {
		t.Errorf("Overwrite failed. got %+v", c)
	}
}
