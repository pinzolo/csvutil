package main

import "testing"

func Example_runAppend() {
	runAppend([]string{testFilePath("utf8.csv")})
	// Output: 名前,個数,column1
	// りんご,1,
	// みかん,2,
}

func Test_runAppend(t *testing.T) {
	if c := runAppend([]string{testFilePath("utf8.csv")}); c != 0 {
		t.Errorf("Invalid success exit code: %d", c)
	}
}

func Test_runAppendOnNoFile(t *testing.T) {
	if c := runAppend([]string{testFilePath("no-file.csv")}); c == 0 {
		t.Errorf("Invalid failed exit code: %d", c)
	}
}

func Test_runAppendOnFail(t *testing.T) {
	if c := runAppend([]string{testFilePath("broken.csv")}); c == 0 {
		t.Errorf("Invalid failed exit code: %d", c)
	}
}

func Test_runAppendOnBackup(t *testing.T) {
	f, err := prepareWritingTest()
	defer f()
	if err != nil {
		t.Error(err)
	}
	appendOpt.Overwrite = true
	appendOpt.Backup = true
	runAppend([]string{tempFilePath()})
	appendOpt.Backup = false
	appendOpt.Overwrite = false
	if b, err := existsBackup(); err != nil || !b {
		t.Errorf("Failed backup")
	}
}

func Test_runAppendOnOverwrite(t *testing.T) {
	f, err := prepareWritingTest()
	defer f()
	if err != nil {
		t.Error(err)
	}
	appendOpt.Overwrite = true
	runAppend([]string{tempFilePath()})
	appendOpt.Overwrite = false
	c, err := overwriteContent()
	if err != nil {
		t.Error(err)
	}
	if len(c[0]) != 3 {
		t.Errorf("Overwrite failed. got %+v", c)
	}
	if c[0][0] != "名前" || c[0][1] != "個数" || c[0][2] != "column1" {
		t.Errorf("Overwrite failed. got %+v", c)
	}
	if c[1][0] != "りんご" || c[1][1] != "1" || c[1][2] != "" {
		t.Errorf("Overwrite failed. got %+v", c)
	}
	if c[2][0] != "みかん" || c[2][1] != "2" || c[2][2] != "" {
		t.Errorf("Overwrite failed. got %+v", c)
	}
}
