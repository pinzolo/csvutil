package main

import "testing"

func Example_runConvert() {
	runConvert([]string{testFilePath("utf8.csv")})
	// Output: |  名前  | 個数 |
	// |--------|------|
	// | りんご |    1 |
	// | みかん |    2 |
}

func Example_runConvertWithTemplate() {
	convertOpt.Template = testFilePath("template.tmpl")
	runConvert([]string{testFilePath("utf8.csv")})
	convertOpt.Template = ""
	// Output: [
	//   { '名前' => 'りんご', '個数' => '1' }
	//   { '名前' => 'みかん', '個数' => '2' }
	// ]
}

func Test_runConvert(t *testing.T) {
	if c := runConvert([]string{testFilePath("utf8.csv")}); c != 0 {
		t.Errorf("Invalid success exit code: %d", c)
	}
}

func Test_runConvertOnNoFile(t *testing.T) {
	if c := runConvert([]string{testFilePath("no-file.csv")}); c == 0 {
		t.Errorf("Invalid failed exit code: %d", c)
	}
}

func Test_runConvertOnFail(t *testing.T) {
	if c := runConvert([]string{testFilePath("broken.csv")}); c == 0 {
		t.Errorf("Invalid failed exit code: %d", c)
	}
}

func Test_runConvertOnTemplateNotFound(t *testing.T) {
	convertOpt.Template = testFilePath("no-template.tmpl")
	if c := runConvert([]string{testFilePath("utf8.csv")}); c == 0 {
		t.Errorf("Invalid failed exit code: %d", c)
	}
	convertOpt.Template = ""
}
