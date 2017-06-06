package main

import (
	"github.com/pinzolo/csvutil"
)

var cmdSort = &Command{
	Run:       runSort,
	UsageLine: "sort [OPTIONS...] [FILE]",
	Short:     "ソート",
	Long: `DESCRIPTION
        指定した列を基準にしてソートした CSV を出力します

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
            ソート対象となる列のシンボルを指定します。
            列のシンボルとは列のインデックス（0開始）、もしくはヘッダーテキストです。
            --no-header オプションが指定された場合、インデックスしか受け入れません。

        -d, --desc, --descending
            このオプションを指定するとソートを降順で行います。

        -dt, --data-type TYPE
            ソート対象列の値をどのように扱うかを指定します。
            対応している値:
                text   : 文字列としてソートします（初期値）
                number : 数値としてソートします

        -em, --empty HANDLING
            値が空の場合にどのように処理するかを指定します。
            対応している値:
                natural : システムに任せます（初期値）
                          --data-type が number の場合、空値は 0 として扱われます
                first   : 強制的に先頭に移動します
                last    : 強制的に末尾に移動します
	`,
}

type cmdSortOption struct {
	csvutil.SortOption
	Overwrite bool
	Backup    bool
}

var sortOpt = cmdSortOption{}

func init() {
	cmdSort.Flag.BoolVar(&sortOpt.Overwrite, "overwrite", false, "Overwrite to source.")
	cmdSort.Flag.BoolVar(&sortOpt.Overwrite, "w", false, "Overwrite to source.")
	cmdSort.Flag.BoolVar(&sortOpt.NoHeader, "no-header", false, "Source file does not have header line.")
	cmdSort.Flag.BoolVar(&sortOpt.NoHeader, "H", false, "Source file does not have header line.")
	cmdSort.Flag.BoolVar(&sortOpt.Backup, "backup", false, "Backup source file.")
	cmdSort.Flag.BoolVar(&sortOpt.Backup, "b", false, "Backup source file.")
	cmdSort.Flag.StringVar(&sortOpt.Encoding, "encoding", "utf8", "Encoding of source file")
	cmdSort.Flag.StringVar(&sortOpt.Encoding, "e", "utf8", "Encoding of source file")
	cmdSort.Flag.StringVar(&sortOpt.OutputEncoding, "output-encoding", "", "Encoding for output")
	cmdSort.Flag.StringVar(&sortOpt.OutputEncoding, "oe", "", "Encoding for output")
	cmdSort.Flag.StringVar(&sortOpt.Column, "column", "", "Home column symbol")
	cmdSort.Flag.StringVar(&sortOpt.Column, "c", "", "Home column symbol")
	cmdSort.Flag.StringVar(&sortOpt.DataType, "data-type", csvutil.SortDataTypeText, "Data type")
	cmdSort.Flag.StringVar(&sortOpt.DataType, "dt", csvutil.SortDataTypeText, "Data type")
	cmdSort.Flag.BoolVar(&sortOpt.Descending, "descending", false, "Order in descending")
	cmdSort.Flag.BoolVar(&sortOpt.Descending, "desc", false, "Order in descending")
	cmdSort.Flag.BoolVar(&sortOpt.Descending, "d", false, "Order in descending")
	cmdSort.Flag.StringVar(&sortOpt.EmptyHandling, "empty", csvutil.EmptyNatural, "Empty handling method")
	cmdSort.Flag.StringVar(&sortOpt.EmptyHandling, "em", csvutil.EmptyNatural, "Empty handling method")
}

// runSort executes sort command and return exit code.
func runSort(args []string) int {
	success := false
	w, wf, r, rf, err := prepare(args, sortOpt.Overwrite)
	if wf != nil {
		defer wf(&success, sortOpt.Backup)
	}
	if rf != nil {
		defer rf()
	}
	if err != nil {
		return handleError(err)
	}

	err = csvutil.Sort(r, w, sortOpt.SortOption)
	if err != nil {
		return handleError(err)
	}

	success = true
	return 0
}
