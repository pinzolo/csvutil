package main

import (
	"io/ioutil"
	"os"

	"github.com/pinzolo/csvutil"
)

var cmdConvert = &Command{
	Run:       runConvert,
	UsageLine: "convert [OPTIONS...] [FILE]",
	Short:     "形式変換",
	Long: `DESCRIPTION
        CSVのデータを特定の形式に変換して出力します。出力エンコードは UTF-8 固定です。

ARGUMENTS
        FILE
            ソースとなる CSV ファイルのパスを指定します。
            パスが指定されていない場合、標準入力が対象となりパイプでの使用ができます。

OPTIONS
        -H, --no-header
            ソースとなるCSVの1行目をヘッダー列として扱いません。

        -e, --encoding ENCODING
            ソースとなるCSVの文字エンコーディングを指定します。
            このオプションが指定されていない場合、csvutil はUTF-8とみなして処理を行います。
            UTF-8であった場合、BOMのあるなしは自動的に判別されます。
            対応している値:
                sjis : Shift_JISとして扱います
                eucjp: EUC_JPとして扱います

        -f, --format FORMAT
            変換先の組み込み済みフォーマットを指定します。初期値は markdown です。
            対応している値:
                markdown: Markdownのテーブル書式に変換します
                json    : JSONに変換します
                yaml    : YAMLに変換します
                html    : HTMLのテーブル書式に変換します
                xml     : XMLに変換します

        -t, --template FILE
            テンプレートが記述されたファイルを指定します。
            このオプションは --format オプションより優先されます。
            テンプレートの書式は Go 言語に準拠します。
            https://golang.org/pkg/text/template/
            https://golang.org/pkg/html/template/

            テンプレートにコンテキストとして渡されるオブジェクトには、Headers プロパティと Data プロパティがあります。
            Headers は []string, Data は [][]string で定義されており、それぞれヘッダーと全データが格納されています。

        -h, --html
            このオプションを指定すると --template オプションで指定されたテンプレートを適用する際に、
            text/template パッケージではなく html/template パッケージを使用して変換します。
	`,
}

type cmdConvertOption struct {
	csvutil.ConvertOption
}

var convertOpt = cmdConvertOption{}

func init() {
	cmdConvert.Flag.BoolVar(&convertOpt.NoHeader, "no-header", false, "Source file does not have header line.")
	cmdConvert.Flag.BoolVar(&convertOpt.NoHeader, "H", false, "Source file does not have header line.")
	cmdConvert.Flag.StringVar(&convertOpt.Encoding, "encoding", "utf8", "Encoding of source file")
	cmdConvert.Flag.StringVar(&convertOpt.Encoding, "e", "utf8", "Encoding of source file")
	cmdConvert.Flag.StringVar(&convertOpt.Format, "format", "markdown", "Format of converter")
	cmdConvert.Flag.StringVar(&convertOpt.Format, "f", "markdown", "Format of converter")
	cmdConvert.Flag.StringVar(&convertOpt.Template, "template", "", "Template file path")
	cmdConvert.Flag.StringVar(&convertOpt.Template, "t", "", "Template file path")
	cmdConvert.Flag.BoolVar(&convertOpt.HTML, "html", false, "Template is HTML")
	cmdConvert.Flag.BoolVar(&convertOpt.HTML, "h", false, "Template is HTML")
}

// runConvert executes convert command and return exit code.
func runConvert(args []string) int {
	r, rf, err := prepareReader(args)
	if rf != nil {
		defer rf()
	}
	if err != nil {
		return handleError(err)
	}

	if convertOpt.Template != "" {
		var f *os.File
		f, err = os.Open(convertOpt.Template)
		if err != nil {
			return handleError(err)
		}
		var p []byte
		p, err = ioutil.ReadAll(f)
		if err != nil {
			return handleError(err)
		}
		convertOpt.Template = string(p)
	}
	err = csvutil.Convert(r, os.Stdout, convertOpt.ConvertOption)
	if err != nil {
		return handleError(err)
	}

	return 0
}
