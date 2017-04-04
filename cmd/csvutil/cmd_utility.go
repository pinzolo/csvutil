package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"

	"golang.org/x/crypto/ssh/terminal"
)

func handleError(err error) int {
	fmt.Fprintln(os.Stderr, err)
	return 2
}

func path(args []string) (string, error) {
	if len(args) == 0 {
		if !terminal.IsTerminal(0) {
			return "", nil
		}
		return "", errors.New("Required file path or CSV source.")
	}
	return args[0], nil
}

func writer(path string, overwrite bool) (io.Writer, func(), error) {
	if path == "" {
		return os.Stdout, nil, nil
	}
	if overwrite {
		tmp, err := ioutil.TempFile("", "")
		if err != nil {
			return nil, nil, errors.Wrap(err, "Failed create tempfile")
		}
		finisher := func() {
			tmp.Close()
			os.Rename(tmp.Name(), path)
		}
		return tmp, finisher, nil
	}
	return os.Stdout, nil, nil
}

func reader(path string, bak bool) (io.Reader, func(), error) {
	if path == "" {
		return os.Stdin, nil, nil
	}
	if bak {
		src, err := backup(path)
		if err != nil {
			return nil, nil, errors.Wrap(err, "Failed backup")
		}
		return openWithCloser(src)
	}

	src, err := os.Open(path)
	if err != nil {
		return nil, nil, errors.Wrap(err, "Failed open file")
	}
	return openWithCloser(src)
}

func openWithCloser(f *os.File) (io.Reader, func(), error) {
	closer := func() {
		f.Close()
	}

	return f, closer, nil
}

func backup(path string) (*os.File, error) {
	ext := filepath.Ext(path)
	dst := strings.TrimSuffix(path, ext) + "." + time.Now().Format("20060102150405") + ext
	err := os.Rename(path, dst)
	if err != nil {
		return nil, errors.Wrap(err, "Cannot move")
	}
	return os.Open(dst)
}
