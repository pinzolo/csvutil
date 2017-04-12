package main

import "testing"

func Example_runCombine() {
	combineOpt.Source = "0:1"
	combineOpt.Destination = "1"
	runCombine([]string{testFilePath("utf8.csv")})
	combineOpt.Source = ""
	combineOpt.Destination = ""
	// Output: 名前,個数
	// りんご,りんご1
	// みかん,みかん2
}

func Test_runCombine(t *testing.T) {
	combineOpt.Source = "0:1"
	combineOpt.Destination = "1"
	runCombine([]string{testFilePath("utf8.csv")})
	if c := runCombine([]string{testFilePath("utf8.csv")}); c != 0 {
		t.Errorf("Invalid success exit code: %d", c)
	}
	combineOpt.Source = ""
	combineOpt.Destination = ""
}

func Test_runCombineOnNoFile(t *testing.T) {
	combineOpt.Source = "0:1"
	combineOpt.Destination = "1"
	runCombine([]string{testFilePath("utf8.csv")})
	if c := runCombine([]string{testFilePath("no-file.csv")}); c == 0 {
		t.Errorf("Invalid failed exit code: %d", c)
	}
	combineOpt.Source = ""
	combineOpt.Destination = ""
}

func Test_runCombineOnFail(t *testing.T) {
	combineOpt.Source = "0:1"
	combineOpt.Destination = "1"
	runCombine([]string{testFilePath("utf8.csv")})
	if c := runCombine([]string{testFilePath("broken.csv")}); c == 0 {
		t.Errorf("Invalid failed exit code: %d", c)
	}
	combineOpt.Source = ""
	combineOpt.Destination = ""
}

func Test_runCombineOnBackup(t *testing.T) {
	f, err := prepareWritingTest()
	defer f()
	if err != nil {
		t.Error(err)
	}
	combineOpt.Source = "1:0"
	combineOpt.Destination = "1"
	combineOpt.Overwrite = true
	combineOpt.Backup = true
	runCombine([]string{tempFilePath()})
	combineOpt.Backup = false
	combineOpt.Overwrite = false
	combineOpt.Destination = ""
	combineOpt.Source = ""
	if b, err := existsBackup(); err != nil || !b {
		t.Errorf("Failed backup")
	}
}

func Test_runCombineOnOverwrite(t *testing.T) {
	f, err := prepareWritingTest()
	defer f()
	if err != nil {
		t.Error(err)
	}
	combineOpt.Source = "1:0"
	combineOpt.Destination = "1"
	combineOpt.Overwrite = true
	runCombine([]string{tempFilePath()})
	combineOpt.Overwrite = false
	combineOpt.Destination = ""
	combineOpt.Source = ""
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
	if c[1][0] != "りんご" || c[1][1] != "1りんご" {
		t.Errorf("Overwrite failed. got %+v", c)
	}
	if c[2][0] != "みかん" || c[2][1] != "2みかん" {
		t.Errorf("Overwrite failed. got %+v", c)
	}
}
