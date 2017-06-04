package csvutil

import (
	"bytes"
	"io/ioutil"
	"strconv"
	"strings"
	"testing"
)

func BenchmarkNumeric(b *testing.B) {
	p, err := ioutil.ReadFile("testdata/bench.csv")
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		r := bytes.NewBuffer(p)
		w := &bytes.Buffer{}
		o := NumericOption{
			Column: "パスワード",
			Max:    1000000,
			Min:    0,
		}
		Numeric(r, w, o)
	}
}

func TestNumericWithoutColumn(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := NumericOption{
		Max: 10000,
		Min: -10000,
	}

	if err := Numeric(r, w, o); err == nil {
		t.Fatal("Numeric without column symbol should raise error.")
	}
}

func TestNumericWithMaxLessThanMin(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := NumericOption{
		Column: "aaa",
		Max:    -10000,
		Min:    10000,
	}

	if err := Numeric(r, w, o); err == nil {
		t.Fatal("Numeric with max that is less than min should raise error.")
	}
}

func TestNumericWithZeroDecimalDigit(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := NumericOption{
		Column:       "aaa",
		Max:          10000,
		Min:          -10000,
		Decimal:      true,
		DecimalDigit: 0,
	}

	if err := Numeric(r, w, o); err == nil {
		t.Fatal("Numeric with zero decimal digit should raise error.")
	}
}

func TestNumericWithNegativeDecimalDigit(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := NumericOption{
		Column:       "aaa",
		Max:          10000,
		Min:          -10000,
		Decimal:      true,
		DecimalDigit: -1,
	}

	if err := Numeric(r, w, o); err == nil {
		t.Fatal("Numeric with negative decimal digit should raise error.")
	}
}

func TestNumericWithNoHeaderButColumnNotNumber(t *testing.T) {
	s := `1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := NumericOption{
		NoHeader: true,
		Column:   "foo",
		Max:      10000,
		Min:      -10000,
	}

	if err := Numeric(r, w, o); err == nil {
		t.Fatal("Numeric with not number column symbol for no header CSV should raise error.")
	}
}

func TestNumericWithUnknownColumn(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := NumericOption{
		Column: "ddd",
		Max:    10000,
		Min:    -10000,
	}

	if err := Numeric(r, w, o); err == nil {
		t.Fatal("Numeric with unknown column should raise error.")
	}
}

func TestNumericWithBrokenCSV(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := NumericOption{
		Column: "aaa",
		Max:    10000,
		Min:    -10000,
	}

	if err := Numeric(r, w, o); err == nil {
		t.Fatal("Numeric with broken csv should raise error.")
	}
}

func TestNumericInteger(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := NumericOption{
		Column: "aaa",
		Max:    10000,
		Min:    0,
	}

	if err := Numeric(r, w, o); err != nil {
		t.Fatal(err)
	}

	data := readCSV(w.String())
	f := func(s string) bool {
		if !isDigit(s) {
			return false
		}
		if i, err := strconv.Atoi(s); i >= 10000 || err != nil {
			return false
		}
		return true
	}
	if ok := allOK(data, 0, f); !ok {
		t.Fatalf("Numeric failed updating with integer. %+v", data)
	}
}

func TestNumericNegativeInteger(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := NumericOption{
		Column: "aaa",
		Max:    0,
		Min:    -10000,
	}

	if err := Numeric(r, w, o); err != nil {
		t.Fatal(err)
	}

	data := readCSV(w.String())
	f := func(s string) bool {
		if !strings.HasPrefix(s, "-") {
			return false
		}
		abs := string([]rune(s)[1:])
		return isDigit(abs)
	}
	if ok := allOK(data, 0, f); !ok {
		t.Fatalf("Numeric failed updating with negative integer. %+v", data)
	}
}

func TestNumericWithNoHeader(t *testing.T) {
	s := `1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := NumericOption{
		NoHeader: true,
		Column:   "0",
		Max:      10000,
		Min:      -10000,
	}

	if err := Numeric(r, w, o); err != nil {
		t.Fatal(err)
	}

	data := readCSV(w.String())
	f := func(s string) bool {
		if i, err := strconv.Atoi(s); err == nil && -10000 <= i && i < 10000 {
			return true
		}
		return false
	}
	if ok := allOKNoHeader(data, 0, f); !ok {
		t.Fatalf("Numeric failed updating with integer. %+v", data)
	}
}

func TestNumericWithDecimal(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := NumericOption{
		Column:       "aaa",
		Max:          100,
		Min:          0,
		Decimal:      true,
		DecimalDigit: 3,
	}

	if err := Numeric(r, w, o); err != nil {
		t.Fatal(err)
	}

	data := readCSV(w.String())
	f := func(s string) bool {
		if !strings.Contains(s, ".") {
			t.Errorf("Not contains period: %s", s)
			return false
		}
		if f, err := strconv.ParseFloat(s, 32); err != nil || f < 0.0 {
			t.Errorf("Less than min: %s", s)
			return false
		}
		if f, err := strconv.ParseFloat(s, 32); err != nil || 100.0 < f {
			t.Errorf("Greater than max: %s", s)
			return false
		}
		if len(strings.Split(s, ".")) != 2 {
			t.Errorf("Invalid decimal format: %s", s)
			return false
		}
		if len(strings.Split(s, ".")[1]) != 3 {
			t.Errorf("Invalid decimal digits (%d): %s", o.DecimalDigit, s)
			return false
		}
		return true
	}
	if ok := allOK(data, 0, f); !ok {
		t.Fatalf("Numeric failed updating with decimal. %+v", data)
	}
}

func TestNumericWithNegativeDecimal(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := NumericOption{
		Column:       "aaa",
		Max:          0,
		Min:          -100,
		Decimal:      true,
		DecimalDigit: 3,
	}

	if err := Numeric(r, w, o); err != nil {
		t.Fatal(err)
	}

	data := readCSV(w.String())
	f := func(s string) bool {
		if !strings.Contains(s, ".") {
			t.Errorf("Not contains period: %s", s)
			return false
		}
		if !strings.HasPrefix(s, "-") {
			t.Errorf("Not negative decimal: %s", s)
			return false
		}
		if f, err := strconv.ParseFloat(s, 32); err != nil || f < -100.0 {
			t.Errorf("Less than min: %s", s)
			return false
		}
		if f, err := strconv.ParseFloat(s, 32); err != nil || 0.0 < f {
			t.Errorf("Greater than max: %s", s)
			return false
		}
		if len(strings.Split(s, ".")) != 2 {
			t.Errorf("Invalid decimal format: %s", s)
			return false
		}
		if len(strings.Split(s, ".")[1]) != 3 {
			t.Errorf("Invalid decimal digits (%d): %s", o.DecimalDigit, s)
			return false
		}
		return true
	}
	if ok := allOK(data, 0, f); !ok {
		t.Fatalf("Numeric failed updating with negative decimal. %+v", data)
	}
}

func TestNumericWithDecimalLessThanOne(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := NumericOption{
		Column:       "aaa",
		Max:          1,
		Min:          0,
		Decimal:      true,
		DecimalDigit: 3,
	}

	if err := Numeric(r, w, o); err != nil {
		t.Fatal(err)
	}

	data := readCSV(w.String())
	f := func(s string) bool {
		if !strings.Contains(s, ".") {
			t.Errorf("Not contains period: %s", s)
			return false
		}
		if !strings.HasPrefix(s, "0.") {
			t.Errorf("Integer part of %s should be zero", s)
			return false
		}
		if f, err := strconv.ParseFloat(s, 32); err != nil || f < 0.0 {
			t.Errorf("Less than min: %s", s)
			return false
		}
		if f, err := strconv.ParseFloat(s, 32); err != nil || 100.0 < f {
			t.Errorf("Greater than max: %s", s)
			return false
		}
		if len(strings.Split(s, ".")) != 2 {
			t.Errorf("Invalid decimal format: %s", s)
			return false
		}
		if len(strings.Split(s, ".")[1]) != 3 {
			t.Errorf("Invalid decimal digits (%d): %s", o.DecimalDigit, s)
			return false
		}
		return true
	}
	if ok := allOK(data, 0, f); !ok {
		t.Fatalf("Numeric failed updating with decimal. %+v", data)
	}
}

func TestNumericWithNegativeDecimalGreaterThanMinusOne(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := NumericOption{
		Column:       "aaa",
		Max:          0,
		Min:          -1,
		Decimal:      true,
		DecimalDigit: 3,
	}

	if err := Numeric(r, w, o); err != nil {
		t.Fatal(err)
	}

	data := readCSV(w.String())
	f := func(s string) bool {
		if !strings.Contains(s, ".") {
			t.Errorf("Not contains period: %s", s)
			return false
		}
		if !strings.HasPrefix(s, "-0.") {
			t.Errorf("Not negative decimal or integer part is not zero: %s", s)
			return false
		}
		if f, err := strconv.ParseFloat(s, 32); err != nil || f < -100.0 {
			t.Errorf("Less than min: %s", s)
			return false
		}
		if f, err := strconv.ParseFloat(s, 32); err != nil || 0.0 < f {
			t.Errorf("Greater than max: %s", s)
			return false
		}
		if len(strings.Split(s, ".")) != 2 {
			t.Errorf("Invalid decimal format: %s", s)
			return false
		}
		if len(strings.Split(s, ".")[1]) != 3 {
			t.Errorf("Invalid decimal digits (%d): %s", o.DecimalDigit, s)
			return false
		}
		return true
	}
	if ok := allOK(data, 0, f); !ok {
		t.Fatalf("Numeric failed updating with negative decimal. %+v", data)
	}
}
