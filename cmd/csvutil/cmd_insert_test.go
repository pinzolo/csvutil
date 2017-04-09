package main

import "testing"

func Example_runInsert() {
	runInsert([]string{"../../testdata/utf8.csv"})
	// Output: column1,名前,個数
	// ,りんご,1
	// ,みかん,2
}

func Test_runInsert(t *testing.T) {
	if c := runInsert([]string{"../../testdata/utf8.csv"}); c != 0 {
		t.Errorf("Invalid success exit code: %d", c)
	}
}

func Test_runInsertOnNoFile(t *testing.T) {
	if c := runInsert([]string{"../../testdata/no-file.csv"}); c == 0 {
		t.Errorf("Invalid failed exit code: %d", c)
	}
}

func Test_runInsertOnFail(t *testing.T) {
	if c := runInsert([]string{"../../testdata/broken.csv"}); c == 0 {
		t.Errorf("Invalid failed exit code: %d", c)
	}
}
