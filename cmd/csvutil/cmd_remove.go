package main

import (
	"github.com/pinzolo/csvutil"
)

var cmdRemove = &Command{
	Run:       runRemove,
	UsageLine: "remove [OPTIONS...] [FILE]",
	Short:     "Remove column(s) in CSV.",
	Long: `DESCRIPTION
        Remove column(s) in CSV.

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
            To target multi columns, use semicolon separated value like foo:bar and 1:2.
	`,
}

type cmdRemoveOption struct {
	csvutil.RemoveOption
	// Overwrite to source. (default false)
	Overwrite bool
	// Backup source file. (default false)
	Backup bool
	// Column header or column index separated by semicolon.
	Column string
}

var removeOpt = cmdRemoveOption{}

func init() {
	cmdRemove.Flag.BoolVar(&removeOpt.Overwrite, "overwrite", false, "Overwrite to source.")
	cmdRemove.Flag.BoolVar(&removeOpt.Overwrite, "w", false, "Overwrite to source.")
	cmdRemove.Flag.BoolVar(&removeOpt.NoHeader, "no-header", false, "Source file does not have header line.")
	cmdRemove.Flag.BoolVar(&removeOpt.NoHeader, "H", false, "Source file does not have header line.")
	cmdRemove.Flag.BoolVar(&removeOpt.Backup, "backup", false, "Backup source file.")
	cmdRemove.Flag.BoolVar(&removeOpt.Backup, "b", false, "Backup source file.")
	cmdRemove.Flag.StringVar(&removeOpt.Encoding, "encoding", "utf8", "Encoding of source file")
	cmdRemove.Flag.StringVar(&removeOpt.Encoding, "e", "utf8", "Encoding of source file")
	cmdRemove.Flag.StringVar(&removeOpt.Column, "column", "", "Column symbol")
	cmdRemove.Flag.StringVar(&removeOpt.Column, "c", "", "Column symbol")
}

// runRemove executes remove command and return exit code.
func runRemove(args []string) int {
	success := false
	path, err := path(args)
	if err != nil {
		return handleError(err)
	}

	w, wf, err := writer(path, removeOpt.Overwrite)
	if err != nil {
		return handleError(err)
	}
	if wf != nil {
		defer wf(&success, removeOpt.Backup)
	}

	r, rf, err := reader(path)
	if err != nil {
		return handleError(err)
	}
	if rf != nil {
		defer rf()
	}

	opt := removeOpt.RemoveOption
	opt.ColumnSyms = split(removeOpt.Column)
	err = csvutil.Remove(r, w, removeOpt.RemoveOption)
	if err != nil {
		return handleError(err)
	}

	success = true
	return 0
}
