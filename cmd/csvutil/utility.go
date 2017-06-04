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

var errWriter = io.Writer(os.Stderr)

func handleError(err error) int {
	fmt.Fprintln(errWriter, err)
	return 2
}

func path(args []string) (string, error) {
	if len(args) == 0 {
		if !terminal.IsTerminal(0) {
			return "", nil
		}
		return "", errors.New("no file and source")
	}
	return args[0], nil
}

func writer(path string, overwrite bool) (io.Writer, func(*bool, bool), error) {
	if path == "" {
		return os.Stdout, nil, nil
	}
	if overwrite {
		tmp, err := ioutil.TempFile("", "")
		if err != nil {
			return nil, nil, errors.Wrap(err, "failed create tempfile")
		}
		finisher := func(success *bool, bak bool) {
			tmp.Close()
			if *success {
				if bak {
					backup(path)
				}
				if overwrite {
					os.Rename(tmp.Name(), path)
				}
			}
		}
		return tmp, finisher, nil
	}
	return os.Stdout, nil, nil
}

func reader(path string) (io.Reader, func(), error) {
	if path == "" {
		return os.Stdin, nil, nil
	}

	src, err := os.Open(path)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed open")
	}
	return withCloser(src)
}

func withCloser(f *os.File) (io.Reader, func(), error) {
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
		return nil, errors.Wrap(err, "failed backup")
	}
	return os.Open(dst)
}

func split(s string) []string {
	if s == "" {
		return nil
	}
	return strings.Split(s, ":")
}

func prepare(args []string, ow bool) (io.Writer, func(*bool, bool), io.Reader, func(), error) {
	path, err := path(args)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	w, wf, err := writer(path, ow)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	r, rf, err := reader(path)
	if err != nil {
		return w, wf, nil, nil, err
	}
	return w, wf, r, rf, nil
}

func prepareReader(args []string) (io.Reader, func(), error) {
	path, err := path(args)
	if err != nil {
		return nil, nil, err
	}

	r, rf, err := reader(path)
	if err != nil {
		return nil, nil, err
	}
	return r, rf, nil
}
