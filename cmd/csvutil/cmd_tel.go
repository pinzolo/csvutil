package main

import (
	"github.com/pinzolo/csvutil"
)

var cmdTel = &Command{
	Run:       runTel,
	UsageLine: "tel [OPTIONS...] [FILE]",
	Short:     "電話番号埋め込み",
	Long: `DESCRIPTION
        指定した列にダミーの電話番号を埋め込みます。

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

        -mr, --mobile-rate
            ダミーの携帯番号を設定する割合を指定します。0〜100までの整数を指定して下さい。（初期値: 0）
            csvutilが埋め込む携帯番号は090,080,070,050で始まるランダムな電話番号です。
	`,
}

type cmdTelOption struct {
	csvutil.TelOption
	Overwrite bool
	Backup    bool
}

var telOpt = cmdTelOption{}

func init() {
	cmdTel.Flag.BoolVar(&telOpt.Overwrite, "overwrite", false, "Overwrite to source.")
	cmdTel.Flag.BoolVar(&telOpt.Overwrite, "w", false, "Overwrite to source.")
	cmdTel.Flag.BoolVar(&telOpt.NoHeader, "no-header", false, "Source file does not have header line.")
	cmdTel.Flag.BoolVar(&telOpt.NoHeader, "H", false, "Source file does not have header line.")
	cmdTel.Flag.BoolVar(&telOpt.Backup, "backup", false, "Backup source file.")
	cmdTel.Flag.BoolVar(&telOpt.Backup, "b", false, "Backup source file.")
	cmdTel.Flag.StringVar(&telOpt.Encoding, "encoding", "utf8", "Encoding of source file")
	cmdTel.Flag.StringVar(&telOpt.Encoding, "e", "utf8", "Encoding of source file")
	cmdTel.Flag.StringVar(&telOpt.OutputEncoding, "output-encoding", "", "Encoding for output")
	cmdTel.Flag.StringVar(&telOpt.OutputEncoding, "oe", "", "Encoding for output")
	cmdTel.Flag.StringVar(&telOpt.Column, "column", "", "Home column symbol")
	cmdTel.Flag.StringVar(&telOpt.Column, "c", "", "Home column symbol")
	cmdTel.Flag.IntVar(&telOpt.MobileRate, "mobile-rate", 0, "Mobile tel number rate")
	cmdTel.Flag.IntVar(&telOpt.MobileRate, "mr", 0, "Mobile tel number rate")
}

// runTel executes tel command and return exit code.
func runTel(args []string) int {
	success := false
	w, wf, r, rf, err := prepare(args, insertOpt.Overwrite)
	if wf != nil {
		defer wf(&success, insertOpt.Backup)
	}
	if rf != nil {
		defer rf()
	}
	if err != nil {
		return handleError(err)
	}

	err = csvutil.Tel(r, w, telOpt.TelOption)
	if err != nil {
		return handleError(err)
	}

	success = true
	return 0
}
