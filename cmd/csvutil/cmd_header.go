package main

import (
	"github.com/pinzolo/csvutil"
)

var cmdHeader = &Command{
	Run:       runHeader,
	UsageLine: "header [OPTIONS...] [FILE]",
	Short:     "Header CSV that have no values.",
	Long: `DESCRIPTION
        Header CSV that have no values.

ARGUMENTS
        FILE
            Source CSV file.
            Without FILE argument, read from STDIN.

OPTIONS
        -e, --encoding
            Encoding of source file.
            This option accepts 'sjis' or 'eucjp'.
            Without this option, csvutil treats CSV file is encoded by UTF-8.

        -i, --index
            Print header with index.

        -io, --index-origin
            Start number of index.
            If --index option is not given, this option is ignored.
	`,
}

type cmdHeaderOption struct {
	csvutil.HeaderOption
}

var headerOpt = cmdHeaderOption{}

func init() {
	cmdHeader.Flag.StringVar(&headerOpt.Encoding, "encoding", "utf8", "Encoding of source file")
	cmdHeader.Flag.StringVar(&headerOpt.Encoding, "e", "utf8", "Encoding of source file")
	cmdHeader.Flag.BoolVar(&headerOpt.Index, "index", false, "Print index")
	cmdHeader.Flag.BoolVar(&headerOpt.Index, "i", false, "Print index")
	cmdHeader.Flag.IntVar(&headerOpt.IndexOrigin, "index-origin", 0, "Index origin number")
	cmdHeader.Flag.IntVar(&headerOpt.IndexOrigin, "io", 0, "Index origin number")
}

// runHeader executes header command and return exit code.
func runHeader(args []string) int {
	path, err := path(args)
	if err != nil {
		return handleError(err)
	}

	w, wf, err := writer(path, false)
	if err != nil {
		return handleError(err)
	}
	if wf != nil {
		f := false
		defer wf(&f, false)
	}

	r, rf, err := reader(path)
	if err != nil {
		return handleError(err)
	}
	if rf != nil {
		defer rf()
	}

	err = csvutil.Header(r, w, headerOpt.HeaderOption)
	if err != nil {
		return handleError(err)
	}

	return 0
}
