package main

import "github.com/pinzolo/csvutil"

var cmdTop = &Command{
	Run:       runTop,
	UsageLine: "top [OPTIONS...]",
	Short:     "先頭取得",
	Long: `DESCRIPTION
        先頭から指定した数の行だけを抽出したCSVを作成します。

OPTIONS
        -w, --overwrite
            指定されたCSVファイルを実行結果で上書きします。
            ファイルパスが渡されていない場合には無視されます。

        -H, --no-header
            このオプションを指定すると、ヘッダーの無いCSVを出力します。

        -b, --backup
            処理が成功した場合に、指定されたCSVファイルをバックアップします。
            --overwrite オプションと同時に使用されることを想定しているため、ファイルパスが渡されていない場合には無視されます。

        -oe, --output-encoding ENCODING
            出力するCSVの文字エンコーディングを指定します。
            このオプションが指定されていない場合、UTF-8とみなして処理を行います。
            対応している値:
                utf8bom : UTF-8として出力します（BOMは出力します）
                sjis    : Shift_JISとして出力します
                eucjp   : EUC_JPとして出力します

        -h, --header HEADER(S)
            新規に追加する列のヘッダーテキストを指定します。
            複数のヘッダーテキストを指定する場合には、foo:bar のようにコロン区切りにします。
            このオプションが指定されていない場合、もしくは指定したヘッダーが --size オプションの値に足りない場合には、
            column1,column2... のように連番のヘッダーが自動で付与されます。
            また --size オプションの値を超えた場合は、超えた分が無視されます。

        -c, --count NUMBER
            抽出する行の数を指定します。初期値は 1 です。
	`,
}

type cmdTopOption struct {
	csvutil.TopOption
	Overwrite bool
	Backup    bool
}

var topOpt = cmdTopOption{}

func init() {
	cmdTop.Flag.BoolVar(&topOpt.Overwrite, "overwrite", false, "Overwrite to source.")
	cmdTop.Flag.BoolVar(&topOpt.Overwrite, "w", false, "Overwrite to source.")
	cmdTop.Flag.BoolVar(&topOpt.NoHeader, "no-header", false, "Source file does not have header line")
	cmdTop.Flag.BoolVar(&topOpt.NoHeader, "H", false, "Source file does not have header line")
	cmdTop.Flag.BoolVar(&topOpt.Backup, "backup", false, "Backup source file.")
	cmdTop.Flag.BoolVar(&topOpt.Backup, "b", false, "Backup source file.")
	cmdTop.Flag.StringVar(&topOpt.Encoding, "encoding", "utf8", "Encoding of source file")
	cmdTop.Flag.StringVar(&topOpt.Encoding, "e", "utf8", "Encoding of source file")
	cmdTop.Flag.StringVar(&topOpt.OutputEncoding, "output-encoding", "utf8", "Encoding for output")
	cmdTop.Flag.StringVar(&topOpt.OutputEncoding, "oe", "utf8", "Encoding for output")
	cmdTop.Flag.IntVar(&topOpt.Count, "count", 1, "Toping line count")
	cmdTop.Flag.IntVar(&topOpt.Count, "c", 1, "Toping line count")
}

// runTop executes Top command and return exit code.
func runTop(args []string) int {
	success := false
	w, wf, r, rf, err := prepare(args, topOpt.Overwrite)
	if wf != nil {
		defer wf(&success, topOpt.Backup)
	}
	if rf != nil {
		defer rf()
	}
	if err != nil {
		return handleError(err)
	}

	err = csvutil.Top(r, w, topOpt.TopOption)
	if err != nil {
		return handleError(err)
	}

	success = true
	return 0
}
