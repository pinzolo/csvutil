package main

import "testing"

func Example_runCount() {
	runCount([]string{testFilePath("utf8.csv")})
	// Output: 2
}

func Test_runCount(t *testing.T) {
	if c := runCount([]string{testFilePath("utf8.csv")}); c != 0 {
		t.Fatalf("Invalid success exit code: %d", c)
	}
}

func Test_runCountOnNoFile(t *testing.T) {
	if c := runCount([]string{testFilePath("no-file.csv")}); c == 0 {
		t.Fatalf("Invalid failed exit code: %d", c)
	}
}

func Test_runCountOnFail(t *testing.T) {
	if c := runCount([]string{testFilePath("broken.csv")}); c == 0 {
		t.Fatalf("Invalid failed exit code: %d", c)
	}
}
