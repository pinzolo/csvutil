package main

import (
	"testing"
)

func Test_runBuilding(t *testing.T) {
	buildingOpt.Column = "名前"
	if c := runBuilding([]string{testFilePath("utf8.csv")}); c != 0 {
		t.Fatalf("Invalid success exit code: %d", c)
	}
	buildingOpt.Column = ""
}

func Test_runBuildingOnNoFile(t *testing.T) {
	if c := runBuilding([]string{testFilePath("no-file.csv")}); c == 0 {
		t.Fatalf("Invalid failed exit code: %d", c)
	}
}

func Test_runBuildingOnFail(t *testing.T) {
	if c := runBuilding([]string{testFilePath("broken.csv")}); c == 0 {
		t.Fatalf("Invalid failed exit code: %d", c)
	}
}

func Test_runBuildingOnBackup(t *testing.T) {
	f, err := prepareWritingTest()
	defer f()
	if err != nil {
		t.Fatal(err)
	}
	buildingOpt.Column = "名前"
	buildingOpt.Overwrite = true
	buildingOpt.Backup = true
	runBuilding([]string{tempFilePath()})
	buildingOpt.Backup = false
	buildingOpt.Overwrite = false
	buildingOpt.Column = ""
	if b, err := existsBackup(); err != nil || !b {
		t.Fatalf("Failed backup")
	}
}

func Test_runBuildingOnOverwrite(t *testing.T) {
	f, err := prepareWritingTest()
	defer f()
	if err != nil {
		t.Fatal(err)
	}
	buildingOpt.Column = "名前"
	buildingOpt.Overwrite = true
	runBuilding([]string{tempFilePath()})
	buildingOpt.Overwrite = false
	buildingOpt.Column = ""
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
	if c[1][0] == "りんご" || c[1][1] != "1" {
		t.Fatalf("Overwrite failed. got %+v", c)
	}
	if c[2][0] == "みかん" || c[2][1] != "2" {
		t.Fatalf("Overwrite failed. got %+v", c)
	}
}
