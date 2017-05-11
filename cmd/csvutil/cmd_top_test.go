package main

import "testing"

func Example_runTop() {
	topOpt.Count = 1
	runTop([]string{testFilePath("utf8.csv")})
	topOpt.Count = 0
	// Output: 名前,個数
	// りんご,1
}

func Test_runTop(t *testing.T) {
	topOpt.Count = 1
	if c := runTop([]string{testFilePath("utf8.csv")}); c != 0 {
		t.Errorf("Invalid success exit code: %d", c)
	}
	topOpt.Count = 0
}

func Test_runTopOnNoFile(t *testing.T) {
	if c := runTop([]string{testFilePath("no-file.csv")}); c == 0 {
		t.Errorf("Invalid failed exit code: %d", c)
	}
}

func Test_runTopOnFail(t *testing.T) {
	if c := runTop([]string{testFilePath("broken.csv")}); c == 0 {
		t.Errorf("Invalid failed exit code: %d", c)
	}
}

func Test_runTopOnBackup(t *testing.T) {
	f, err := prepareWritingTest()
	defer f()
	if err != nil {
		t.Error(err)
	}
	topOpt.Count = 1
	topOpt.Overwrite = true
	topOpt.Backup = true
	runTop([]string{tempFilePath()})
	topOpt.Backup = false
	topOpt.Overwrite = false
	topOpt.Count = 0
	if b, err := existsBackup(); err != nil || !b {
		t.Errorf("Failed backup")
	}
}

func Test_runTopOnOverwrite(t *testing.T) {
	f, err := prepareWritingTest()
	defer f()
	if err != nil {
		t.Error(err)
	}
	topOpt.Count = 1
	topOpt.Overwrite = true
	runTop([]string{tempFilePath()})
	topOpt.Overwrite = false
	topOpt.Count = 0
	c, err := overwriteContent()
	if err != nil {
		t.Error(err)
	}
	if len(c) != 2 || len(c[0]) != 2 || c[0][0] != "名前" || c[0][1] != "個数" || c[1][0] != "りんご" || c[1][1] != "1" {
		t.Errorf("Overwrite failed. got %+v", c)
	}
}
