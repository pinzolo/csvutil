package main

import (
	"github.com/pinzolo/csvutil"
)

var cmdRemove = &Command{
	Run:       runRemove,
	UsageLine: "remove [OPTIONS...] [FILE]",
	Short:     "列削除",
	Long: `DESCRIPTION
        指定された列を削除したCSVを出力します。

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
            削除する列のシンボルを指定します。
            列のシンボルとは列のインデックス（0開始）、もしくはヘッダーテキストです。
            --no-header オプションが指定された場合、インデックスしか受け入れません。
            複数列を対象としたい場合は、foo:bar や 1:2のようにコロン区切りで指定して下さい。
	`,
}

type cmdRemoveOption struct {
	csvutil.RemoveOption
	// Overwrite to source. (default false)
	Overwrite bool
	// Backup source file. (default false)
	Backup bool
	// Column header or column index separated by semicolon.
	Column string
}

var removeOpt = cmdRemoveOption{}

func init() {
	cmdRemove.Flag.BoolVar(&removeOpt.Overwrite, "overwrite", false, "Overwrite to source.")
	cmdRemove.Flag.BoolVar(&removeOpt.Overwrite, "w", false, "Overwrite to source.")
	cmdRemove.Flag.BoolVar(&removeOpt.NoHeader, "no-header", false, "Source file does not have header line.")
	cmdRemove.Flag.BoolVar(&removeOpt.NoHeader, "H", false, "Source file does not have header line.")
	cmdRemove.Flag.BoolVar(&removeOpt.Backup, "backup", false, "Backup source file.")
	cmdRemove.Flag.BoolVar(&removeOpt.Backup, "b", false, "Backup source file.")
	cmdRemove.Flag.StringVar(&removeOpt.Encoding, "encoding", "utf8", "Encoding of source file")
	cmdRemove.Flag.StringVar(&removeOpt.Encoding, "e", "utf8", "Encoding of source file")
	cmdRemove.Flag.StringVar(&removeOpt.OutputEncoding, "output-encoding", "", "Encoding for output")
	cmdRemove.Flag.StringVar(&removeOpt.OutputEncoding, "oe", "", "Encoding for output")
	cmdRemove.Flag.StringVar(&removeOpt.Column, "column", "", "Column symbol")
	cmdRemove.Flag.StringVar(&removeOpt.Column, "c", "", "Column symbol")
}

// runRemove executes remove command and return exit code.
func runRemove(args []string) int {
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

	opt := removeOpt.RemoveOption
	opt.ColumnSyms = split(removeOpt.Column)
	err = csvutil.Remove(r, w, opt)
	if err != nil {
		return handleError(err)
	}

	success = true
	return 0
}
