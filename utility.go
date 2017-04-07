package csvutil

import (
	"encoding/csv"
	"io"
	"math/rand"

	"golang.org/x/text/encoding/japanese"
)

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

func lot(n int) bool {
	if n == 100 {
		return true
	}
	if n == 0 {
		return false
	}
	return rand.Intn(100) < n
}
