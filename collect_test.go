package csvutil

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func BenchmarkCollect(b *testing.B) {
	p, err := ioutil.ReadFile("testdata/bench.csv")
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		r := bytes.NewBuffer(p)
		w := &bytes.Buffer{}
		o := CollectOption{
			Column: "都道府県",
		}
		Collect(r, w, o)
	}
}

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
		t.Fatal("Collect without column symbol should raise error.")
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
		t.Fatal("Collect with not number column symbol for no header CSV should raise error.")
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
		t.Fatal("Collect with broken CSV should raise error.")
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
		t.Fatal("Collect with broken header CSV should raise error.")
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
		t.Fatal("Collect with unsupported sort key should raise error.")
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
		t.Fatal(err)
	}

	expected := `1
4
2
`
	if actual := w.String(); actual != expected {
		t.Fatalf("Expectd: %s, but got %s", expected, actual)
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
		t.Fatal(err)
	}

	expected := `1
4
2
`
	if actual := w.String(); actual != expected {
		t.Fatalf("Expectd: %s, but got %s", expected, actual)
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
		t.Fatal(err)
	}

	expected := `1
4
`
	if actual := w.String(); actual != expected {
		t.Fatalf("Expectd: %s, but got %s", expected, actual)
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
		t.Fatal(err)
	}

	expected := `1
4
`
	if actual := w.String(); actual != expected {
		t.Fatalf("Expectd: %s, but got %s", expected, actual)
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
		t.Fatal(err)
	}

	expected := `1
4

`
	if actual := w.String(); actual != expected {
		t.Fatalf("Expectd: %s, but got %s", expected, actual)
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
		t.Fatal(err)
	}

	expected := `2	1
1	4
`
	if actual := w.String(); actual != expected {
		t.Fatalf("Expectd: %s, but got %s", expected, actual)
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
		Column:  "aaa",
		Sort:    true,
		SortKey: "value",
	}

	if err := Collect(r, w, o); err != nil {
		t.Fatal(err)
	}

	expected := `1
3
4
`
	if actual := w.String(); actual != expected {
		t.Fatalf("Expectd: %s, but got %s", expected, actual)
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
		Column:     "aaa",
		Sort:       true,
		SortKey:    "value",
		Descending: true,
	}

	if err := Collect(r, w, o); err != nil {
		t.Fatal(err)
	}

	expected := `4
3
1
`
	if actual := w.String(); actual != expected {
		t.Fatalf("Expectd: %s, but got %s", expected, actual)
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
		Column:     "aaa",
		PrintCount: true,
		Sort:       true,
		SortKey:    "count",
	}

	if err := Collect(r, w, o); err != nil {
		t.Fatal(err)
	}

	expected := `1	3
2	4
3	1
`
	if actual := w.String(); actual != expected {
		t.Fatalf("Expectd: %s, but got %s", expected, actual)
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
		Column:     "aaa",
		PrintCount: true,
		Sort:       true,
		SortKey:    "count",
		Descending: true,
	}

	if err := Collect(r, w, o); err != nil {
		t.Fatal(err)
	}

	expected := `3	1
2	4
1	3
`
	if actual := w.String(); actual != expected {
		t.Fatalf("Expectd: %s, but got %s", expected, actual)
	}
}
