package main

import (
	"github.com/pinzolo/csvutil"
)

var cmdNumeric = &Command{
	Run:       runNumeric,
	UsageLine: "numeric [OPTIONS...] [FILE]",
	Short:     "数値生成",
	Long: `DESCRIPTION
        指定した列にmin <= n < max となるランダムな数値を出力します。

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

        -mx, --max NUMBER
            出力する数値の最大値を指定します。ただし、ここで指定した値は出力されません。

        -mn, --min NUMBER
            出力する数値の最小値を指定します。

        -d, --decimal
            出力する数値を小数として出力します。

        -dd, --decimal-digit NUMBER
            出力する小数の有効桁数を指定します。値は正の整数でなければいけません。（初期値: 3）
	`,
}

type cmdNumericOption struct {
	csvutil.NumericOption
	Overwrite bool
	Backup    bool
}

var numericOpt = cmdNumericOption{}

func init() {
	cmdNumeric.Flag.BoolVar(&numericOpt.Overwrite, "overwrite", false, "Overwrite to source.")
	cmdNumeric.Flag.BoolVar(&numericOpt.Overwrite, "w", false, "Overwrite to source.")
	cmdNumeric.Flag.BoolVar(&numericOpt.NoHeader, "no-header", false, "Source file does not have header line.")
	cmdNumeric.Flag.BoolVar(&numericOpt.NoHeader, "H", false, "Source file does not have header line.")
	cmdNumeric.Flag.BoolVar(&numericOpt.Backup, "backup", false, "Backup source file.")
	cmdNumeric.Flag.BoolVar(&numericOpt.Backup, "b", false, "Backup source file.")
	cmdNumeric.Flag.StringVar(&numericOpt.Encoding, "encoding", "utf8", "Encoding of source file")
	cmdNumeric.Flag.StringVar(&numericOpt.Encoding, "e", "utf8", "Encoding of source file")
	cmdNumeric.Flag.StringVar(&numericOpt.OutputEncoding, "output-encoding", "", "Encoding for output")
	cmdNumeric.Flag.StringVar(&numericOpt.OutputEncoding, "oe", "", "Encoding for output")
	cmdNumeric.Flag.StringVar(&numericOpt.Column, "column", "", "Target column symbol")
	cmdNumeric.Flag.StringVar(&numericOpt.Column, "c", "", "Target column symbol")
	cmdNumeric.Flag.IntVar(&numericOpt.Max, "max", 100, "Maximum value")
	cmdNumeric.Flag.IntVar(&numericOpt.Max, "mx", 100, "Maximum value")
	cmdNumeric.Flag.IntVar(&numericOpt.Min, "min", 0, "Minimum value")
	cmdNumeric.Flag.IntVar(&numericOpt.Min, "mn", 0, "Minmum value")
	cmdNumeric.Flag.BoolVar(&numericOpt.Decimal, "decimal", false, "Output decimal number")
	cmdNumeric.Flag.BoolVar(&numericOpt.Decimal, "d", false, "Output decimal number")
	cmdNumeric.Flag.IntVar(&numericOpt.DecimalDigit, "decimal-digit", 3, "Decimal digit number")
	cmdNumeric.Flag.IntVar(&numericOpt.DecimalDigit, "dd", 3, "Decimal digit number")
}

// runNumeric executes numeric command and return exit code.
func runNumeric(args []string) int {
	success := false
	w, wf, r, rf, err := prepare(args, numericOpt.Overwrite)
	if wf != nil {
		defer wf(&success, numericOpt.Backup)
	}
	if rf != nil {
		defer rf()
	}
	if err != nil {
		return handleError(err)
	}

	err = csvutil.Numeric(r, w, numericOpt.NumericOption)
	if err != nil {
		return handleError(err)
	}

	success = true
	return 0
}
