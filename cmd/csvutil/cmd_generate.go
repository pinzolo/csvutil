package main

import (
	"github.com/pinzolo/csvutil"
)

var cmdGenerate = &Command{
	Run:       runGenerate,
	UsageLine: "generate [OPTIONS...]",
	Short:     "Generate CSV that have no values.",
	Long: `DESCRIPTION
        Generate CSV that have no values.

OPTIONS
        -H, --no-header
            Tel given CSV does not have header line.

        -e, --encoding
            Encoding of source file.
            This option accepts 'sjis' or 'eucjp'.
            Without this option, csvutil treats CSV file is encoded by UTF-8.

        -h, --header
            Generateing header text.
            To target multi headers, use semicolon separated value like foo:bar.
            If this option is not given, new header texts are set with column1, column2...

        -s, --size
            Generateing column size. Default size is 3.
            If size is less than header length, ignore unused header(s).
            If size is greater than header length, append default header(s).

        -c, --count
            Generateing line count. Default count is 3.

        -b, --bom
            Append bom to top of CSV.
            If encoding is not UTF-8, this option is ignored.
	`,
}

type cmdGenerateOption struct {
	csvutil.GenerateOption
	// Header symbols.
	Header string
}

var generateOpt = cmdGenerateOption{}

func init() {
	cmdGenerate.Flag.BoolVar(&generateOpt.NoHeader, "no-header", false, "Source file does not have header line")
	cmdGenerate.Flag.BoolVar(&generateOpt.NoHeader, "H", false, "Source file does not have header line")
	cmdGenerate.Flag.StringVar(&generateOpt.Encoding, "encoding", "utf8", "Encoding of source file")
	cmdGenerate.Flag.StringVar(&generateOpt.Encoding, "e", "utf8", "Encoding of source file")
	cmdGenerate.Flag.StringVar(&generateOpt.Header, "header", "", "Generateing header(s)")
	cmdGenerate.Flag.StringVar(&generateOpt.Header, "h", "", "Generateing header(s)")
	cmdGenerate.Flag.IntVar(&generateOpt.Size, "size", 3, "Generateing column size")
	cmdGenerate.Flag.IntVar(&generateOpt.Size, "s", 3, "Generateing column size")
	cmdGenerate.Flag.IntVar(&generateOpt.Count, "count", 3, "Generateing line count")
	cmdGenerate.Flag.IntVar(&generateOpt.Count, "c", 3, "Generateing line count")
	cmdGenerate.Flag.BoolVar(&generateOpt.BOM, "bom", false, "Add bom")
	cmdGenerate.Flag.BoolVar(&generateOpt.BOM, "b", false, "Add bom")
}

// runGenerate executes generate command and return exit code.
func runGenerate(args []string) int {
	success := false
	w, wf, err := writer("", false)
	if err != nil {
		return handleError(err)
	}
	if wf != nil {
		defer wf(&success, false)
	}

	opt := generateOpt.GenerateOption
	opt.Headers = split(generateOpt.Header)
	err = csvutil.Generate(w, opt)
	if err != nil {
		return handleError(err)
	}

	success = true
	return 0
}
