package main

import (
	"github.com/pinzolo/csvutil"
)

var cmdBlank = &Command{
	Run:       runBlank,
	UsageLine: "blank [OPTIONS...] [FILE]",
	Short:     "列の空白可",
	Long: `DESCRIPTION
        指定した列の値を削除、もしくは特定の空白文字で置き換えます。

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

        -c, --column
            対象となる列のシンボルを指定します。
            列のシンボルとは列のインデックス（0開始）、もしくはヘッダーテキストです。
            --no-header オプションが指定された場合、インデックスしか受け入れません。
            複数列を対象としたい場合は、foo:bar や 1:2のようにコロン区切りで指定して下さい。

        -r, --rate
            空白可する割合を指定します。0〜100までの整数を指定して下さい。
            指定しない場合、100%空白化します。

        -sw, --space-width
            空白化に指定する空白文字を数字で指定します。
                0: 空文字（初期値）
                1: 半角スペース [0x20]
                2: 全角スペース [0xE3 0x80 0x80]

        --ss, --space-size
            空白化する際に --space-width で指定した文字を何回繰り返すかを指定します。
            初期値が0なので --space-width も指定しないと必ず空文字での空白可を行います。
            例:
                --space-size が 0 で --space-width が 1 ならば空文字
                --space-size が 1 で --space-width が 0 ならば空文字
                --space-size が 2 で --space-width が 1 ならば "  " (半角スペース2つ）
                --space-size が 3 で --space-width が 2 ならば "　　　" (全角スペース3つ）
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
	cmdBlank.Flag.StringVar(&blankOpt.OutputEncoding, "output-encoding", "", "Encoding for output")
	cmdBlank.Flag.StringVar(&blankOpt.OutputEncoding, "oe", "", "Encoding for output")
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
	w, wf, r, rf, err := prepare(args, blankOpt.Overwrite)
	if wf != nil {
		defer wf(&success, blankOpt.Backup)
	}
	if rf != nil {
		defer rf()
	}
	if err != nil {
		return handleError(err)
	}

	opt := blankOpt.BlankOption
	opt.ColumnSyms = split(blankOpt.Column)
	err = csvutil.Blank(r, w, opt)
	if err != nil {
		return handleError(err)
	}

	success = true
	return 0
}
