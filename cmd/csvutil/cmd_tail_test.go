package main

import "testing"

func Example_runTail() {
	tailOpt.Count = 1
	runTail([]string{testFilePath("utf8.csv")})
	tailOpt.Count = 0
	// Output: 名前,個数
	// みかん,2
}

func Test_runTail(t *testing.T) {
	tailOpt.Count = 1
	if c := runTail([]string{testFilePath("utf8.csv")}); c != 0 {
		t.Errorf("Invalid success exit code: %d", c)
	}
	tailOpt.Count = 0
}

func Test_runTailOnNoFile(t *testing.T) {
	if c := runTail([]string{testFilePath("no-file.csv")}); c == 0 {
		t.Errorf("Invalid failed exit code: %d", c)
	}
}

func Test_runTailOnFail(t *testing.T) {
	if c := runTail([]string{testFilePath("broken.csv")}); c == 0 {
		t.Errorf("Invalid failed exit code: %d", c)
	}
}

func Test_runTailOnBackup(t *testing.T) {
	f, err := prepareWritingTest()
	defer f()
	if err != nil {
		t.Error(err)
	}
	tailOpt.Count = 1
	tailOpt.Overwrite = true
	tailOpt.Backup = true
	runTail([]string{tempFilePath()})
	tailOpt.Backup = false
	tailOpt.Overwrite = false
	tailOpt.Count = 0
	if b, err := existsBackup(); err != nil || !b {
		t.Errorf("Failed backup")
	}
}

func Test_runTailOnOverwrite(t *testing.T) {
	f, err := prepareWritingTest()
	defer f()
	if err != nil {
		t.Error(err)
	}
	tailOpt.Count = 1
	tailOpt.Overwrite = true
	runTail([]string{tempFilePath()})
	tailOpt.Overwrite = false
	tailOpt.Count = 0
	c, err := overwriteContent()
	if err != nil {
		t.Error(err)
	}
	if len(c) != 2 || len(c[0]) != 2 || c[0][0] != "名前" || c[0][1] != "個数" || c[1][0] != "みかん" || c[1][1] != "2" {
		t.Errorf("Overwrite failed. got %+v", c)
	}
}
