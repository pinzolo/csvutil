package main

import (
	"testing"
)

func Example_runName() {
	nameOpt.Gender = "0"
	nameOpt.MaleRate = 100
	runName([]string{testFilePath("utf8.csv")})
	nameOpt.Gender = ""
	nameOpt.MaleRate = 50
	// Output: 名前,個数
	// 男,1
	// 男,2
}

func Test_runName(t *testing.T) {
	nameOpt.Name = "0"
	if c := runName([]string{testFilePath("utf8.csv")}); c != 0 {
		t.Fatalf("Invalid success exit code: %d", c)
	}
	nameOpt.Name = ""
}

func Test_runNameOnNoFile(t *testing.T) {
	nameOpt.Name = "0"
	if c := runName([]string{testFilePath("no-file.csv")}); c == 0 {
		t.Fatalf("Invalid failed exit code: %d", c)
	}
	nameOpt.Name = ""
}

func Test_runNameOnFail(t *testing.T) {
	nameOpt.Name = "0"
	if c := runName([]string{testFilePath("broken.csv")}); c == 0 {
		t.Fatalf("Invalid failed exit code: %d", c)
	}
	nameOpt.Name = ""
}

func Test_runNameOnBackup(t *testing.T) {
	f, err := prepareWritingTest()
	defer f()
	if err != nil {
		t.Fatal(err)
	}
	nameOpt.Name = "0"
	nameOpt.Overwrite = true
	nameOpt.Backup = true
	runName([]string{tempFilePath()})
	nameOpt.Backup = false
	nameOpt.Overwrite = false
	nameOpt.Name = ""
	if b, err := existsBackup(); err != nil || !b {
		t.Fatalf("Failed backup")
	}
}

func Test_runNameOnOverwrite(t *testing.T) {
	f, err := prepareWritingTest()
	defer f()
	if err != nil {
		t.Fatal(err)
	}
	nameOpt.Gender = "0"
	nameOpt.Overwrite = true
	runName([]string{tempFilePath()})
	nameOpt.Overwrite = false
	nameOpt.Gender = ""
	c, err := overwriteContent()
	if err != nil {
		t.Fatal(err)
	}
	if len(c[0]) != 2 {
		t.Fatalf("Overwrite failed. got %+v", c)
	}
	if c[0][0] != "名前" || c[0][1] != "個数" {
		t.Fatalf("Overwrite failed. got %+v", c)
	}
	if (c[1][0] != "女" && c[1][0] != "男") || c[1][1] != "1" {
		t.Fatalf("Overwrite failed. got %+v", c)
	}
	if (c[2][0] != "女" && c[2][0] != "男") || c[1][1] != "1" {
		t.Fatalf("Overwrite failed. got %+v", c)
	}
}
