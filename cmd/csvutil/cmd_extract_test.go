package main

import "testing"

func Example_runExtract() {
	extractOpt.Column = "名前"
	runExtract([]string{testFilePath("utf8.csv")})
	extractOpt.Column = ""
	// Output: 名前
	// りんご
	// みかん
}

func Test_runExtract(t *testing.T) {
	extractOpt.Column = "名前"
	if c := runExtract([]string{testFilePath("utf8.csv")}); c != 0 {
		t.Fatalf("Invalid success exit code: %d", c)
	}
	extractOpt.Column = ""
}

func Test_runExtractOnNoFile(t *testing.T) {
	if c := runExtract([]string{testFilePath("no-file.csv")}); c == 0 {
		t.Fatalf("Invalid failed exit code: %d", c)
	}
}

func Test_runExtractOnFail(t *testing.T) {
	if c := runExtract([]string{testFilePath("broken.csv")}); c == 0 {
		t.Fatalf("Invalid failed exit code: %d", c)
	}
}

func Test_runExtractOnBackup(t *testing.T) {
	f, err := prepareWritingTest()
	defer f()
	if err != nil {
		t.Fatal(err)
	}
	extractOpt.Column = "名前"
	extractOpt.Overwrite = true
	extractOpt.Backup = true
	runExtract([]string{tempFilePath()})
	extractOpt.Backup = false
	extractOpt.Overwrite = false
	extractOpt.Column = ""
	if b, err := existsBackup(); err != nil || !b {
		t.Fatalf("Failed backup")
	}
}

func Test_runExtractOnOverwrite(t *testing.T) {
	f, err := prepareWritingTest()
	defer f()
	if err != nil {
		t.Fatal(err)
	}
	extractOpt.Column = "名前"
	extractOpt.Overwrite = true
	runExtract([]string{tempFilePath()})
	extractOpt.Overwrite = false
	extractOpt.Column = ""
	c, err := overwriteContent()
	if err != nil {
		t.Fatal(err)
	}
	if len(c[0]) != 1 || c[0][0] != "名前" || c[1][0] != "りんご" || c[2][0] != "みかん" {
		t.Fatalf("Overwrite failed. got %+v", c)
	}
}
