package main

import "testing"

func Example_runCollect() {
	collectOpt.Column = "1"
	runCollect([]string{testFilePath("utf8.csv")})
	collectOpt.Column = ""
	// Output: 1
	// 2
}

func Test_runCollect(t *testing.T) {
	collectOpt.Column = "1"
	if c := runCollect([]string{testFilePath("utf8.csv")}); c != 0 {
		t.Errorf("Invalid success exit code: %d", c)
	}
	collectOpt.Column = ""
}

func Test_runCollectOnNoFile(t *testing.T) {
	collectOpt.Column = "1"
	if c := runCollect([]string{testFilePath("no-file.csv")}); c == 0 {
		t.Errorf("Invalid failed exit code: %d", c)
	}
	collectOpt.Column = ""
}

func Test_runCollectOnFail(t *testing.T) {
	collectOpt.Column = "1"
	if c := runCollect([]string{testFilePath("broken.csv")}); c == 0 {
		t.Errorf("Invalid failed exit code: %d", c)
	}
	collectOpt.Column = ""
}

func Test_runCollectOnBackup(t *testing.T) {
	f, err := prepareWritingTest()
	defer f()
	if err != nil {
		t.Error(err)
	}
	collectOpt.Column = "1"
	collectOpt.Overwrite = true
	collectOpt.Backup = true
	runCollect([]string{tempFilePath()})
	collectOpt.Backup = false
	collectOpt.Overwrite = false
	collectOpt.Column = ""
	if b, err := existsBackup(); err != nil || !b {
		t.Errorf("Failed backup")
	}
}

func Test_runCollectOnOverwrite(t *testing.T) {
	f, err := prepareWritingTest()
	defer f()
	if err != nil {
		t.Error(err)
	}
	collectOpt.Column = "1"
	collectOpt.Overwrite = true
	runCollect([]string{tempFilePath()})
	collectOpt.Overwrite = false
	collectOpt.Column = ""
	c, err := overwriteContent()
	if err != nil {
		t.Error(err)
	}
	if len(c[0]) != 1 {
		t.Errorf("Overwrite failed. got %+v", c)
	}
	if c[0][0] != "1" {
		t.Errorf("Overwrite failed. got %+v", c)
	}
	if c[1][0] != "2" {
		t.Errorf("Overwrite failed. got %+v", c)
	}
}
