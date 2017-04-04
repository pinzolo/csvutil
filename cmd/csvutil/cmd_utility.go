package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/pinzolo/csvutil"
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
			return nil, nil, err
		}
		return tmp, func() { os.Rename(tmp.Name(), path) }, nil
	}
	return os.Stdout, nil, nil
}

func reader(path string, bak bool) (io.Reader, error) {
	if path == "" {
		return os.Stdin, nil
	}
	if !bak {
		return os.Open(path)
	}

	bp, err := csvutil.Backup(path)
	if err != nil {
		return nil, err
	}
	return os.Open(bp)
}
