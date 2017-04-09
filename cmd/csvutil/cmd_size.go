package main

import (
	"fmt"

	"github.com/pinzolo/csvutil"
)

var cmdSize = &Command{
	Run:       runSize,
	UsageLine: "size [OPTIONS...] [FILE]",
	Short:     "列数カウント",
	Long: `DESCRIPTION
        CSVの列数をカウントします。

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
                sjis : Shift_JIS として扱います
                eucjp: EUC_JPとして扱います
	`,
}

type cmdSizeOption struct {
	csvutil.SizeOption
}

var sizeOpt = cmdSizeOption{}

func init() {
	cmdSize.Flag.StringVar(&sizeOpt.Encoding, "encoding", "utf8", "Encoding of source file")
	cmdSize.Flag.StringVar(&sizeOpt.Encoding, "e", "utf8", "Encoding of source file")
}

// runSize executes size command and return exit code.
func runSize(args []string) int {
	r, rf, err := prepareReader(args)
	if rf != nil {
		defer rf()
	}
	if err != nil {
		return handleError(err)
	}

	i, err := csvutil.Size(r, sizeOpt.SizeOption)
	if err != nil {
		return handleError(err)
	}
	fmt.Println(i)

	return 0
}
