package main

import (
	"github.com/pinzolo/csvutil"
)

var cmdPassword = &Command{
	Run:       runPassword,
	UsageLine: "password [OPTIONS...] [FILE]",
	Short:     "パスワード出力",
	Long: `DESCRIPTION
        指定した列にダミーのパスワードを出力します。

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

        -mn, --min-length LENGTH
            生成するパスワードの最少長を指定します。（初期値：8）
            この値は正の整数でなければなりません。

        -mx, --max-length LENGTH
            生成するパスワードの最大長を指定します。（初期値：16）
            この値は --min-length 以上の正の整数でなければなりません。

        -N, --no-numeric
            生成するパスワードに数字を含めません。

        -U, --no-upper
            生成するパスワードに大文字アルファベットを含めません。

        -S, --no-special
            生成するパスワードに特殊文字（記号）を含めません。
	`,
}

type cmdPasswordOption struct {
	csvutil.PasswordOption
	Overwrite bool
	Backup    bool
}

var passwordOpt = cmdPasswordOption{}

func init() {
	cmdPassword.Flag.BoolVar(&passwordOpt.Overwrite, "overwrite", false, "Overwrite to source.")
	cmdPassword.Flag.BoolVar(&passwordOpt.Overwrite, "w", false, "Overwrite to source.")
	cmdPassword.Flag.BoolVar(&passwordOpt.NoHeader, "no-header", false, "Source file does not have header line.")
	cmdPassword.Flag.BoolVar(&passwordOpt.NoHeader, "H", false, "Source file does not have header line.")
	cmdPassword.Flag.BoolVar(&passwordOpt.Backup, "backup", false, "Backup source file.")
	cmdPassword.Flag.BoolVar(&passwordOpt.Backup, "b", false, "Backup source file.")
	cmdPassword.Flag.StringVar(&passwordOpt.Encoding, "encoding", "utf8", "Encoding of source file")
	cmdPassword.Flag.StringVar(&passwordOpt.Encoding, "e", "utf8", "Encoding of source file")
	cmdPassword.Flag.StringVar(&passwordOpt.OutputEncoding, "output-encoding", "", "Encoding for output")
	cmdPassword.Flag.StringVar(&passwordOpt.OutputEncoding, "oe", "", "Encoding for output")
	cmdPassword.Flag.StringVar(&passwordOpt.Column, "column", "", "Target column symbol")
	cmdPassword.Flag.StringVar(&passwordOpt.Column, "c", "", "Target column symbol")
	cmdPassword.Flag.IntVar(&passwordOpt.MinLength, "min-length", 8, "Min length of password")
	cmdPassword.Flag.IntVar(&passwordOpt.MinLength, "mn", 8, "Min length of password")
	cmdPassword.Flag.IntVar(&passwordOpt.MaxLength, "max-length", 16, "Max length of password")
	cmdPassword.Flag.IntVar(&passwordOpt.MaxLength, "mx", 16, "Max length of password")
	cmdPassword.Flag.BoolVar(&passwordOpt.NoNumeric, "no-numeric", false, "Not use numeric in password")
	cmdPassword.Flag.BoolVar(&passwordOpt.NoNumeric, "N", false, "Not use numeric in password")
	cmdPassword.Flag.BoolVar(&passwordOpt.NoUpper, "no-upper", false, "Not use upper alphabet in password")
	cmdPassword.Flag.BoolVar(&passwordOpt.NoUpper, "U", false, "Not use upper alphabet in password")
	cmdPassword.Flag.BoolVar(&passwordOpt.NoSpecial, "no-special", false, "Not use special char in password")
	cmdPassword.Flag.BoolVar(&passwordOpt.NoSpecial, "S", false, "Not use special char in password")
}

// runPassword executes password command and return exit code.
func runPassword(args []string) int {
	success := false
	w, wf, r, rf, err := prepare(args, passwordOpt.Overwrite)
	if wf != nil {
		defer wf(&success, passwordOpt.Backup)
	}
	if rf != nil {
		defer rf()
	}
	if err != nil {
		return handleError(err)
	}

	err = csvutil.Password(r, w, passwordOpt.PasswordOption)
	if err != nil {
		return handleError(err)
	}

	success = true
	return 0
}
