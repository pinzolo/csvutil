package main

import "testing"

var passwordRunes = []rune(`abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!'@#$%^&*()_+-=[]{};:",./?`)

func Test_runPassword(t *testing.T) {
	passwordOpt.Column = "名前"
	if c := runPassword([]string{testFilePath("utf8.csv")}); c != 0 {
		t.Fatalf("Invalid success exit code: %d", c)
	}
	passwordOpt.Column = ""
}

func Test_runPasswordOnNoFile(t *testing.T) {
	if c := runPassword([]string{testFilePath("no-file.csv")}); c == 0 {
		t.Fatalf("Invalid failed exit code: %d", c)
	}
}

func Test_runPasswordOnFail(t *testing.T) {
	if c := runPassword([]string{testFilePath("broken.csv")}); c == 0 {
		t.Fatalf("Invalid failed exit code: %d", c)
	}
}

func Test_runPasswordOnBackup(t *testing.T) {
	f, err := prepareWritingTest()
	defer f()
	if err != nil {
		t.Fatal(err)
	}
	passwordOpt.Column = "名前"
	passwordOpt.Overwrite = true
	passwordOpt.Backup = true
	runPassword([]string{tempFilePath()})
	passwordOpt.Backup = false
	passwordOpt.Overwrite = false
	passwordOpt.Column = ""
	if b, err := existsBackup(); err != nil || !b {
		t.Fatalf("Failed backup")
	}
}

func Test_runPasswordOnOverwrite(t *testing.T) {
	f, err := prepareWritingTest()
	defer f()
	if err != nil {
		t.Fatal(err)
	}
	passwordOpt.Column = "名前"
	passwordOpt.Overwrite = true
	runPassword([]string{tempFilePath()})
	passwordOpt.Overwrite = false
	passwordOpt.Column = ""
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
	if !isPassword(c[1][0]) || c[1][1] != "1" {
		t.Fatalf("Overwrite failed. got %+v", c)
	}
	if !isPassword(c[2][0]) || c[2][1] != "2" {
		t.Fatalf("Overwrite failed. got %+v", c)
	}
}

func isPassword(s string) bool {
	for _, r := range []rune(s) {
		b := false
		for _, r2 := range passwordRunes {
			if r2 == r {
				b = true
			}
		}
		if !b {
			return false
		}
	}
	return true
}
