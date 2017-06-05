package main

import "testing"

func Example_runInsert() {
	runInsert([]string{testFilePath("utf8.csv")})
	// Output: column1,名前,個数
	// ,りんご,1
	// ,みかん,2
}

func Test_runInsert(t *testing.T) {
	if c := runInsert([]string{testFilePath("utf8.csv")}); c != 0 {
		t.Fatalf("Invalid success exit code: %d", c)
	}
}

func Test_runInsertOnNoFile(t *testing.T) {
	if c := runInsert([]string{testFilePath("no-file.csv")}); c == 0 {
		t.Fatalf("Invalid failed exit code: %d", c)
	}
}

func Test_runInsertOnFail(t *testing.T) {
	if c := runInsert([]string{testFilePath("broken.csv")}); c == 0 {
		t.Fatalf("Invalid failed exit code: %d", c)
	}
}

func Test_runInsertOnBackup(t *testing.T) {
	f, err := prepareWritingTest()
	defer f()
	if err != nil {
		t.Fatal(err)
	}
	insertOpt.Overwrite = true
	insertOpt.Backup = true
	runInsert([]string{tempFilePath()})
	insertOpt.Backup = false
	insertOpt.Overwrite = false
	if b, err := existsBackup(); err != nil || !b {
		t.Fatalf("Failed backup")
	}
}

func Test_runInsertOnOverwrite(t *testing.T) {
	f, err := prepareWritingTest()
	defer f()
	if err != nil {
		t.Fatal(err)
	}
	insertOpt.Overwrite = true
	runInsert([]string{tempFilePath()})
	insertOpt.Overwrite = false
	c, err := overwriteContent()
	if err != nil {
		t.Fatal(err)
	}
	if len(c[0]) != 3 {
		t.Fatalf("Overwrite failed. got %+v", c)
	}
	if c[0][0] != "column1" || c[0][1] != "名前" || c[0][2] != "個数" {
		t.Fatalf("Overwrite failed. got %+v", c)
	}
	if c[1][0] != "" || c[1][1] != "りんご" || c[1][2] != "1" {
		t.Fatalf("Overwrite failed. got %+v", c)
	}
	if c[2][0] != "" || c[2][1] != "みかん" || c[2][2] != "2" {
		t.Fatalf("Overwrite failed. got %+v", c)
	}
}
