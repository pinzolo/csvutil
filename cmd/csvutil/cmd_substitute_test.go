package main

import "testing"

func Example_runSubstitute() {
	substituteOpt.Column = "1"
	substituteOpt.Pattern = `\d`
	substituteOpt.Replacement = "xxx"
	substituteOpt.Regexp = true
	runSubstitute([]string{testFilePath("utf8.csv")})
	substituteOpt.Regexp = false
	substituteOpt.Replacement = ""
	substituteOpt.Pattern = ""
	substituteOpt.Column = ""
	// Output: 名前,個数
	// りんご,xxx
	// みかん,xxx
}

func Test_runSubstitute(t *testing.T) {
	substituteOpt.Column = "1"
	substituteOpt.Pattern = `\d`
	substituteOpt.Replacement = "xxx"
	substituteOpt.Regexp = true
	if c := runSubstitute([]string{testFilePath("utf8.csv")}); c != 0 {
		t.Errorf("Invalid success exit code: %d", c)
	}
	substituteOpt.Regexp = false
	substituteOpt.Replacement = ""
	substituteOpt.Pattern = ""
	substituteOpt.Column = ""
}

func Test_runSubstituteOnNoFile(t *testing.T) {
	substituteOpt.Column = "1"
	substituteOpt.Pattern = `\d`
	substituteOpt.Replacement = "xxx"
	substituteOpt.Regexp = true
	if c := runSubstitute([]string{testFilePath("no-file.csv")}); c == 0 {
		t.Errorf("Invalid failed exit code: %d", c)
	}
	substituteOpt.Regexp = false
	substituteOpt.Replacement = ""
	substituteOpt.Pattern = ""
	substituteOpt.Column = ""
}

func Test_runSubstituteOnFail(t *testing.T) {
	substituteOpt.Column = "1"
	substituteOpt.Pattern = `\d`
	substituteOpt.Replacement = "xxx"
	substituteOpt.Regexp = true
	if c := runSubstitute([]string{testFilePath("broken.csv")}); c == 0 {
		t.Errorf("Invalid failed exit code: %d", c)
	}
	substituteOpt.Regexp = false
	substituteOpt.Replacement = ""
	substituteOpt.Pattern = ""
	substituteOpt.Column = ""
}

func Test_runSubstituteOnBackup(t *testing.T) {
	f, err := prepareWritingTest()
	defer f()
	if err != nil {
		t.Error(err)
	}
	substituteOpt.Column = "1"
	substituteOpt.Pattern = `\d`
	substituteOpt.Replacement = "xxx"
	substituteOpt.Regexp = true
	substituteOpt.Overwrite = true
	substituteOpt.Backup = true
	runSubstitute([]string{tempFilePath()})
	substituteOpt.Backup = false
	substituteOpt.Overwrite = false
	substituteOpt.Regexp = false
	substituteOpt.Replacement = ""
	substituteOpt.Pattern = ""
	substituteOpt.Column = ""
	if b, err := existsBackup(); err != nil || !b {
		t.Errorf("Failed backup")
	}
}

func Test_runSubstituteOnOverwrite(t *testing.T) {
	f, err := prepareWritingTest()
	defer f()
	if err != nil {
		t.Error(err)
	}
	substituteOpt.Column = "1"
	substituteOpt.Pattern = `\d`
	substituteOpt.Replacement = "xxx"
	substituteOpt.Regexp = true
	substituteOpt.Overwrite = true
	runSubstitute([]string{tempFilePath()})
	substituteOpt.Overwrite = false
	substituteOpt.Regexp = false
	substituteOpt.Replacement = ""
	substituteOpt.Pattern = ""
	substituteOpt.Column = ""
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
	if c[1][0] != "りんご" || c[1][1] != "xxx" {
		t.Errorf("Overwrite failed. got %+v", c)
	}
	if c[2][0] != "みかん" || c[2][1] != "xxx" {
		t.Errorf("Overwrite failed. got %+v", c)
	}
}
