package main

import (
	"github.com/pinzolo/csvutil"
)

var cmdFilter = &Command{
	Run:       runFilter,
	UsageLine: "filter [OPTIONS...] [FILE]",
	Short:     "行抽出",
	Long: `DESCRIPTION
        指定したパターンに合致する値を持つ行のみのCSVを出力します。

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
            複数列を対象としたい場合は、foo:bar や 1:2のようにコロン区切りで指定して下さい。
            このオプションを指定しない場合、すべての列が対象になります。

        -p, --pattern PATTERN
            置換対象のパターンです。ここに指定したパターンにマッチする文字列を持つ行を抽出します。

        -re, --regex, --regexp
            このオプションが指定されると --pattern に指定された値は正規表現と見なされます。
            初期値は false で、単純な曖昧検索を行います。
	`,
}

type cmdFilterOption struct {
	csvutil.FilterOption
	Overwrite bool
	Backup    bool
	Column    string
}

var filterOpt = cmdFilterOption{}

func init() {
	cmdFilter.Flag.BoolVar(&filterOpt.Overwrite, "overwrite", false, "Overwrite to source.")
	cmdFilter.Flag.BoolVar(&filterOpt.Overwrite, "w", false, "Overwrite to source.")
	cmdFilter.Flag.BoolVar(&filterOpt.NoHeader, "no-header", false, "Source file does not have header line.")
	cmdFilter.Flag.BoolVar(&filterOpt.NoHeader, "H", false, "Source file does not have header line.")
	cmdFilter.Flag.BoolVar(&filterOpt.Backup, "backup", false, "Backup source file.")
	cmdFilter.Flag.BoolVar(&filterOpt.Backup, "b", false, "Backup source file.")
	cmdFilter.Flag.StringVar(&filterOpt.Encoding, "encoding", "utf8", "Encoding of source file")
	cmdFilter.Flag.StringVar(&filterOpt.Encoding, "e", "utf8", "Encoding of source file")
	cmdFilter.Flag.StringVar(&filterOpt.OutputEncoding, "output-encoding", "", "Encoding for output")
	cmdFilter.Flag.StringVar(&filterOpt.OutputEncoding, "oe", "", "Encoding for output")
	cmdFilter.Flag.StringVar(&filterOpt.Column, "column", "", "Target column symbol")
	cmdFilter.Flag.StringVar(&filterOpt.Column, "c", "", "Home column symbol")
	cmdFilter.Flag.StringVar(&filterOpt.Pattern, "pattern", "", "Pattern")
	cmdFilter.Flag.StringVar(&filterOpt.Pattern, "p", "", "Pattern")
	cmdFilter.Flag.BoolVar(&filterOpt.Regexp, "regexp", false, "Pattern is regex")
	cmdFilter.Flag.BoolVar(&filterOpt.Regexp, "regex", false, "Pattern is regex")
	cmdFilter.Flag.BoolVar(&filterOpt.Regexp, "re", false, "Pattern is regex")
}

// runFilter executes filter command and return exit code.
func runFilter(args []string) int {
	success := false
	w, wf, r, rf, err := prepare(args, filterOpt.Overwrite)
	if wf != nil {
		defer wf(&success, filterOpt.Backup)
	}
	if rf != nil {
		defer rf()
	}
	if err != nil {
		return handleError(err)
	}

	opt := filterOpt.FilterOption
	opt.ColumnSyms = split(filterOpt.Column)
	err = csvutil.Filter(r, w, opt)
	if err != nil {
		return handleError(err)
	}

	success = true
	return 0
}
