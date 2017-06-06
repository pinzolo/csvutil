package csvutil

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func BenchmarkSort(b *testing.B) {
	p, err := ioutil.ReadFile("testdata/bench.csv")
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		r := bytes.NewBuffer(p)
		w := &bytes.Buffer{}
		o := SortOption{
			Column: "都道府県",
		}
		Sort(r, w, o)
	}
}

func TestSortWithoutColumn(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := SortOption{
		DataType:      SortDataTypeText,
		EmptyHandling: EmptyNatural,
	}

	if err := Sort(r, w, o); err == nil {
		t.Fatal("Sort without column symbol should raise error.")
	}
}

func TestSortWithNoHeaderButColumnNotNumber(t *testing.T) {
	s := `1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := SortOption{
		NoHeader:      true,
		Column:        "aaa",
		DataType:      SortDataTypeText,
		EmptyHandling: EmptyNatural,
	}

	if err := Sort(r, w, o); err == nil {
		t.Fatal("Sort with not number column symbol for no header CSV should raise error.")
	}
}

func TestSortWithUnsupportedDataType(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := SortOption{
		Column:        "aaa",
		DataType:      "date",
		EmptyHandling: EmptyNatural,
	}

	if err := Sort(r, w, o); err == nil {
		t.Fatal("Sort with unsupported data type should raise error.")
	}
}

func TestSortWithUnsupportedEmptyHandling(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := SortOption{
		Column:        "aaa",
		DataType:      SortDataTypeText,
		EmptyHandling: "unknown",
	}

	if err := Sort(r, w, o); err == nil {
		t.Fatal("Sort with unsupported empty handling should raise error.")
	}
}

func TestSortWithBrokenCSV(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := SortOption{
		Column:        "aaa",
		DataType:      SortDataTypeText,
		EmptyHandling: EmptyNatural,
	}

	if err := Sort(r, w, o); err == nil {
		t.Fatal("Sort with broken CSV should raise error.")
	}
}

func TestSortWithBrokenHeaderCSV(t *testing.T) {
	s := `a"aa",bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := SortOption{
		Column:        "aaa",
		DataType:      SortDataTypeText,
		EmptyHandling: EmptyNatural,
	}

	if err := Sort(r, w, o); err == nil {
		t.Fatal("Sort with broken header CSV should raise error.")
	}
}

func TestSortWithEmptyCSV(t *testing.T) {
	r := bytes.NewBufferString("")
	w := &bytes.Buffer{}
	o := SortOption{
		Column:        "aaa",
		DataType:      SortDataTypeText,
		EmptyHandling: EmptyNatural,
	}

	if err := Sort(r, w, o); err == nil {
		t.Fatal("Sort with empty CSV should raise error.")
	}
}

func TestSortWithNotNumberValue(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
a,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := SortOption{
		Column:        "aaa",
		DataType:      SortDataTypeNumber,
		EmptyHandling: EmptyNatural,
	}

	if err := Sort(r, w, o); err == nil {
		t.Fatal("Sort with number data type on CSV has not number data should raise error.")
	}
}

func TestSort(t *testing.T) {
	tests := []struct {
		dataType string
		desc     bool
		empty    string
		sorted   []string
	}{
		{SortDataTypeText, false, EmptyNatural, []string{"", "1", "10", "2"}},
		{SortDataTypeText, false, EmptyFirst, []string{"", "1", "10", "2"}},
		{SortDataTypeText, false, EmptyLast, []string{"1", "10", "2", ""}},
		{SortDataTypeText, true, EmptyNatural, []string{"2", "10", "1", ""}},
		{SortDataTypeText, true, EmptyFirst, []string{"", "2", "10", "1"}},
		{SortDataTypeText, true, EmptyLast, []string{"2", "10", "1", ""}},
		{SortDataTypeNumber, false, EmptyNatural, []string{"", "1", "2", "10"}},
		{SortDataTypeNumber, false, EmptyFirst, []string{"", "1", "2", "10"}},
		{SortDataTypeNumber, false, EmptyLast, []string{"1", "2", "10", ""}},
		{SortDataTypeNumber, true, EmptyNatural, []string{"10", "2", "1", ""}},
		{SortDataTypeNumber, true, EmptyFirst, []string{"", "10", "2", "1"}},
		{SortDataTypeNumber, true, EmptyLast, []string{"10", "2", "1", ""}},
	}

	s := `aaa,bbb,ccc
2,2,3
10,5,6
,8,9
1,11,12
`

	for _, tt := range tests {
		r := bytes.NewBufferString(s)
		w := &bytes.Buffer{}
		o := SortOption{
			Column:        "aaa",
			DataType:      tt.dataType,
			Descending:    tt.desc,
			EmptyHandling: tt.empty,
		}

		Sort(r, w, o)

		result := readCSV(w.String())
		for i, s := range tt.sorted {
			if row := result[i+1]; s != row[0] {
				t.Errorf("expected %s but got %s (index: %d, test: %#v)", s, row[0], i, tt)
			}
		}
	}
}

func TestSortWithNoHeader(t *testing.T) {
	tests := []struct {
		dataType string
		desc     bool
		empty    string
		sorted   []string
	}{
		{SortDataTypeText, false, EmptyNatural, []string{"", "1", "10", "2"}},
		{SortDataTypeText, false, EmptyFirst, []string{"", "1", "10", "2"}},
		{SortDataTypeText, false, EmptyLast, []string{"1", "10", "2", ""}},
		{SortDataTypeText, true, EmptyNatural, []string{"2", "10", "1", ""}},
		{SortDataTypeText, true, EmptyFirst, []string{"", "2", "10", "1"}},
		{SortDataTypeText, true, EmptyLast, []string{"2", "10", "1", ""}},
		{SortDataTypeNumber, false, EmptyNatural, []string{"", "1", "2", "10"}},
		{SortDataTypeNumber, false, EmptyFirst, []string{"", "1", "2", "10"}},
		{SortDataTypeNumber, false, EmptyLast, []string{"1", "2", "10", ""}},
		{SortDataTypeNumber, true, EmptyNatural, []string{"10", "2", "1", ""}},
		{SortDataTypeNumber, true, EmptyFirst, []string{"", "10", "2", "1"}},
		{SortDataTypeNumber, true, EmptyLast, []string{"10", "2", "1", ""}},
	}

	s := `2,2,3
10,5,6
,8,9
1,11,12
`

	for _, tt := range tests {
		r := bytes.NewBufferString(s)
		w := &bytes.Buffer{}
		o := SortOption{
			NoHeader:      true,
			Column:        "0",
			DataType:      tt.dataType,
			Descending:    tt.desc,
			EmptyHandling: tt.empty,
		}

		Sort(r, w, o)

		result := readCSV(w.String())
		for i, s := range tt.sorted {
			if row := result[i]; s != row[0] {
				t.Errorf("expected %s but got %s (index: %d, test: %#v)", s, row[0], i, tt)
			}
		}
	}
}
