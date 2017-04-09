package main

import "testing"

func Example_runGenerate() {
	runGenerate([]string{})
	// Output: column1,column2,column3
	// ,,
	// ,,
	// ,,
}

func Test_runGenerate(t *testing.T) {
	if c := runGenerate([]string{}); c != 0 {
		t.Errorf("Invalid success exit code: %d", c)
	}
}

func Test_runGenerateOnFail(t *testing.T) {
	s := generateOpt.Size
	generateOpt.Size = -1
	if c := runGenerate([]string{}); c == 0 {
		t.Errorf("Invalid failed exit code: %d", c)
	}
	generateOpt.Size = s
}
