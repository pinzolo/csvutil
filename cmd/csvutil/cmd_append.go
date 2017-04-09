package main

import (
	"github.com/pinzolo/csvutil"
)

var cmdAppend = &Command{
	Run:       runAppend,
	UsageLine: "append [OPTIONS...] [FILE]",
	Short:     "列追加",
	Long: `DESCRIPTION
        指定されたCSVの末尾に新規に列を追加したCSVを出力します。
        追加された列の値は空です。

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

        -h, --header
            新規に追加する列のヘッダーテキストを指定します。
            複数のヘッダーテキストを指定する場合には、foo:bar のようにコロン区切りにします。
            このオプションが指定されていない場合、もしくは指定したヘッダーが --size オプションの値に足りない場合には、
            column1,column2... のように連番のヘッダーが自動で付与されます。
            また --size オプションの値を超えた場合は、超えた分が無視されます。

        -s, --size
            追加する列の数を指定します。初期値は 1 です。
	`,
}

type cmdAppendOption struct {
	csvutil.AppendOption
	// Overwrite to source. (default false)
	Overwrite bool
	// Backup source file. (default false)
	Backup bool
	// Header symbols.
	Header string
}

var appendOpt = cmdAppendOption{}

func init() {
	cmdAppend.Flag.BoolVar(&appendOpt.Overwrite, "overwrite", false, "Overwrite to source.")
	cmdAppend.Flag.BoolVar(&appendOpt.Overwrite, "w", false, "Overwrite to source.")
	cmdAppend.Flag.BoolVar(&appendOpt.NoHeader, "no-header", false, "Source file does not have header line.")
	cmdAppend.Flag.BoolVar(&appendOpt.NoHeader, "H", false, "Source file does not have header line.")
	cmdAppend.Flag.BoolVar(&appendOpt.Backup, "backup", false, "Backup source file.")
	cmdAppend.Flag.BoolVar(&appendOpt.Backup, "b", false, "Backup source file.")
	cmdAppend.Flag.StringVar(&appendOpt.Encoding, "encoding", "utf8", "Encoding of source file")
	cmdAppend.Flag.StringVar(&appendOpt.Encoding, "e", "utf8", "Encoding of source file")
	cmdAppend.Flag.StringVar(&appendOpt.OutputEncoding, "output-encoding", "", "Encoding for output")
	cmdAppend.Flag.StringVar(&appendOpt.OutputEncoding, "oe", "", "Encoding for output")
	cmdAppend.Flag.StringVar(&appendOpt.Header, "header", "", "Appending header(s)")
	cmdAppend.Flag.StringVar(&appendOpt.Header, "h", "", "Appending header(s)")
	cmdAppend.Flag.IntVar(&appendOpt.Size, "size", 1, "Appending column size")
	cmdAppend.Flag.IntVar(&appendOpt.Size, "s", 1, "Appending column size")
}

// runAppend executes append command and return exit code.
func runAppend(args []string) int {
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

	opt := appendOpt.AppendOption
	opt.Headers = split(appendOpt.Header)
	err = csvutil.Append(r, w, opt)
	if err != nil {
		return handleError(err)
	}

	success = true
	return 0
}
