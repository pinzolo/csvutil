package main

import (
	"github.com/pinzolo/csvutil"
)

var cmdBuilding = &Command{
	Run:       runBuilding,
	UsageLine: "building [OPTIONS...] [FILE]",
	Short:     "建物生成",
	Long: `DESCRIPTION
        指定した列にダミーの建物を出力します。
        自宅の場合は部屋番号、勤務先の場合はフロアを建物名似合わせて出力します。

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

        -nw, --number-width NUMBER
            このオプションに 1 を渡すと半角で、2 を渡すと全角で部屋番号・フロアを出力します。
            初期値は 1 です。

        -or, --office-rate PERCENTAGE
            ダミーの勤務先向け建物を設定する割合を指定します。0〜100までの整数を指定して下さい。（初期値: 0）
            このオプションはランダムなビル名とランダムなフロアを出力します。
	`,
}

type cmdBuildingOption struct {
	csvutil.BuildingOption
	Overwrite bool
	Backup    bool
}

var buildingOpt = cmdBuildingOption{}

func init() {
	cmdBuilding.Flag.BoolVar(&buildingOpt.Overwrite, "overwrite", false, "Overwrite to source.")
	cmdBuilding.Flag.BoolVar(&buildingOpt.Overwrite, "w", false, "Overwrite to source.")
	cmdBuilding.Flag.BoolVar(&buildingOpt.NoHeader, "no-header", false, "Source file does not have header line.")
	cmdBuilding.Flag.BoolVar(&buildingOpt.NoHeader, "H", false, "Source file does not have header line.")
	cmdBuilding.Flag.BoolVar(&buildingOpt.Backup, "backup", false, "Backup source file.")
	cmdBuilding.Flag.BoolVar(&buildingOpt.Backup, "b", false, "Backup source file.")
	cmdBuilding.Flag.StringVar(&buildingOpt.Encoding, "encoding", "utf8", "Encoding of source file")
	cmdBuilding.Flag.StringVar(&buildingOpt.Encoding, "e", "utf8", "Encoding of source file")
	cmdBuilding.Flag.StringVar(&buildingOpt.OutputEncoding, "output-encoding", "", "Encoding for output")
	cmdBuilding.Flag.StringVar(&buildingOpt.OutputEncoding, "oe", "", "Encoding for output")
	cmdBuilding.Flag.StringVar(&buildingOpt.Column, "column", "", "Target column symbol")
	cmdBuilding.Flag.StringVar(&buildingOpt.Column, "c", "", "Target column symbol")
	cmdBuilding.Flag.IntVar(&buildingOpt.OfficeRate, "office-rate", 0, "Office rate")
	cmdBuilding.Flag.IntVar(&buildingOpt.OfficeRate, "or", 0, "Office rate")
	cmdBuilding.Flag.IntVar(&buildingOpt.NumberWidth, "number-width", 1, "Number character width")
	cmdBuilding.Flag.IntVar(&buildingOpt.NumberWidth, "nw", 1, "Number character width")
	cmdBuilding.Flag.BoolVar(&buildingOpt.Append, "append", false, "Appen to source value")
	cmdBuilding.Flag.BoolVar(&buildingOpt.Append, "a", false, "Appen to source value")
}

// runBuilding executes building command and return exit code.
func runBuilding(args []string) int {
	success := false
	w, wf, r, rf, err := prepare(args, buildingOpt.Overwrite)
	if wf != nil {
		defer wf(&success, buildingOpt.Backup)
	}
	if rf != nil {
		defer rf()
	}
	if err != nil {
		return handleError(err)
	}

	err = csvutil.Building(r, w, buildingOpt.BuildingOption)
	if err != nil {
		return handleError(err)
	}

	success = true
	return 0
}
