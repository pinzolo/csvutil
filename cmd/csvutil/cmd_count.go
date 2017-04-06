package main

import (
	"fmt"

	"github.com/pinzolo/csvutil"
)

var cmdCount = &Command{
	Run:       runCount,
	UsageLine: "count [OPTIONS...] [FILE]",
	Short:     "行数カウント",
	Long: `DESCRIPTION
        CSVの行数をカウントします。
        ヘッダの存在や、値に改行が含まれていることを考慮したカウントが行えます。

ARGUMENTS
        FILE
            ソースとなる CSV ファイルのパスを指定します。
            パスが指定されていない場合、標準入力が対象となりパイプでの使用ができます。

OPTIONS
        -H, --no-header
            ソースとなるCSVの1行目をヘッダー列として扱いません。

        -e, --encoding
            ソースとなるCSVの文字エンコーディングを指定します。
            このオプションが指定されていない場合、csvutil はUTF-8とみなして処理を行います。
            UTF-8であった場合、BOMのあるなしは自動的に判別されます。
            対応している値:
                sjis : Shift_JISとして扱います
                eucjp: EUC_JP として扱います
	`,
}

type cmdCountOption struct {
	csvutil.CountOption
}

var countOpt = cmdCountOption{}

func init() {
	cmdCount.Flag.BoolVar(&countOpt.NoHeader, "no-header", false, "Source file does not have header line.")
	cmdCount.Flag.BoolVar(&countOpt.NoHeader, "H", false, "Source file does not have header line.")
	cmdCount.Flag.StringVar(&countOpt.Encoding, "encoding", "utf8", "Encoding of source file")
	cmdCount.Flag.StringVar(&countOpt.Encoding, "e", "utf8", "Encoding of source file")
}

// runCount executes count command and return exit code.
func runCount(args []string) int {
	path, err := path(args)
	if err != nil {
		return handleError(err)
	}

	r, rf, err := reader(path)
	if err != nil {
		return handleError(err)
	}
	if rf != nil {
		defer rf()
	}

	i, err := csvutil.Count(r, countOpt.CountOption)
	if err != nil {
		return handleError(err)
	}
	fmt.Println(i)

	return 0
}
