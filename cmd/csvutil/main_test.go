package main

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	"github.com/pinzolo/csvutil"
)

const (
	testdataDir = "../../testdata"
	tempDir     = "../../testdata/temp"
)

var bakFileRegex = regexp.MustCompile(`utf8\.\d{14}\.csv`)

func copyTestFile() error {
	src, err := os.Open(testFilePath("utf8.csv"))
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(tempFilePath())
	if err != nil {
		return nil
	}
	defer dst.Close()
	_, err = io.Copy(dst, src)
	return err
}

func prepareWritingTest() (func(), error) {
	if _, err := os.Stat(tempDir); err == nil {
		os.RemoveAll(tempDir)
	}
	err := os.Mkdir(tempDir, 0777)
	if err != nil {
		return nil, err
	}
	err = copyTestFile()
	f := func() {
		ferr := os.RemoveAll(tempDir)
		if ferr != nil {
			panic(ferr)
		}
	}
	return f, err
}

func testFilePath(n string) string {
	return filepath.Join(testdataDir, n)
}

func tempFilePath() string {
	return filepath.Join(filepath.Join(tempDir, "utf8.csv"))
}

func overwriteContent() ([][]string, error) {
	f, err := os.Open(tempFilePath())
	if err != nil {
		return nil, err
	}
	r, _ := csvutil.NewReader(f)
	return r.ReadAll()
}

func existsBackup() (bool, error) {
	files, err := ioutil.ReadDir(tempDir)
	if err != nil {
		return false, err
	}
	for _, f := range files {
		if bakFileRegex.MatchString(f.Name()) {
			return true, nil
		}
	}
	return false, nil
}
