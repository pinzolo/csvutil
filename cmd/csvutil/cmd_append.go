package main

import (
	"github.com/pinzolo/csvutil"
)

var cmdAppend = &Command{
	Run:       runAppend,
	UsageLine: "append [OPTIONS...] [FILE]",
	Short:     "Append empty values to CSV each line.",
	Long: `DESCRIPTION
        Append empty values to CSV each line.

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

        -oe, --output-encoding
            Encoding for output.
            This option accepts 'sjis', 'eucjp', 'utf8' or 'utf8bom'.
            Without this option, using --encoding option (or default).

        -h, --header
            Appending header text.
            To target multi headers, use semicolon separated value like foo:bar.
            If this option is not given, new header texts are set with column1, column2...

        -s, --size
            Appending column size. Default is 1
            If size is less than header length, ignore unused header(s).
            If size is greater than header length, append default header(s).
	`,
}

type cmdAppendOption struct {
	csvutil.AppendOption
	// Overwrite to source. (default false)
	Overwrite bool
	// Backup source file. (default false)
	Backup bool
	// Header symbols.
	Header string
}

var appendOpt = cmdAppendOption{}

func init() {
	cmdAppend.Flag.BoolVar(&appendOpt.Overwrite, "overwrite", false, "Overwrite to source.")
	cmdAppend.Flag.BoolVar(&appendOpt.Overwrite, "w", false, "Overwrite to source.")
	cmdAppend.Flag.BoolVar(&appendOpt.NoHeader, "no-header", false, "Source file does not have header line.")
	cmdAppend.Flag.BoolVar(&appendOpt.NoHeader, "H", false, "Source file does not have header line.")
	cmdAppend.Flag.BoolVar(&appendOpt.Backup, "backup", false, "Backup source file.")
	cmdAppend.Flag.BoolVar(&appendOpt.Backup, "b", false, "Backup source file.")
	cmdAppend.Flag.StringVar(&appendOpt.Encoding, "encoding", "utf8", "Encoding of source file")
	cmdAppend.Flag.StringVar(&appendOpt.Encoding, "e", "utf8", "Encoding of source file")
	cmdAppend.Flag.StringVar(&appendOpt.OutputEncoding, "output-encoding", "", "Encoding for output")
	cmdAppend.Flag.StringVar(&appendOpt.OutputEncoding, "oe", "", "Encoding for output")
	cmdAppend.Flag.StringVar(&appendOpt.Header, "header", "", "Appending header(s)")
	cmdAppend.Flag.StringVar(&appendOpt.Header, "h", "", "Appending header(s)")
	cmdAppend.Flag.IntVar(&appendOpt.Size, "size", 1, "Appending column size")
	cmdAppend.Flag.IntVar(&appendOpt.Size, "s", 1, "Appending column size")
}

// runAppend executes append command and return exit code.
func runAppend(args []string) int {
	success := false
	path, err := path(args)
	if err != nil {
		return handleError(err)
	}

	w, wf, err := writer(path, appendOpt.Overwrite)
	if err != nil {
		return handleError(err)
	}
	if wf != nil {
		defer wf(&success, appendOpt.Backup)
	}

	r, rf, err := reader(path)
	if err != nil {
		return handleError(err)
	}
	if rf != nil {
		defer rf()
	}

	opt := appendOpt.AppendOption
	opt.Headers = split(appendOpt.Header)
	err = csvutil.Append(r, w, opt)
	if err != nil {
		return handleError(err)
	}

	success = true
	return 0
}
