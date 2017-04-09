package main

import "testing"

func Example_runSize() {
	runSize([]string{testFilePath("utf8.csv")})
	// Output: 2
}

func Test_runSize(t *testing.T) {
	if c := runSize([]string{testFilePath("utf8.csv")}); c != 0 {
		t.Errorf("Invalid success exit code: %d", c)
	}
}

func Test_runSizeOnNoFile(t *testing.T) {
	if c := runSize([]string{testFilePath("no-file.csv")}); c == 0 {
		t.Errorf("Invalid failed exit code: %d", c)
	}
}

func Test_runSizeOnFail(t *testing.T) {
	if c := runSize([]string{testFilePath("broken.csv")}); c == 0 {
		t.Errorf("Invalid failed exit code: %d", c)
	}
}
