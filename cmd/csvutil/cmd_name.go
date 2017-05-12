package main

import (
	"github.com/pinzolo/csvutil"
)

var cmdName = &Command{
	Run:       runName,
	UsageLine: "name [OPTIONS...] [FILE]",
	Short:     "名前出力",
	Long: `DESCRIPTION
        指定した列にダミーの名前を出力します。

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

        -n, --name COLUMN_SYMBOL
            フルネーム（漢字）を出力する列のシンボルを指定します。

        -ln, --last-name COLUMN_SYMBOL
            姓（漢字）を出力する列のシンボルを指定します。

        -fn, --first-name COLUMN_SYMBOL
            名（漢字）を出力する列のシンボルを指定します。

        -k, --kana COLUMN_SYMBOL
            フルネーム（カナ）を出力する列のシンボルを指定します。

        -lk, --last-kana COLUMN_SYMBOL
            姓（カナ）を出力する列のシンボルを指定します。

        -fk, --first-kana COLUMN_SYMBOL
            名（カナ）を出力する列のシンボルを指定します。

        -h, --hiragana
            このオプションを指定すると仮名をひらがなで出力します。

        -g, --gender COLUMN_SYMBOL
            性別を出力する列のシンボルを指定します。

        -gf, --gender-format FORMAT
            性別を出力するフォーマットを指定します。初期値は jp_short です。
                 FORMAT   |  男  |   女
                --------- | ---- | ------
                 code     | 1    | 2
                 en_short | M    | F
                 en_long  | Male | Female
                 jp_short | 男   | 女
                 jp_long  | 男性 | 女性
                 symbol   | ♂   | ♀

        -mr, --male-rate PERCENTAGE
            男性を出力する割合を指定します。初期値は 50 です。

        -r, --reference COLUMN_SYMBOL
            保護者・家族として出力するために参照する列のシンボルを指定します。
            この列は漢字・ひらがな・カタカナで出力されている必要があります。
            また、姓と名の間は半角スペースから全角スペースで区切られている必要があります。
            スペースで区切られていない場合、列の値全体が姓として扱われます。

        -rr, --ristrict-reference
            --reference オプションで参照した値から名前を生成できなかった場合、エラーを起こします。
            指定しない場合、名前を生成できなかった場合には名前には空白が出力されます。

        -sw, --space-width NUMBER
            姓と名の区切り文字となる空白文字を数字で指定します。
                0: 空文字（初期値）
                1: 半角スペース [0x20]
                2: 全角スペース [0xE3 0x80 0x80]
	`,
}

type cmdNameOption struct {
	csvutil.NameOption
	Overwrite bool
	Backup    bool
}

var nameOpt = cmdNameOption{}

func init() {
	cmdName.Flag.BoolVar(&nameOpt.Overwrite, "overwrite", false, "Overwrite to source.")
	cmdName.Flag.BoolVar(&nameOpt.Overwrite, "w", false, "Overwrite to source.")
	cmdName.Flag.BoolVar(&nameOpt.NoHeader, "no-header", false, "Source file does not have header line.")
	cmdName.Flag.BoolVar(&nameOpt.NoHeader, "H", false, "Source file does not have header line.")
	cmdName.Flag.BoolVar(&nameOpt.Backup, "backup", false, "Backup source file.")
	cmdName.Flag.BoolVar(&nameOpt.Backup, "b", false, "Backup source file.")
	cmdName.Flag.StringVar(&nameOpt.Encoding, "encoding", "utf8", "Encoding of source file")
	cmdName.Flag.StringVar(&nameOpt.Encoding, "e", "utf8", "Encoding of source file")
	cmdName.Flag.StringVar(&nameOpt.OutputEncoding, "output-encoding", "", "Encoding for output")
	cmdName.Flag.StringVar(&nameOpt.OutputEncoding, "oe", "", "Encoding for output")
	cmdName.Flag.StringVar(&nameOpt.Name, "name", "", "Name column symbol")
	cmdName.Flag.StringVar(&nameOpt.Name, "n", "", "Name column symbol")
	cmdName.Flag.StringVar(&nameOpt.FirstName, "first-name", "", "First name column symbol")
	cmdName.Flag.StringVar(&nameOpt.FirstName, "fn", "", "First name column symbol")
	cmdName.Flag.StringVar(&nameOpt.LastName, "last-name", "", "Last name column symbol")
	cmdName.Flag.StringVar(&nameOpt.LastName, "ln", "", "Last name column symbol")
	cmdName.Flag.StringVar(&nameOpt.Kana, "kana", "", "Kana column symbol")
	cmdName.Flag.StringVar(&nameOpt.Kana, "k", "", "Kana column symbol")
	cmdName.Flag.StringVar(&nameOpt.FirstKana, "first-kana", "", "First kana column symbol")
	cmdName.Flag.StringVar(&nameOpt.FirstKana, "fk", "", "First kana column symbol")
	cmdName.Flag.StringVar(&nameOpt.LastKana, "last-kana", "", "Last kana column symbol")
	cmdName.Flag.StringVar(&nameOpt.LastKana, "lk", "", "Last kana column symbol")
	cmdName.Flag.BoolVar(&nameOpt.Hiragana, "hiragana", false, "Output hiragana as kana")
	cmdName.Flag.BoolVar(&nameOpt.Hiragana, "h", false, "Output hiragana as kana")
	cmdName.Flag.StringVar(&nameOpt.Gender, "gender", "", "Gender column symbol")
	cmdName.Flag.StringVar(&nameOpt.Gender, "g", "", "Gender column symbol")
	cmdName.Flag.IntVar(&nameOpt.MaleRate, "male-rate", 50, "Male rate")
	cmdName.Flag.IntVar(&nameOpt.MaleRate, "mr", 50, "Male rate")
	cmdName.Flag.StringVar(&nameOpt.GenderFormat, "gender-format", "jp_short", "Gender format")
	cmdName.Flag.StringVar(&nameOpt.GenderFormat, "gf", "jp_short", "Gender format")
	cmdName.Flag.StringVar(&nameOpt.Reference, "reference", "", "Reference column symbol")
	cmdName.Flag.StringVar(&nameOpt.Reference, "r", "", "Reference column symbol")
	cmdName.Flag.BoolVar(&nameOpt.RistrictReference, "ristrict-reference", false, "Raise error reference not found")
	cmdName.Flag.BoolVar(&nameOpt.RistrictReference, "rr", false, "Raise error reference not found")
	cmdName.Flag.IntVar(&nameOpt.SpaceWidth, "space-width", 1, "Delimiter space width")
	cmdName.Flag.IntVar(&nameOpt.SpaceWidth, "sw", 1, "Delimiter space width")
}

// runName executes name command and return exit code.
func runName(args []string) int {
	success := false
	w, wf, r, rf, err := prepare(args, nameOpt.Overwrite)
	if wf != nil {
		defer wf(&success, nameOpt.Backup)
	}
	if rf != nil {
		defer rf()
	}
	if err != nil {
		return handleError(err)
	}

	err = csvutil.Name(r, w, nameOpt.NameOption)
	if err != nil {
		return handleError(err)
	}

	success = true
	return 0
}
