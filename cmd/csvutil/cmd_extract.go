package main

import (
	"github.com/pinzolo/csvutil"
)

var cmdExtract = &Command{
	Run:       runExtract,
	UsageLine: "extract [OPTIONS...] [FILE]",
	Short:     "列抽出",
	Long: `DESCRIPTION
        指定された列を抽出したCSVを出力します。

ARGUMENTS
        FILE
            ソースとなる CSV ファイルのパスを指定します。
            パスが指定されていない場合、標準入力が対象となりパイプでの使用ができます。

OPTIONS
        -w, --overwrite
            指定されたCSVファイルを実行結果で上書きします。
            ファイルパスが渡されていない場合には無視されます。

        -H, --no-header
            ソースとなるCSVの1行目をヘッダー列として扱いません。

        -b, --backup
            処理が成功した場合に、指定されたCSVファイルをバックアップします。
            --overwrite オプションと同時に使用されることを想定しているため、ファイルパスが渡されていない場合には無視されます。

        -e, --encoding ENCODING
            ソースとなるCSVの文字エンコーディングを指定します。
            このオプションが指定されていない場合、csvutil はUTF-8とみなして処理を行います。
            UTF-8であった場合、BOMのあるなしは自動的に判別されます。
            対応している値:
                sjis : Shift_JISとして扱います
                eucjp: EUC_JPとして扱います

        -oe, --output-encoding ENCODING
            出力するCSVの文字エンコーディングを指定します。
            このオプションが指定されていない場合 --encoding オプションで指定されたエンコーディングとして出力します。
            対応している値:
                utf8    : UTF-8として出力します（BOMは出力しません）
                utf8bom : UTF-8として出力します（BOMは出力します）
                sjis    : Shift_JISとして出力します
                eucjp   : EUC_JPとして出力します

        -c, --column COLUMN_SYMBOL(S)
            抽出する列のシンボルを指定します。
            列のシンボルとは列のインデックス（0開始）、もしくはヘッダーテキストです。
            --no-header オプションが指定された場合、インデックスしか受け入れません。
            複数列を対象としたい場合は、foo:bar や 1:2のようにコロン区切りで指定して下さい。
	`,
}

type cmdExtractOption struct {
	csvutil.ExtractOption
	// Overwrite to source. (default false)
	Overwrite bool
	// Backup source file. (default false)
	Backup bool
	// Column header or column index separated by semicolon.
	Column string
}

var extractOpt = cmdExtractOption{}

func init() {
	cmdExtract.Flag.BoolVar(&extractOpt.Overwrite, "overwrite", false, "Overwrite to source.")
	cmdExtract.Flag.BoolVar(&extractOpt.Overwrite, "w", false, "Overwrite to source.")
	cmdExtract.Flag.BoolVar(&extractOpt.NoHeader, "no-header", false, "Source file does not have header line.")
	cmdExtract.Flag.BoolVar(&extractOpt.NoHeader, "H", false, "Source file does not have header line.")
	cmdExtract.Flag.BoolVar(&extractOpt.Backup, "backup", false, "Backup source file.")
	cmdExtract.Flag.BoolVar(&extractOpt.Backup, "b", false, "Backup source file.")
	cmdExtract.Flag.StringVar(&extractOpt.Encoding, "encoding", "utf8", "Encoding of source file")
	cmdExtract.Flag.StringVar(&extractOpt.Encoding, "e", "utf8", "Encoding of source file")
	cmdExtract.Flag.StringVar(&extractOpt.OutputEncoding, "output-encoding", "", "Encoding for output")
	cmdExtract.Flag.StringVar(&extractOpt.OutputEncoding, "oe", "", "Encoding for output")
	cmdExtract.Flag.StringVar(&extractOpt.Column, "column", "", "Column symbol")
	cmdExtract.Flag.StringVar(&extractOpt.Column, "c", "", "Column symbol")
}

// runExtract executes extract command and return exit code.
func runExtract(args []string) int {
	success := false
	w, wf, r, rf, err := prepare(args, extractOpt.Overwrite)
	if wf != nil {
		defer wf(&success, extractOpt.Backup)
	}
	if rf != nil {
		defer rf()
	}
	if err != nil {
		return handleError(err)
	}

	opt := extractOpt.ExtractOption
	opt.ColumnSyms = split(extractOpt.Column)
	err = csvutil.Extract(r, w, opt)
	if err != nil {
		return handleError(err)
	}

	success = true
	return 0
}
