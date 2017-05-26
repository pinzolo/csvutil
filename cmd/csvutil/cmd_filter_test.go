package main

import "testing"

func Example_runFilter() {
	filterOpt.Pattern = `[A-E]`
	filterOpt.Regexp = true
	runFilter([]string{testFilePath("filter.csv")})
	filterOpt.Regexp = false
	filterOpt.Pattern = ""
	// Output: aaa,bbb,ccc
	// A1,B2,C3
	// D4,E1,F6
}

func Test_runFilter(t *testing.T) {
	filterOpt.Pattern = `[A-E]`
	filterOpt.Regexp = true
	if c := runFilter([]string{testFilePath("filter.csv")}); c != 0 {
		t.Errorf("Invalid success exit code: %d", c)
	}
	filterOpt.Regexp = false
	filterOpt.Pattern = ""
}

func Test_runFilterOnNoFile(t *testing.T) {
	filterOpt.Pattern = `[A-E]`
	filterOpt.Regexp = true
	if c := runFilter([]string{testFilePath("no-file.csv")}); c == 0 {
		t.Errorf("Invalid failed exit code: %d", c)
	}
	filterOpt.Regexp = false
	filterOpt.Pattern = ""
}

func Test_runFilterOnFail(t *testing.T) {
	filterOpt.Column = "1"
	filterOpt.Pattern = `\d`
	filterOpt.Regexp = true
	if c := runFilter([]string{testFilePath("broken.csv")}); c == 0 {
		t.Errorf("Invalid failed exit code: %d", c)
	}
	filterOpt.Regexp = false
	filterOpt.Pattern = ""
	filterOpt.Column = ""
}

func Test_runFilterOnBackup(t *testing.T) {
	f, err := prepareWritingTest()
	defer f()
	if err != nil {
		t.Error(err)
	}
	filterOpt.Pattern = "2"
	filterOpt.Overwrite = true
	filterOpt.Backup = true
	runFilter([]string{tempFilePath()})
	filterOpt.Backup = false
	filterOpt.Overwrite = false
	filterOpt.Pattern = ""
	if b, err := existsBackup(); err != nil || !b {
		t.Errorf("Failed backup")
	}
}

func Test_runFilterOnOverwrite(t *testing.T) {
	f, err := prepareWritingTest()
	defer f()
	if err != nil {
		t.Error(err)
	}
	filterOpt.Pattern = "2"
	filterOpt.Regexp = true
	filterOpt.Overwrite = true
	runFilter([]string{tempFilePath()})
	filterOpt.Overwrite = false
	filterOpt.Regexp = false
	filterOpt.Pattern = ""
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
	if c[1][0] != "みかん" || c[1][1] != "2" {
		t.Errorf("Overwrite failed. got %+v", c)
	}
}
