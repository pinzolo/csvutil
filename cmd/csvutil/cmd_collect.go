package main

import (
	"github.com/pinzolo/csvutil"
)

var cmdCollect = &Command{
	Run:       runCollect,
	UsageLine: "collect [OPTIONS...] [FILE]",
	Short:     "値抽出",
	Long: `DESCRIPTION
        指定した列の値を収集し、重複した値を除いて出力します。

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

        -c, --column COLUMN_SYMBOL
            対象となる列のシンボルを指定します。
            列のシンボルとは列のインデックス（0開始）、もしくはヘッダーテキストです。
            --no-header オプションが指定された場合、インデックスしか受け入れません。

        -ae, --allow-empty
            このオプションを指定すると、値が空でも値として扱います。

        -pc, --print-count
            このオプションを指定すると、値が出現する数を数えて出力します。

        -s, --sort
            このオプションを指定すると、収集した値をソートして出力します。

        -sk, --sort-key KEY
            ソートする際に、値と出現数のどちらをキーにしてソートするかを指定します。
            対応している値:
                value : 値をキーにしてソートします（初期値）
                count : 出現数をキーにしてソートします

        -d, --descending
            ソートする際に降順で並び替えます。
	`,
}

type cmdCollectOption struct {
	csvutil.CollectOption
	Overwrite bool
	Backup    bool
}

var collectOpt = cmdCollectOption{}

func init() {
	cmdCollect.Flag.BoolVar(&collectOpt.Overwrite, "overwrite", false, "Overwrite to source.")
	cmdCollect.Flag.BoolVar(&collectOpt.Overwrite, "w", false, "Overwrite to source.")
	cmdCollect.Flag.BoolVar(&collectOpt.NoHeader, "no-header", false, "Source file does not have header line.")
	cmdCollect.Flag.BoolVar(&collectOpt.NoHeader, "H", false, "Source file does not have header line.")
	cmdCollect.Flag.BoolVar(&collectOpt.Backup, "backup", false, "Backup source file.")
	cmdCollect.Flag.BoolVar(&collectOpt.Backup, "b", false, "Backup source file.")
	cmdCollect.Flag.StringVar(&collectOpt.Encoding, "encoding", "utf8", "Encoding of source file")
	cmdCollect.Flag.StringVar(&collectOpt.Encoding, "e", "utf8", "Encoding of source file")
	cmdCollect.Flag.StringVar(&collectOpt.OutputEncoding, "output-encoding", "", "Encoding for output")
	cmdCollect.Flag.StringVar(&collectOpt.OutputEncoding, "oe", "", "Encoding for output")
	cmdCollect.Flag.StringVar(&collectOpt.Column, "column", "", "Target column symbol")
	cmdCollect.Flag.StringVar(&collectOpt.Column, "c", "", "Home column symbol")
	cmdCollect.Flag.BoolVar(&collectOpt.AllowEmpty, "allow-empty", false, "Allow empty")
	cmdCollect.Flag.BoolVar(&collectOpt.AllowEmpty, "ae", false, "Allow empty")
	cmdCollect.Flag.BoolVar(&collectOpt.PrintCount, "print-count", false, "Print count")
	cmdCollect.Flag.BoolVar(&collectOpt.PrintCount, "pc", false, "Print count")
	cmdCollect.Flag.BoolVar(&collectOpt.Sort, "sort", false, "Sort")
	cmdCollect.Flag.BoolVar(&collectOpt.Sort, "s", false, "Sort")
	cmdCollect.Flag.StringVar(&collectOpt.SortKey, "sort-key", "value", "Sort key")
	cmdCollect.Flag.StringVar(&collectOpt.SortKey, "sk", "value", "Sort key")
	cmdCollect.Flag.BoolVar(&collectOpt.Descending, "descending", false, "Sort in descending order")
	cmdCollect.Flag.BoolVar(&collectOpt.Descending, "desc", false, "Sort in descending order")
	cmdCollect.Flag.BoolVar(&collectOpt.Descending, "d", false, "Sort in descending order")
}

// runCollect executes collect command and return exit code.
func runCollect(args []string) int {
	success := false
	w, wf, r, rf, err := prepare(args, collectOpt.Overwrite)
	if wf != nil {
		defer wf(&success, collectOpt.Backup)
	}
	if rf != nil {
		defer rf()
	}
	if err != nil {
		return handleError(err)
	}

	err = csvutil.Collect(r, w, collectOpt.CollectOption)
	if err != nil {
		return handleError(err)
	}

	success = true
	return 0
}
