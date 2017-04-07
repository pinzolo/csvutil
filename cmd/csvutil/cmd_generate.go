package main

import (
	"github.com/pinzolo/csvutil"
)

var cmdGenerate = &Command{
	Run:       runGenerate,
	UsageLine: "generate [OPTIONS...]",
	Short:     "CSVの新規生成",
	Long: `DESCRIPTION
        全ての値が空のCSVを新規に出力します。

OPTIONS
        -H, --no-header
            このオプションを指定すると、ヘッダーの無いCSVを出力します。

        -oe, --output-encoding
            出力するCSVの文字エンコーディングを指定します。
            このオプションが指定されていない場合、UTF-8とみなして処理を行います。
            対応している値:
                utf8bom : UTF-8として出力します（BOMは出力します）
                sjis    : Shift_JISとして出力します
                eucjp   : EUC_JPとして出力します

        -h, --header
            新規に追加する列のヘッダーテキストを指定します。
            複数のヘッダーテキストを指定する場合には、foo:bar のようにコロン区切りにします。
            このオプションが指定されていない場合、もしくは指定したヘッダーが --size オプションの値に足りない場合には、
            column1,column2... のように連番のヘッダーが自動で付与されます。
            また --size オプションの値を超えた場合は、超えた分が無視されます。

        -s, --size
            生成する列の数を指定します。初期値は 3 です。

        -c, --count
            生成する行の数を指定します。初期値は 3 です。
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
	cmdGenerate.Flag.StringVar(&generateOpt.OutputEncoding, "output-encoding", "utf8", "Encoding for output")
	cmdGenerate.Flag.StringVar(&generateOpt.OutputEncoding, "oe", "utf8", "Encoding for output")
	cmdGenerate.Flag.StringVar(&generateOpt.Header, "header", "", "Generateing header(s)")
	cmdGenerate.Flag.StringVar(&generateOpt.Header, "h", "", "Generateing header(s)")
	cmdGenerate.Flag.IntVar(&generateOpt.Size, "size", 3, "Generateing column size")
	cmdGenerate.Flag.IntVar(&generateOpt.Size, "s", 3, "Generateing column size")
	cmdGenerate.Flag.IntVar(&generateOpt.Count, "count", 3, "Generateing line count")
	cmdGenerate.Flag.IntVar(&generateOpt.Count, "c", 3, "Generateing line count")
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
