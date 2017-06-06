package main

import "testing"

func Example_runSort() {
	sortOpt.Column = "aaa"
	runSort([]string{testFilePath("sort.csv")})
	sortOpt.Column = ""
	// Output: aaa,bbb,ccc
	// ,b2,c2
	// ,b4,c4
	// 1,b5,c5
	// 10,b3,c3
	// 2,b1,c1
}

func Test_runSort(t *testing.T) {
	sortOpt.Column = "aaa"
	if c := runSort([]string{testFilePath("sort.csv")}); c != 0 {
		t.Fatalf("Invalid success exit code: %d", c)
	}
	sortOpt.Column = ""
}

func Test_runSortOnNoFile(t *testing.T) {
	if c := runSort([]string{testFilePath("no-file.csv")}); c == 0 {
		t.Fatalf("Invalid failed exit code: %d", c)
	}
}

func Test_runSortOnFail(t *testing.T) {
	if c := runSort([]string{testFilePath("broken.csv")}); c == 0 {
		t.Fatalf("Invalid failed exit code: %d", c)
	}
}

func Test_runSortOnBackup(t *testing.T) {
	f, err := prepareWritingTest()
	defer f()
	if err != nil {
		t.Fatal(err)
	}
	sortOpt.Column = "個数"
	sortOpt.Overwrite = true
	sortOpt.Backup = true
	runSort([]string{tempFilePath()})
	sortOpt.Backup = false
	sortOpt.Overwrite = false
	sortOpt.Column = ""
	if b, err := existsBackup(); err != nil || !b {
		t.Fatalf("Failed backup")
	}
}

func Test_runSortOnOverwrite(t *testing.T) {
	f, err := prepareWritingTest()
	defer f()
	if err != nil {
		t.Fatal(err)
	}
	sortOpt.Column = "個数"
	sortOpt.Descending = true
	sortOpt.Overwrite = true
	runSort([]string{tempFilePath()})
	sortOpt.Overwrite = false
	sortOpt.Descending = false
	sortOpt.Column = ""
	c, err := overwriteContent()
	if err != nil {
		t.Fatal(err)
	}
	sorted := []string{"2", "1"}
	for i, s := range sorted {
		if row := c[i+1]; row[1] != s {
			t.Errorf("expected %s but got %s (index: %d)", s, row[0], i)
		}
	}
}
