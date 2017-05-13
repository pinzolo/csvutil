package csvutil

import (
	"encoding/csv"
	"io"
	"math/rand"
	"strings"

	"golang.org/x/text/encoding/japanese"
)

var halfWidthNums = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "-"}
var fullWidthNums = []string{"０", "１", "２", "３", "４", "５", "６", "７", "８", "９", "－"}

func reader(r io.Reader, enc string) (*csv.Reader, bool) {
	if enc == "sjis" {
		return NewReaderWithEnc(r, japanese.ShiftJIS), false
	} else if enc == "eucjp" {
		return NewReaderWithEnc(r, japanese.EUCJP), false
	}
	return NewReader(r)
}

func writer(w io.Writer, bom bool, enc string) *csv.Writer {
	if enc == "sjis" {
		return NewWriterWithEnc(w, japanese.ShiftJIS)
	} else if enc == "eucjp" {
		return NewWriterWithEnc(w, japanese.EUCJP)
	}
	return NewWriter(w, bom)
}

func isDigit(s string) bool {
	if s == "" {
		return false
	}
	for _, b := range []byte(s) {
		if b < 0x30 || 0x39 < b {
			return false
		}
	}
	return true
}

func isEmptyOrDigit(s string) bool {
	if s == "" {
		return true
	}
	return isDigit(s)
}

func lot(n int) bool {
	if n == 100 {
		return true
	}
	if n == 0 {
		return false
	}
	return rand.Intn(100) < n
}

func sampleString(ss []string) string {
	return ss[rand.Intn(len(ss))]
}

func containsString(ss []string, s string) bool {
	for _, s2 := range ss {
		if s2 == s {
			return true
		}
	}
	return false
}

func containsInt(is []int, i int) bool {
	for _, i2 := range is {
		if i2 == i {
			return true
		}
	}
	return false
}

func containsRune(rs []rune, r rune) bool {
	for _, r2 := range rs {
		if r2 == r {
			return true
		}
	}
	return false
}

func toFullWidthNum(s string) string {
	fs := s
	for i, hn := range halfWidthNums {
		fs = strings.Replace(fs, hn, fullWidthNums[i], -1)
	}
	return fs
}
