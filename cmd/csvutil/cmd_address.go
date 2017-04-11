package main

import (
	"github.com/pinzolo/csvutil"
)

var cmdAddress = &Command{
	Run:       runAddress,
	UsageLine: "address [OPTIONS...] [FILE]",
	Short:     "住所出力",
	Long: `DESCRIPTION
        指定した列にダミーの住所を出力します。
        郵便番号、都道府県、都市、町は同じ列を指定すれば追記されます。
        都市と町の区分は厳格なものではなく、都市→町→番地の順に記載されるもの程度の扱いです。

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

        -z, --zip-code COLUMN_SYMBOL
            郵便番号を出力する列のシンボルを指定します。

        -p, --prefecture COLUMN_SYMBOL
            都道府県を出力する列のシンボルを指定します。

        -pc, --prefecture-code
            このオプションを指定すると都道府県を都道府県コードとして出力します。
            都道府県コードは 01 の用にゼロ埋めされています。

        -c, --city COLUMN_SYMBOL
            都市を出力する列のシンボルを指定します。

        -t, --town COLUMN_SYMBOL
            町を出力する列のシンボルを指定します。

        -bn, --block-number
            このオプションを指定すると --town オプションに対して、ダミーの番地を出力します。

        -nw, --number-width NUMBER
            このオプションに 1 を渡すと半角で、2 を渡すと全角で番地を出力します。
            初期値は 1 です。
	`,
}

type cmdAddressOption struct {
	csvutil.AddressOption
	Overwrite bool
	Backup    bool
}

var addressOpt = cmdAddressOption{}

func init() {
	cmdAddress.Flag.BoolVar(&addressOpt.Overwrite, "overwrite", false, "Overwrite to source.")
	cmdAddress.Flag.BoolVar(&addressOpt.Overwrite, "w", false, "Overwrite to source.")
	cmdAddress.Flag.BoolVar(&addressOpt.NoHeader, "no-header", false, "Source file does not have header line.")
	cmdAddress.Flag.BoolVar(&addressOpt.NoHeader, "H", false, "Source file does not have header line.")
	cmdAddress.Flag.BoolVar(&addressOpt.Backup, "backup", false, "Backup source file.")
	cmdAddress.Flag.BoolVar(&addressOpt.Backup, "b", false, "Backup source file.")
	cmdAddress.Flag.StringVar(&addressOpt.Encoding, "encoding", "utf8", "Encoding of source file")
	cmdAddress.Flag.StringVar(&addressOpt.Encoding, "e", "utf8", "Encoding of source file")
	cmdAddress.Flag.StringVar(&addressOpt.OutputEncoding, "output-encoding", "", "Encoding for output")
	cmdAddress.Flag.StringVar(&addressOpt.OutputEncoding, "oe", "", "Encoding for output")
	cmdAddress.Flag.StringVar(&addressOpt.ZipCode, "zip-code", "", "Zip code column symbol")
	cmdAddress.Flag.StringVar(&addressOpt.ZipCode, "z", "", "Zip code column symbol")
	cmdAddress.Flag.StringVar(&addressOpt.Prefecture, "prefecture", "", "Prefecture column symbol")
	cmdAddress.Flag.StringVar(&addressOpt.Prefecture, "p", "", "Prefecture column symbol")
	cmdAddress.Flag.BoolVar(&addressOpt.PrefectureCode, "prefecture-code", false, "Output prefecture code instead of prefecture name")
	cmdAddress.Flag.BoolVar(&addressOpt.PrefectureCode, "pc", false, "Output prefecture code instead of prefecture name")
	cmdAddress.Flag.StringVar(&addressOpt.City, "city", "", "City column symbol")
	cmdAddress.Flag.StringVar(&addressOpt.City, "c", "", "City column symbol")
	cmdAddress.Flag.StringVar(&addressOpt.Town, "town", "", "Town column symbol")
	cmdAddress.Flag.StringVar(&addressOpt.Town, "t", "", "Town column symbol")
	cmdAddress.Flag.BoolVar(&addressOpt.BlockNumber, "block-number", false, "Output block number after town")
	cmdAddress.Flag.BoolVar(&addressOpt.BlockNumber, "bn", false, "Output block number after town")
	cmdAddress.Flag.IntVar(&addressOpt.NumberWidth, "number-width", 1, "Block number character width")
	cmdAddress.Flag.IntVar(&addressOpt.NumberWidth, "nw", 1, "Block number character width")
}

// runAddress executes address command and return exit code.
func runAddress(args []string) int {
	success := false
	w, wf, r, rf, err := prepare(args, addressOpt.Overwrite)
	if wf != nil {
		defer wf(&success, addressOpt.Backup)
	}
	if rf != nil {
		defer rf()
	}
	if err != nil {
		return handleError(err)
	}

	err = csvutil.Address(r, w, addressOpt.AddressOption)
	if err != nil {
		return handleError(err)
	}

	success = true
	return 0
}
