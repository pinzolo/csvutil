package main

import (
	"github.com/pinzolo/csvutil"
)

var cmdEmail = &Command{
	Run:       runEmail,
	UsageLine: "email [OPTIONS...] [FILE]",
	Short:     "メールアドレス出力",
	Long: `DESCRIPTION
        指定した列にダミーのメールアドレスを出力します。

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

        -mr, --mobile-rate PERCENTAGE
            ダミーの携帯電話メールアドレスを設定する割合を指定します。0〜100までの整数を指定して下さい。（初期値: 0）
            csvutilが埋め込む携帯電話メールアドレスは下記のドメインを持つランダムなメールアドレスです。
                docomo.ne.jp
                ezweb.ne.jp
                softbank.ne.jp
                i.softbank.ne.jp
                ymobile.ne.jp
                emobile.ne.jp
	`,
}

type cmdEmailOption struct {
	csvutil.EmailOption
	Overwrite bool
	Backup    bool
}

var emailOpt = cmdEmailOption{}

func init() {
	cmdEmail.Flag.BoolVar(&emailOpt.Overwrite, "overwrite", false, "Overwrite to source.")
	cmdEmail.Flag.BoolVar(&emailOpt.Overwrite, "w", false, "Overwrite to source.")
	cmdEmail.Flag.BoolVar(&emailOpt.NoHeader, "no-header", false, "Source file does not have header line.")
	cmdEmail.Flag.BoolVar(&emailOpt.NoHeader, "H", false, "Source file does not have header line.")
	cmdEmail.Flag.BoolVar(&emailOpt.Backup, "backup", false, "Backup source file.")
	cmdEmail.Flag.BoolVar(&emailOpt.Backup, "b", false, "Backup source file.")
	cmdEmail.Flag.StringVar(&emailOpt.Encoding, "encoding", "utf8", "Encoding of source file")
	cmdEmail.Flag.StringVar(&emailOpt.Encoding, "e", "utf8", "Encoding of source file")
	cmdEmail.Flag.StringVar(&emailOpt.OutputEncoding, "output-encoding", "", "Encoding for output")
	cmdEmail.Flag.StringVar(&emailOpt.OutputEncoding, "oe", "", "Encoding for output")
	cmdEmail.Flag.StringVar(&emailOpt.Column, "column", "", "Target column symbol")
	cmdEmail.Flag.StringVar(&emailOpt.Column, "c", "", "Target column symbol")
	cmdEmail.Flag.IntVar(&emailOpt.MobileRate, "mobile-rate", 0, "Mobile email address rate")
	cmdEmail.Flag.IntVar(&emailOpt.MobileRate, "mr", 0, "Mobile email address rate")
}

// runEmail executes email command and return exit code.
func runEmail(args []string) int {
	success := false
	w, wf, r, rf, err := prepare(args, emailOpt.Overwrite)
	if wf != nil {
		defer wf(&success, emailOpt.Backup)
	}
	if rf != nil {
		defer rf()
	}
	if err != nil {
		return handleError(err)
	}

	err = csvutil.Email(r, w, emailOpt.EmailOption)
	if err != nil {
		return handleError(err)
	}

	success = true
	return 0
}
