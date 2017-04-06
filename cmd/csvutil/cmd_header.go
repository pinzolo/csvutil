package main

import (
	"github.com/pinzolo/csvutil"
)

var cmdHeader = &Command{
	Run:       runHeader,
	UsageLine: "header [OPTIONS...] [FILE]",
	Short:     "ヘッダー読取",
	Long: `DESCRIPTION
        CSVのヘッダーだけを読込、1列1行として出力します。

ARGUMENTS
        FILE
            ソースとなる CSV ファイルのパスを指定します。
            パスが指定されていない場合、標準入力が対象となりパイプでの使用ができます。

OPTIONS
        -e, --encoding
            ソースとなるCSVの文字エンコーディングを指定します。
            このオプションが指定されていない場合、csvutil はUTF-8とみなして処理を行います。
            UTF-8であった場合、BOMのあるなしは自動的に判別されます。
            対応している値:
                sjis : Shift_JISとして扱います
                eucjp: EUC_JPとして扱います

        -oe, --output-encoding
            出力するCSVの文字エンコーディングを指定します。
            このオプションが指定されていない場合 --encoding オプションで指定されたエンコーディングとして出力します。
            対応している値:
                utf8    : UTF-8として出力します（BOMは出力しません）
                utf8bom : UTF-8として出力します（BOMは出力します）
                sjis    : Shift_JISとして出力します
                eucjp   : EUC_JPとして出力します

        -i, --index
            このオプションを指定すると、列のインデックスも合わせて出力します。

        -io, --index-origin
            インデックスの開始値を指定します。初期値は 0 です。
            --index オプションが指定されていない場合、このオプションは無視されます。
	`,
}

type cmdHeaderOption struct {
	csvutil.HeaderOption
}

var headerOpt = cmdHeaderOption{}

func init() {
	cmdHeader.Flag.StringVar(&headerOpt.Encoding, "encoding", "utf8", "Encoding of source file")
	cmdHeader.Flag.StringVar(&headerOpt.Encoding, "e", "utf8", "Encoding of source file")
	cmdHeader.Flag.StringVar(&headerOpt.OutputEncoding, "output-encoding", "", "Encoding for output")
	cmdHeader.Flag.StringVar(&headerOpt.OutputEncoding, "oe", "", "Encoding for output")
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
