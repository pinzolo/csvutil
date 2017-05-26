package csvutil

import (
	"bytes"
	"testing"
)

func TestConvertWithoutFormatAndTemplate(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := ConvertOption{}

	if err := Convert(r, w, o); err == nil {
		t.Error("Convert without format and template should raise error")
	}
}

func TestConvertWithUnsupportedFormat(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := ConvertOption{
		Format: "unknown",
	}

	if err := Convert(r, w, o); err == nil {
		t.Error("Convert with unsupported format should raise error")
	}
}

func TestConvertWithBrokenCSV(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := ConvertOption{
		Format: "markdown",
	}

	if err := Convert(r, w, o); err == nil {
		t.Error("Convert with broken CSV should raise error")
	}
}

func TestConvertWithBrokenTemplate(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := ConvertOption{
		Template: `{{ange $i, $v := .Headers}}{{if ne $i 0}}	{{end}}{{$v}}{{end}}
{{range .Data -}}
	{{range $i, $v := .}}{{if ne $i 0}}	{{end}}{{$v}}{{end}}
{{end}}`,
	}

	if err := Convert(r, w, o); err == nil {
		t.Error("Convert with broken template should raise error")
	}
}

func TestConvertWithEmptyCSV(t *testing.T) {
	s := ""
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := ConvertOption{
		Format: "markdown",
	}

	if err := Convert(r, w, o); err != nil {
		t.Error(err)
	}

	if actual := w.String(); actual != "" {
		t.Errorf("Expectd: empty, but got %s", actual)
	}
}

func TestConvertToMarkdown(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := ConvertOption{
		Format: "markdown",
	}

	if err := Convert(r, w, o); err != nil {
		t.Error(err)
	}

	expected := `| aaa | bbb | ccc |
|-----|-----|-----|
|   1 |   2 |   3 |
|   4 |   5 |   6 |
|   7 |   8 |   9 |
`
	if actual := w.String(); actual != expected {
		t.Errorf("Expectd: %s, but got %s", expected, actual)
	}
}

func TestConvertToMarkdownWithNoHeader(t *testing.T) {
	s := `1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := ConvertOption{
		Format:   "markdown",
		NoHeader: true,
	}

	if err := Convert(r, w, o); err != nil {
		t.Error(err)
	}

	expected := `| column1 | column2 | column3 |
|---------|---------|---------|
|       1 |       2 |       3 |
|       4 |       5 |       6 |
|       7 |       8 |       9 |
`
	if actual := w.String(); actual != expected {
		t.Errorf("Expectd: %s, but got %s", expected, actual)
	}
}

func TestConvertToJSON(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := ConvertOption{
		Format: "json",
	}

	if err := Convert(r, w, o); err != nil {
		t.Error(err)
	}

	expected := `[
	{
		"aaa": "1",
		"bbb": "2",
		"ccc": "3"
	},
	{
		"aaa": "4",
		"bbb": "5",
		"ccc": "6"
	},
	{
		"aaa": "7",
		"bbb": "8",
		"ccc": "9"
	}
]`
	if actual := w.String(); actual != expected {
		t.Errorf("Expectd: %s, but got %s", expected, actual)
	}
}

func TestConvertToJSONWithNoHeader(t *testing.T) {
	s := `1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := ConvertOption{
		Format:   "json",
		NoHeader: true,
	}

	if err := Convert(r, w, o); err != nil {
		t.Error(err)
	}

	expected := `[
	{
		"column1": "1",
		"column2": "2",
		"column3": "3"
	},
	{
		"column1": "4",
		"column2": "5",
		"column3": "6"
	},
	{
		"column1": "7",
		"column2": "8",
		"column3": "9"
	}
]`
	if actual := w.String(); actual != expected {
		t.Errorf("Expectd: %s, but got %s", expected, actual)
	}
}

func TestConvertToYAML(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := ConvertOption{
		Format: "yaml",
	}

	if err := Convert(r, w, o); err != nil {
		t.Error(err)
	}

	expected := `- aaa: "1"
  bbb: "2"
  ccc: "3"
- aaa: "4"
  bbb: "5"
  ccc: "6"
- aaa: "7"
  bbb: "8"
  ccc: "9"
`
	if actual := w.String(); actual != expected {
		t.Errorf("Expectd: %s, but got %s", expected, actual)
	}
}

func TestConvertToYAMLWithNoHeader(t *testing.T) {
	s := `1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := ConvertOption{
		Format:   "yaml",
		NoHeader: true,
	}

	if err := Convert(r, w, o); err != nil {
		t.Error(err)
	}

	expected := `- column1: "1"
  column2: "2"
  column3: "3"
- column1: "4"
  column2: "5"
  column3: "6"
- column1: "7"
  column2: "8"
  column3: "9"
`
	if actual := w.String(); actual != expected {
		t.Errorf("Expectd: %s, but got %s", expected, actual)
	}
}

func TestConvertToHTML(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := ConvertOption{
		Format: "html",
	}

	if err := Convert(r, w, o); err != nil {
		t.Error(err)
	}

	expected := `<table>
	<thead>
		<tr>
			<th>aaa</th>
			<th>bbb</th>
			<th>ccc</th>
		</tr>
	</thead>
	<tbody>
		<tr>
			<td>1</td>
			<td>2</td>
			<td>3</td>
		</tr>
		<tr>
			<td>4</td>
			<td>5</td>
			<td>6</td>
		</tr>
		<tr>
			<td>7</td>
			<td>8</td>
			<td>9</td>
		</tr>
	</tbody>
</table>
`
	if actual := w.String(); actual != expected {
		t.Errorf("Expectd: %s, but got %s", expected, actual)
	}
}

func TestConvertToHTMLWithNoHeader(t *testing.T) {
	s := `1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := ConvertOption{
		Format:   "html",
		NoHeader: true,
	}

	if err := Convert(r, w, o); err != nil {
		t.Error(err)
	}

	expected := `<table>
	<thead>
		<tr>
			<th>column1</th>
			<th>column2</th>
			<th>column3</th>
		</tr>
	</thead>
	<tbody>
		<tr>
			<td>1</td>
			<td>2</td>
			<td>3</td>
		</tr>
		<tr>
			<td>4</td>
			<td>5</td>
			<td>6</td>
		</tr>
		<tr>
			<td>7</td>
			<td>8</td>
			<td>9</td>
		</tr>
	</tbody>
</table>
`
	if actual := w.String(); actual != expected {
		t.Errorf("Expectd: %s, but got %s", expected, actual)
	}
}

func TestConvertToXML(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := ConvertOption{
		Format: "xml",
	}

	if err := Convert(r, w, o); err != nil {
		t.Error(err)
	}

	expected := `<csv>
	<data>
		<aaa>1</aaa>
		<bbb>2</bbb>
		<ccc>3</ccc>
	</data>
	<data>
		<aaa>4</aaa>
		<bbb>5</bbb>
		<ccc>6</ccc>
	</data>
	<data>
		<aaa>7</aaa>
		<bbb>8</bbb>
		<ccc>9</ccc>
	</data>
</csv>`
	if actual := w.String(); actual != expected {
		t.Errorf("Expectd: %s, but got %s", expected, actual)
	}
}

func TestConvertToXMLWithNoHeader(t *testing.T) {
	s := `1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := ConvertOption{
		Format:   "xml",
		NoHeader: true,
	}

	if err := Convert(r, w, o); err != nil {
		t.Error(err)
	}

	expected := `<csv>
	<data>
		<column1>1</column1>
		<column2>2</column2>
		<column3>3</column3>
	</data>
	<data>
		<column1>4</column1>
		<column2>5</column2>
		<column3>6</column3>
	</data>
	<data>
		<column1>7</column1>
		<column2>8</column2>
		<column3>9</column3>
	</data>
</csv>`
	if actual := w.String(); actual != expected {
		t.Errorf("Expectd: %s, but got %s", expected, actual)
	}
}

func TestConvertWithCustomTemplate(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,<5>,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := ConvertOption{
		Template: `{{range $i, $v := .Headers}}{{if ne $i 0}}	{{end}}{{$v}}{{end}}
{{range .Data -}}
	{{range $i, $v := .}}{{if ne $i 0}}	{{end}}{{$v}}{{end}}
{{end}}`,
	}

	if err := Convert(r, w, o); err != nil {
		t.Error(err)
	}

	expected := `aaa	bbb	ccc
1	2	3
4	<5>	6
7	8	9
`
	if actual := w.String(); actual != expected {
		t.Errorf("Expectd: %s, but got %s", expected, actual)
	}
}

func TestConvertWithCustomTemplateAndNoHeader(t *testing.T) {
	s := `1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := ConvertOption{
		NoHeader: true,
		Template: `{{range $i, $v := .Headers}}{{if ne $i 0}}	{{end}}{{$v}}{{end}}
{{range .Data -}}
	{{range $i, $v := .}}{{if ne $i 0}}	{{end}}{{$v}}{{end}}
{{end}}`,
	}

	if err := Convert(r, w, o); err != nil {
		t.Error(err)
	}

	expected := `column1	column2	column3
1	2	3
4	5	6
7	8	9
`
	if actual := w.String(); actual != expected {
		t.Errorf("Expectd: %s, but got %s", expected, actual)
	}
}

func TestConvertWithCustomHTMLTemplate(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,<5>,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := ConvertOption{
		HTML: true,
		Template: `{{range $i, $v := .Headers}}{{if ne $i 0}}	{{end}}{{$v}}{{end}}
{{range .Data -}}
	{{range $i, $v := .}}{{if ne $i 0}}	{{end}}{{$v}}{{end}}
{{end}}`,
	}

	if err := Convert(r, w, o); err != nil {
		t.Error(err)
	}

	expected := `aaa	bbb	ccc
1	2	3
4	&lt;5&gt;	6
7	8	9
`
	if actual := w.String(); actual != expected {
		t.Errorf("Expectd: %s, but got %s", expected, actual)
	}

}

func TestConvertWithFormatAndCustomTemplate(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := ConvertOption{
		Format: "markdown",
		Template: `{{range $i, $v := .Headers}}{{if ne $i 0}}	{{end}}{{$v}}{{end}}
{{range .Data -}}
	{{range $i, $v := .}}{{if ne $i 0}}	{{end}}{{$v}}{{end}}
{{end}}`,
	}

	if err := Convert(r, w, o); err != nil {
		t.Error(err)
	}

	expected := `aaa	bbb	ccc
1	2	3
4	5	6
7	8	9
`
	if actual := w.String(); actual != expected {
		t.Errorf("Expectd: %s, but got %s", expected, actual)
	}
}
