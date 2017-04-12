package main

import (
	"github.com/pinzolo/csvutil"
)

var cmdCombine = &Command{
	Run:       runCombine,
	UsageLine: "combine [OPTIONS...] [FILE]",
	Short:     "列結合",
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

        -s, --source COLUMN_SYMBOL(S)
            結合元の列のシンボルを指定します。
            列のシンボルとは列のインデックス（0開始）、もしくはヘッダーテキストです。
            --no-header オプションが指定された場合、インデックスしか受け入れません。
            複数列を対象としたい場合は、foo:bar や 1:2のようにコロン区切りで指定して下さい。

        -d, --destination COLUMN_SYMBOL
            結合後の値を入力する列のシンボルを指定します。
            列のシンボルとは列のインデックス（0開始）、もしくはヘッダーテキストです。
            --no-header オプションが指定された場合、インデックスしか受け入れません。

        -dl, --delimiter TEXT
            結合時にデリミタとする文字列を指定します。初期値は空文字です。
	`,
}

type cmdCombineOption struct {
	csvutil.CombineOption
	// Overwrite to source. (default false)
	Overwrite bool
	// Backup source file. (default false)
	Backup bool
	// Source header or column index separated by colon.
	Source string
}

var combineOpt = cmdCombineOption{}

func init() {
	cmdCombine.Flag.BoolVar(&combineOpt.Overwrite, "overwrite", false, "Overwrite to source.")
	cmdCombine.Flag.BoolVar(&combineOpt.Overwrite, "w", false, "Overwrite to source.")
	cmdCombine.Flag.BoolVar(&combineOpt.NoHeader, "no-header", false, "Source file does not have header line.")
	cmdCombine.Flag.BoolVar(&combineOpt.NoHeader, "H", false, "Source file does not have header line.")
	cmdCombine.Flag.BoolVar(&combineOpt.Backup, "backup", false, "Backup source file.")
	cmdCombine.Flag.BoolVar(&combineOpt.Backup, "b", false, "Backup source file.")
	cmdCombine.Flag.StringVar(&combineOpt.Encoding, "encoding", "utf8", "Encoding of source file")
	cmdCombine.Flag.StringVar(&combineOpt.Encoding, "e", "utf8", "Encoding of source file")
	cmdCombine.Flag.StringVar(&combineOpt.OutputEncoding, "output-encoding", "", "Encoding for output")
	cmdCombine.Flag.StringVar(&combineOpt.OutputEncoding, "oe", "", "Encoding for output")
	cmdCombine.Flag.StringVar(&combineOpt.Source, "source", "", "Source column symbol")
	cmdCombine.Flag.StringVar(&combineOpt.Source, "s", "", "Source column symbol")
	cmdCombine.Flag.StringVar(&combineOpt.Destination, "destination", "", "Destination column symbol")
	cmdCombine.Flag.StringVar(&combineOpt.Destination, "d", "", "Destination column symbol")
	cmdCombine.Flag.StringVar(&combineOpt.Delimiter, "delimiter", "", "Delimiter")
	cmdCombine.Flag.StringVar(&combineOpt.Delimiter, "dl", "", "Delimiter")
}

// runCombine executes combine command and return exit code.
func runCombine(args []string) int {
	success := false
	w, wf, r, rf, err := prepare(args, combineOpt.Overwrite)
	if wf != nil {
		defer wf(&success, combineOpt.Backup)
	}
	if rf != nil {
		defer rf()
	}
	if err != nil {
		return handleError(err)
	}

	opt := combineOpt.CombineOption
	opt.SourceSyms = split(combineOpt.Source)
	err = csvutil.Combine(r, w, opt)
	if err != nil {
		return handleError(err)
	}

	success = true
	return 0
}
