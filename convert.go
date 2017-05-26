package csvutil

import (
	"bufio"
	"encoding/json"
	"encoding/xml"
	htmpl "html/template"
	"io"
	"strconv"
	txtmpl "text/template"

	yaml "gopkg.in/yaml.v2"

	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
)

var supportedFormats = []string{"markdown", "json", "yaml", "html", "xml"}

type convertContext struct {
	Headers []string
	Data    [][]string
}

func (ctx convertContext) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = "csv"
	e.EncodeToken(start)
	for _, l := range ctx.Data {
		de := xml.StartElement{
			Name: xml.Name{Local: "data"},
		}
		e.EncodeToken(de)
		for i, h := range ctx.Headers {
			se := xml.StartElement{
				Name: xml.Name{Local: h},
			}
			e.EncodeElement(l[i], se)
		}
		e.EncodeToken(de.End())
	}
	e.EncodeToken(start.End())
	return nil
}

// ConvertOption is option holder for Convert.
type ConvertOption struct {
	// Source file does not have header line. (default false)
	NoHeader bool
	// Encoding of source file. (default utf8)
	Encoding string
	// Format of converter
	Format string
	// Template for convert
	Template string
	// HTML custom template is given
	HTML bool
}

func (o ConvertOption) validate() error {
	if o.Format != "" && !containsString(supportedFormats, o.Format) {
		return errors.Errorf("unsupported format: %s", o.Format)
	}
	if o.Format == "" && o.Template == "" {
		return errors.New("required format or template")
	}
	return nil
}

// Convert CSV.
func Convert(r io.Reader, w io.Writer, o ConvertOption) error {
	if err := o.validate(); err != nil {
		return errors.Wrap(err, "invalid option")
	}

	cr, _ := reader(r, o.Encoding)
	bw := bufio.NewWriter(w)
	defer bw.Flush()

	recs, err := cr.ReadAll()
	if err != nil {
		return err
	}
	if len(recs) == 0 {
		bw.WriteString("")
		return nil
	}

	ctx := convertContext{
		Headers: make([]string, len(recs[0])),
		Data:    nil,
	}
	if o.NoHeader {
		for i := 0; i < len(recs[0]); i++ {
			ctx.Headers[i] = "column" + strconv.Itoa(i+1)
		}
		ctx.Data = recs
	} else {
		for i, h := range recs[0] {
			ctx.Headers[i] = h
		}
		ctx.Data = recs[1:]
	}

	var p []byte
	if o.Template != "" {
		if o.HTML {
			var t *htmpl.Template
			if t, err = htmpl.New("text").Parse(o.Template); err == nil {
				err = t.Execute(bw, ctx)
			}
		} else {
			var t *txtmpl.Template
			if t, err = txtmpl.New("text").Parse(o.Template); err == nil {
				err = t.Execute(bw, ctx)
			}
		}
	} else if o.Format == "markdown" {
		tw := tablewriter.NewWriter(w)
		tw.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
		tw.SetCenterSeparator("|")
		tw.SetAutoFormatHeaders(false)
		tw.SetHeader(ctx.Headers)
		tw.AppendBulk(ctx.Data)
		tw.Render()
	} else if o.Format == "json" {
		p, err = json.MarshalIndent(marshalee(ctx), "", "\t")
	} else if o.Format == "yaml" {
		p, err = yaml.Marshal(marshalee(ctx))
	} else if o.Format == "html" {
		var t *htmpl.Template
		if t, err = htmpl.New("html").Parse(htmlTmpl()); err == nil {
			err = t.Execute(bw, ctx)
		}
	} else if o.Format == "xml" {
		p, err = xml.MarshalIndent(ctx, "", "\t")
	}
	if err != nil {
		return err
	}
	_, err = bw.Write(appendLastNewLine(p))
	return err
}

func htmlTmpl() string {
	return `<table>
	<thead>
		<tr>
{{- range .Headers}}
			<th>{{.}}</th>
{{- end}}
		</tr>
	</thead>
	<tbody>
{{- range .Data}}
		<tr>
	{{- range .}}
			<td>{{.}}</td>
	{{- end}}
		</tr>
{{- end}}
	</tbody>
</table>
`
}

func marshalee(ctx convertContext) []map[string]string {
	ret := make([]map[string]string, len(ctx.Data))
	for i, l := range ctx.Data {
		m := make(map[string]string)
		for j, h := range ctx.Headers {
			m[h] = l[j]
		}
		ret[i] = m
	}
	return ret
}

func appendLastNewLine(p []byte) []byte {
	if p == nil {
		return p
	}

	if p[len(p)-1] == 10 {
		return p
	}

	return append(p, 10)
}
