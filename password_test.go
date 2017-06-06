package csvutil

import (
	"bytes"
	"io/ioutil"
	"testing"
)

var lowerLetters = []rune("abcdefghijklmnopqrstuvwxyz")
var upperLetters = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
var numeric = []rune("0123456789")
var specialChars = []rune(`!'@#$%^&*()_+-=[]{};:",./?`)

func BenchmarkPassword(b *testing.B) {
	p, err := ioutil.ReadFile("testdata/bench.csv")
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		r := bytes.NewBuffer(p)
		w := &bytes.Buffer{}
		o := PasswordOption{
			MinLength: 8,
			MaxLength: 16,
			Column:    "パスワード",
		}
		Password(r, w, o)
	}
}

func TestPasswordWithoutColumn(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := PasswordOption{
		MinLength: 8,
		MaxLength: 16,
	}

	if err := Password(r, w, o); err == nil {
		t.Fatal("Password without column symbol should raise error.")
	}
}

func TestPasswordWithNegativeMinLength(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := PasswordOption{
		Column:    "aaa",
		MinLength: -1,
		MaxLength: 16,
	}

	if err := Password(r, w, o); err == nil {
		t.Fatal("Password with negative min length should raise error.")
	}
}

func TestPasswordWithZeroMinLength(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := PasswordOption{
		Column:    "aaa",
		MinLength: 0,
		MaxLength: 16,
	}

	if err := Password(r, w, o); err == nil {
		t.Fatal("Password with zero min length should raise error.")
	}
}

func TestPasswordWithNegativeMaxLength(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := PasswordOption{
		MinLength: 8,
		MaxLength: -1,
	}

	if err := Password(r, w, o); err == nil {
		t.Fatal("Password with negative max length should raise error.")
	}
}

func TestPasswordWithZeroMaxLength(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := PasswordOption{
		Column:    "aaa",
		MinLength: 8,
		MaxLength: 0,
	}

	if err := Password(r, w, o); err == nil {
		t.Fatal("Password with zero max length should raise error.")
	}
}

func TestPasswordWithLessMaxLengthThanMinLength(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := PasswordOption{
		Column:    "aaa",
		MinLength: 8,
		MaxLength: 7,
	}

	if err := Password(r, w, o); err == nil {
		t.Fatal("Password with less max length than min length should raise error.")
	}
}

func TestPasswordWithNoHeaderButColumnNotNumber(t *testing.T) {
	s := `1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := PasswordOption{
		NoHeader:  true,
		Column:    "foo",
		MinLength: 8,
		MaxLength: 16,
	}

	if err := Password(r, w, o); err == nil {
		t.Fatal("Password with not number column symbol for no header CSV should raise error.")
	}

}

func TestPasswordWithUnknownColumn(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := PasswordOption{
		Column:    "ddd",
		MinLength: 8,
		MaxLength: 16,
	}

	if err := Password(r, w, o); err == nil {
		t.Fatal("Password with unknown column should raise error.")
	}
}

func TestPasswordWithBrokenCSV(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := PasswordOption{
		Column:    "aaa",
		MinLength: 8,
		MaxLength: 16,
	}

	if err := Password(r, w, o); err == nil {
		t.Fatal("Password with broken csv should raise error.")
	}
}

func TestPasswordWithBrokenHeaderCSV(t *testing.T) {
	s := `a"aa",bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := PasswordOption{
		Column:    "aaa",
		MinLength: 8,
		MaxLength: 16,
	}

	if err := Password(r, w, o); err == nil {
		t.Fatal("Password with broken header csv should raise error.")
	}
}

func TestPasswordWitColumn(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := PasswordOption{
		Column:    "aaa",
		MinLength: 8,
		MaxLength: 16,
	}

	if err := Password(r, w, o); err != nil {
		t.Fatal(err)
	}

	data := readCSV(w.String())
	if ok := allOK(data, 0, isPasswordFunc(lowerLetters, upperLetters, numeric, specialChars)); !ok {
		t.Fatalf("Password failed updating on password. %+v", data)
	}
	if ok := allOK(data, 0, isValidLengthPasswordFunc(o.MinLength, o.MaxLength)); !ok {
		t.Fatalf("Password is invalid length. %+v", data)
	}
}

func TestPasswordWithNoHeader(t *testing.T) {
	s := `1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := PasswordOption{
		Column:    "0",
		NoHeader:  true,
		MinLength: 8,
		MaxLength: 16,
	}

	if err := Password(r, w, o); err != nil {
		t.Fatal(err)
	}

	data := readCSV(w.String())
	if ok := allOK(data, 0, isPasswordFunc(lowerLetters, upperLetters, numeric, specialChars)); !ok {
		t.Fatalf("Password failed updating on mobile password. %+v", data)
	}
	if ok := allOK(data, 0, isValidLengthPasswordFunc(o.MinLength, o.MaxLength)); !ok {
		t.Fatalf("Password is invalid length. %+v", data)
	}
}

func TestPasswordWithNoNumeric(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := PasswordOption{
		Column:    "aaa",
		NoNumeric: true,
		MinLength: 8,
		MaxLength: 16,
	}

	if err := Password(r, w, o); err != nil {
		t.Fatal(err)
	}

	data := readCSV(w.String())
	if ok := allOK(data, 0, isPasswordFunc(lowerLetters, upperLetters, specialChars)); !ok {
		t.Fatalf("Password failed updating on password. %+v", data)
	}
	if ok := allOK(data, 0, isValidLengthPasswordFunc(o.MinLength, o.MaxLength)); !ok {
		t.Fatalf("Password is invalid length. %+v", data)
	}
}

func TestPasswordWitNoUpper(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := PasswordOption{
		Column:    "aaa",
		NoUpper:   true,
		MinLength: 8,
		MaxLength: 16,
	}

	if err := Password(r, w, o); err != nil {
		t.Fatal(err)
	}

	data := readCSV(w.String())
	if ok := allOK(data, 0, isPasswordFunc(lowerLetters, numeric, specialChars)); !ok {
		t.Fatalf("Password failed updating on password. %+v", data)
	}
	if ok := allOK(data, 0, isValidLengthPasswordFunc(o.MinLength, o.MaxLength)); !ok {
		t.Fatalf("Password is invalid length. %+v", data)
	}
}

func TestPasswordWitNoSpecial(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := PasswordOption{
		Column:    "aaa",
		NoSpecial: true,
		MinLength: 8,
		MaxLength: 16,
	}

	if err := Password(r, w, o); err != nil {
		t.Fatal(err)
	}

	data := readCSV(w.String())
	if ok := allOK(data, 0, isPasswordFunc(lowerLetters, upperLetters, numeric)); !ok {
		t.Fatalf("Password failed updating on password. %+v", data)
	}
	if ok := allOK(data, 0, isValidLengthPasswordFunc(o.MinLength, o.MaxLength)); !ok {
		t.Fatalf("Password is invalid length. %+v", data)
	}
}

func TestPasswordWithSameMinLengthAndMaxLength(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := PasswordOption{
		Column:    "aaa",
		NoSpecial: true,
		MinLength: 8,
		MaxLength: 8,
	}

	if err := Password(r, w, o); err != nil {
		t.Fatal(err)
	}
	data := readCSV(w.String())
	f := func(s string) bool {
		return len([]rune(s)) == 8
	}
	if ok := allOK(data, 0, f); !ok {
		t.Fatalf("Password is invalid length. %+v", data)
	}
}

func isPasswordFunc(rss ...[]rune) func(s string) bool {
	return func(s string) bool {
		for _, r := range []rune(s) {
			b := false
			for _, rs := range rss {
				if containsRune(rs, r) {
					b = true
					continue
				}
			}
			if !b {
				return false
			}
		}
		return true
	}
}

func isValidLengthPasswordFunc(min, max int) func(s string) bool {
	return func(s string) bool {
		l := len([]rune(s))
		return min <= l && l <= max
	}
}
