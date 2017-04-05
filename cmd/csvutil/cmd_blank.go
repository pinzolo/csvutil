package main

import (
	"github.com/pinzolo/csvutil"
)

var cmdBlank = &Command{
	Run:       runBlank,
	UsageLine: "blank [OPTIONS...] [FILE]",
	Short:     "Replace column value(s) by empty or blank string.",
	Long: `DESCRIPTION
        Replace column value(s) by empty or blank string.

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

type cmdBlankOption struct {
	csvutil.BlankOption
	// Overwrite to source. (default false)
	Overwrite bool
	// Backup source file. (default false)
	Backup bool
	// Column header or column index separated by semicolon.
	Column string
}

var blankOpt = cmdBlankOption{}

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
	success := false
	path, err := path(args)
	if err != nil {
		return handleError(err)
	}

	w, wf, err := writer(path, blankOpt.Overwrite)
	if err != nil {
		return handleError(err)
	}
	if wf != nil {
		defer wf(&success, blankOpt.Backup)
	}

	r, rf, err := reader(path)
	if err != nil {
		return handleError(err)
	}
	if rf != nil {
		defer rf()
	}

	opt := blankOpt.BlankOption
	opt.ColumnSyms = split(blankOpt.Column)
	err = csvutil.Blank(r, w, blankOpt.BlankOption)
	if err != nil {
		return handleError(err)
	}

	success = true
	return 0
}
