package main

import (
	"github.com/pinzolo/csvutil"
)

var cmdInsert = &Command{
	Run:       runInsert,
	UsageLine: "insert [OPTIONS...] [FILE]",
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

type cmdInsertOption struct {
	csvutil.InsertOption
	// Overwrite to source. (default false)
	Overwrite bool
	// Backup source file. (default false)
	Backup bool
	// Header symbols.
	Header string
}

var insertOpt = cmdInsertOption{}

func init() {
	cmdInsert.Flag.BoolVar(&insertOpt.Overwrite, "overwrite", false, "Overwrite to source.")
	cmdInsert.Flag.BoolVar(&insertOpt.Overwrite, "w", false, "Overwrite to source.")
	cmdInsert.Flag.BoolVar(&insertOpt.NoHeader, "no-header", false, "Source file does not have header line.")
	cmdInsert.Flag.BoolVar(&insertOpt.NoHeader, "H", false, "Source file does not have header line.")
	cmdInsert.Flag.BoolVar(&insertOpt.Backup, "backup", false, "Backup source file.")
	cmdInsert.Flag.BoolVar(&insertOpt.Backup, "b", false, "Backup source file.")
	cmdInsert.Flag.StringVar(&insertOpt.Encoding, "encoding", "utf8", "Encoding of source file")
	cmdInsert.Flag.StringVar(&insertOpt.Encoding, "e", "utf8", "Encoding of source file")
	cmdInsert.Flag.StringVar(&insertOpt.OutputEncoding, "output-encoding", "", "Encoding for output")
	cmdInsert.Flag.StringVar(&insertOpt.OutputEncoding, "oe", "", "Encoding for output")
	cmdInsert.Flag.StringVar(&insertOpt.Header, "header", "", "Inserting header(s)")
	cmdInsert.Flag.StringVar(&insertOpt.Header, "h", "", "Inserting header(s)")
	cmdInsert.Flag.StringVar(&insertOpt.Before, "before", "", "Insert before this column")
	cmdInsert.Flag.StringVar(&insertOpt.Before, "bf", "", "Insert before this column")
	cmdInsert.Flag.IntVar(&insertOpt.Size, "size", 1, "Inserting column size")
	cmdInsert.Flag.IntVar(&insertOpt.Size, "s", 1, "Inserting column size")
}

// runInsert executes insert command and return exit code.
func runInsert(args []string) int {
	success := false
	path, err := path(args)
	if err != nil {
		return handleError(err)
	}

	w, wf, err := writer(path, insertOpt.Overwrite)
	if err != nil {
		return handleError(err)
	}
	if wf != nil {
		defer wf(&success, insertOpt.Backup)
	}

	r, rf, err := reader(path)
	if err != nil {
		return handleError(err)
	}
	if rf != nil {
		defer rf()
	}

	opt := insertOpt.InsertOption
	opt.Headers = split(insertOpt.Header)
	err = csvutil.Insert(r, w, opt)
	if err != nil {
		return handleError(err)
	}

	success = true
	return 0
}
