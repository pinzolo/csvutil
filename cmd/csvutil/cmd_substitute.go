package main

import (
	"github.com/pinzolo/csvutil"
)

var cmdSubstitute = &Command{
	Run:       runSubstitute,
	UsageLine: "substitute [OPTIONS...] [FILE]",
	Short:     "文字列置換",
	Long: `DESCRIPTION
        指定した列の値を置換したCSVを出力します。

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

        -p, --pattern PATTERN
            置換対象のパターンです。ここに指定したパターンにマッチする文字列が置き換えられます。

        -r, --replacement TEXT
            実行すると --pattern に指定された文字列がこの値に置き換えられます。
            初期値は空文字なので、このオプションを指定せずに実行するとマッチする文字が削除されます。

        -re, --regex, --regexp
            このオプションが指定されると --pattern に指定された値は正規表現と見なされます。
            初期値は false で、単純な文字列置換を行います。
	`,
}

type cmdSubstituteOption struct {
	csvutil.SubstituteOption
	Overwrite bool
	Backup    bool
}

var substituteOpt = cmdSubstituteOption{}

func init() {
	cmdSubstitute.Flag.BoolVar(&substituteOpt.Overwrite, "overwrite", false, "Overwrite to source.")
	cmdSubstitute.Flag.BoolVar(&substituteOpt.Overwrite, "w", false, "Overwrite to source.")
	cmdSubstitute.Flag.BoolVar(&substituteOpt.NoHeader, "no-header", false, "Source file does not have header line.")
	cmdSubstitute.Flag.BoolVar(&substituteOpt.NoHeader, "H", false, "Source file does not have header line.")
	cmdSubstitute.Flag.BoolVar(&substituteOpt.Backup, "backup", false, "Backup source file.")
	cmdSubstitute.Flag.BoolVar(&substituteOpt.Backup, "b", false, "Backup source file.")
	cmdSubstitute.Flag.StringVar(&substituteOpt.Encoding, "encoding", "utf8", "Encoding of source file")
	cmdSubstitute.Flag.StringVar(&substituteOpt.Encoding, "e", "utf8", "Encoding of source file")
	cmdSubstitute.Flag.StringVar(&substituteOpt.OutputEncoding, "output-encoding", "", "Encoding for output")
	cmdSubstitute.Flag.StringVar(&substituteOpt.OutputEncoding, "oe", "", "Encoding for output")
	cmdSubstitute.Flag.StringVar(&substituteOpt.Column, "column", "", "Home column symbol")
	cmdSubstitute.Flag.StringVar(&substituteOpt.Column, "c", "", "Home column symbol")
	cmdSubstitute.Flag.StringVar(&substituteOpt.Pattern, "pattern", "", "Pattern")
	cmdSubstitute.Flag.StringVar(&substituteOpt.Pattern, "p", "", "Pattern")
	cmdSubstitute.Flag.StringVar(&substituteOpt.Replacement, "replacement", "", "Replacement")
	cmdSubstitute.Flag.StringVar(&substituteOpt.Replacement, "r", "", "Replacement")
	cmdSubstitute.Flag.BoolVar(&substituteOpt.Regexp, "regexp", false, "Pattern is regex")
	cmdSubstitute.Flag.BoolVar(&substituteOpt.Regexp, "regex", false, "Pattern is regex")
	cmdSubstitute.Flag.BoolVar(&substituteOpt.Regexp, "re", false, "Pattern is regex")
}

// runSubstitute executes substitute command and return exit code.
func runSubstitute(args []string) int {
	success := false
	w, wf, r, rf, err := prepare(args, substituteOpt.Overwrite)
	if wf != nil {
		defer wf(&success, substituteOpt.Backup)
	}
	if rf != nil {
		defer rf()
	}
	if err != nil {
		return handleError(err)
	}

	err = csvutil.Substitute(r, w, substituteOpt.SubstituteOption)
	if err != nil {
		return handleError(err)
	}

	success = true
	return 0
}
