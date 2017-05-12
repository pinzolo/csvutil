package csvutil

import (
	"bytes"
	"testing"
)

func TestCollectWithoutColumn(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := CollectOption{}

	if err := Collect(r, w, o); err == nil {
		t.Error("Collect without column symbol should raise error.")
	}
}

func TestCollectWithNoHeaderButColumnNotNumber(t *testing.T) {
	s := `1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := CollectOption{
		NoHeader: true,
		Column:   "aaa",
	}

	if err := Collect(r, w, o); err == nil {
		t.Error("Collect with not number column symbol for no header CSV should raise error.")
	}
}

func TestCollectWithBrokenCSV(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := CollectOption{
		Column: "aaa",
	}

	if err := Collect(r, w, o); err == nil {
		t.Error("Collect with broken CSV should raise error.")
	}
}

func TestCollectWithBrokenHeaderCSV(t *testing.T) {
	s := `a"aa",bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := CollectOption{
		Column: "aaa",
	}

	if err := Collect(r, w, o); err == nil {
		t.Error("Collect with broken header CSV should raise error.")
	}
}

func TestCollectWithUnsupportedSortKey(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := CollectOption{
		Column:  "aaa",
		Sort:    true,
		SortKey: "foo",
	}

	if err := Collect(r, w, o); err == nil {
		t.Error("Collect with unsupported sort key should raise error.")
	}
}

func TestCollectWithUnsupportedSortDirection(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := CollectOption{
		Column:        "aaa",
		Sort:          true,
		SortKey:       "count",
		SortDirection: "foo",
	}

	if err := Collect(r, w, o); err == nil {
		t.Error("Collect with unsupported sort direction should raise error.")
	}
}

func TestCollectWithoutSort(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
2,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := CollectOption{
		Column: "aaa",
	}

	if err := Collect(r, w, o); err != nil {
		t.Error(err)
	}

	expected := `1
4
2
`
	if actual := w.String(); actual != expected {
		t.Errorf("Expectd: %s, but got %s", expected, actual)
	}
}

func TestCollectWitNoHeader(t *testing.T) {
	s := `1,2,3
4,5,6
2,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := CollectOption{
		Column:   "0",
		NoHeader: true,
	}

	if err := Collect(r, w, o); err != nil {
		t.Error(err)
	}

	expected := `1
4
2
`
	if actual := w.String(); actual != expected {
		t.Errorf("Expectd: %s, but got %s", expected, actual)
	}
}

func TestCollectOnDuplicateValues(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
1,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := CollectOption{
		Column: "aaa",
	}

	if err := Collect(r, w, o); err != nil {
		t.Error(err)
	}

	expected := `1
4
`
	if actual := w.String(); actual != expected {
		t.Errorf("Expectd: %s, but got %s", expected, actual)
	}
}

func TestCollectOnIgnoreEmpty(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := CollectOption{
		Column: "aaa",
	}

	if err := Collect(r, w, o); err != nil {
		t.Error(err)
	}

	expected := `1
4
`
	if actual := w.String(); actual != expected {
		t.Errorf("Expectd: %s, but got %s", expected, actual)
	}
}

func TestCollectWithAllowEmpty(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := CollectOption{
		Column:     "aaa",
		AllowEmpty: true,
	}

	if err := Collect(r, w, o); err != nil {
		t.Error(err)
	}

	expected := `1
4

`
	if actual := w.String(); actual != expected {
		t.Errorf("Expectd: %s, but got %s", expected, actual)
	}
}

func TestCollectWithPrintCount(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
1,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := CollectOption{
		Column:     "aaa",
		PrintCount: true,
	}

	if err := Collect(r, w, o); err != nil {
		t.Error(err)
	}

	expected := `2	1
1	4
`
	if actual := w.String(); actual != expected {
		t.Errorf("Expectd: %s, but got %s", expected, actual)
	}
}

func TestCollectWithSortByValueAsc(t *testing.T) {
	s := `aaa,bbb,ccc
3,2,3
4,5,6
1,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := CollectOption{
		Column:        "aaa",
		Sort:          true,
		SortKey:       "value",
		SortDirection: "asc",
	}

	if err := Collect(r, w, o); err != nil {
		t.Error(err)
	}

	expected := `1
3
4
`
	if actual := w.String(); actual != expected {
		t.Errorf("Expectd: %s, but got %s", expected, actual)
	}
}

func TestCollectWithSortByValueDesc(t *testing.T) {
	s := `aaa,bbb,ccc
3,2,3
4,5,6
1,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := CollectOption{
		Column:        "aaa",
		Sort:          true,
		SortKey:       "value",
		SortDirection: "desc",
	}

	if err := Collect(r, w, o); err != nil {
		t.Error(err)
	}

	expected := `4
3
1
`
	if actual := w.String(); actual != expected {
		t.Errorf("Expectd: %s, but got %s", expected, actual)
	}
}

func TestCollectWithSortByCountAsc(t *testing.T) {
	s := `aaa,bbb,ccc
3,2,3
4,5,6
1,8,9
1,8,9
1,8,9
4,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := CollectOption{
		Column:        "aaa",
		PrintCount:    true,
		Sort:          true,
		SortKey:       "count",
		SortDirection: "asc",
	}

	if err := Collect(r, w, o); err != nil {
		t.Error(err)
	}

	expected := `1	3
2	4
3	1
`
	if actual := w.String(); actual != expected {
		t.Errorf("Expectd: %s, but got %s", expected, actual)
	}
}

func TestCollectWithSortByCountDesc(t *testing.T) {
	s := `aaa,bbb,ccc
3,2,3
4,5,6
1,8,9
1,8,9
1,8,9
4,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := CollectOption{
		Column:        "aaa",
		PrintCount:    true,
		Sort:          true,
		SortKey:       "count",
		SortDirection: "desc",
	}

	if err := Collect(r, w, o); err != nil {
		t.Error(err)
	}

	expected := `3	1
2	4
1	3
`
	if actual := w.String(); actual != expected {
		t.Errorf("Expectd: %s, but got %s", expected, actual)
	}
}
