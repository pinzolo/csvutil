package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/pinzolo/csvutil"
	"github.com/pkg/errors"
)

var cmdBlank = &Command{
	Run:       runBlank,
	UsageLine: "blank [OPTIONS...] [FILE]",
	Short:     "Replace column value(s) by empty or blank string.",
	Long: `Replace column value(s) by empty or blank string.
ARGUMENTS
        FILE
            Source CSV file.
            Without FILE argument, read from STDIN.

OPTIONS
        -w, --overwrite
            Overwrite source file by replaced CSV.
            This option does not work when file is not given.

        -H, --no-header
            Tel given CSV does not have header line.

        -b, --backup
            Create backup file before replace.
            This option should be used with --overwrite option.

        -e, --encoding
            Encoding of source file.
            This option accepts 'sjis' or 'eucjp'.
            Without this option, csvutil treats CSV file is encoded by UTF-8.

        -c, --column
            Target column symbol(s).
            Column symbol accepts column index or column header text.
            If --no-header option is used, this option accepts only column index.
            To target multi columns, use semicolon seperated value like foo:bar and 1:2.

        -r, --rate
            Percentage of replace rate. Without this option, always replace CSV value.
            Use this option to make discreta data.

        -sw, --space-width
            Width of space character.
            0: empty string (default)
            1: ASCII space [0x20]
            2: Multi byte space [0xE3 0x80 0x80]

        --ss, --space-size
            Count of space characters. (default 0)
            If space size is 2 and space width is 1 then value replaced by "  ". (2 ASCII space characters).
            If space size is 3 and space width is 2 then value replaced by "　　　". (3 multi byte space characters).
	`,
}

var blankOpt = csvutil.BlankOption{}

func init() {
	cmdBlank.Flag.BoolVar(&blankOpt.Overwrite, "overwrite", false, "Overwrite to source.")
	cmdBlank.Flag.BoolVar(&blankOpt.Overwrite, "w", false, "Overwrite to source.")
	cmdBlank.Flag.BoolVar(&blankOpt.NoHeader, "no-header", false, "Source file does not have header line.")
	cmdBlank.Flag.BoolVar(&blankOpt.NoHeader, "H", false, "Source file does not have header line.")
	cmdBlank.Flag.BoolVar(&blankOpt.Backup, "backup", false, "Backup source file.")
	cmdBlank.Flag.BoolVar(&blankOpt.Backup, "b", false, "Backup source file.")
	cmdBlank.Flag.StringVar(&blankOpt.Encoding, "encoding", "utf8", "Encoding of source file")
	cmdBlank.Flag.StringVar(&blankOpt.Encoding, "e", "utf8", "Encoding of source file")
	cmdBlank.Flag.StringVar(&blankOpt.Column, "column", "", "Column symbol")
	cmdBlank.Flag.StringVar(&blankOpt.Column, "c", "", "Column symbol")
	cmdBlank.Flag.IntVar(&blankOpt.Rate, "rate", 100, "Filling rate")
	cmdBlank.Flag.IntVar(&blankOpt.Rate, "r", 100, "Filling rate")
	cmdBlank.Flag.IntVar(&blankOpt.SpaceWidth, "space-width", 0, "Space character width")
	cmdBlank.Flag.IntVar(&blankOpt.SpaceWidth, "sw", 0, "Space character width")
	cmdBlank.Flag.IntVar(&blankOpt.SpaceSize, "space-size", 0, "Space character count")
	cmdBlank.Flag.IntVar(&blankOpt.SpaceSize, "ss", 0, "Space character count")
}

// runBlank executes blank command and return exit code.
func runBlank(args []string) int {
	path, err := path(args)
	if err != nil {
		return handleError(err)
	}

	w, cleanup, err := writer(path, blankOpt.Overwrite)
	if err != nil {
		return handleError(err)
	}
	if cleanup != nil {
		defer cleanup()
	}

	r, err := reader(path, blankOpt.Backup)
	if err != nil {
		return handleError(err)
	}

	err = csvutil.Blank(r, w, blankOpt)
	if err != nil {
		return handleError(err)
	}

	return 0
}

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
	if blankOpt.Overwrite {
		tmp, err := ioutil.TempFile("", "")
		if err != nil {
			return nil, nil, err
		}
		return tmp, func() { os.Rename(tmp.Name(), path) }, nil
	}
	return os.Stdout, nil, nil
}

func reader(path string, backup bool) (io.Reader, error) {
	if path == "" {
		return os.Stdin, nil
	}
	if !backup {
		return os.Open(path)
	}

	bp, err := csvutil.Backup(path)
	if err != nil {
		return nil, err
	}
	return os.Open(bp)
}
