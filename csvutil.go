package csvutil

import (
	"bufio"
	"encoding/csv"
	"io"

	"golang.org/x/text/encoding"
	"golang.org/x/text/transform"
)

var utf8bom = [...]byte{0xEF, 0xBB, 0xBF}

// UTF8BOM is bytes for byte order mark of UTF8.
func UTF8BOM() []byte {
	return utf8bom[:]
}

// NewReader returns *csv.Reader for UTF8.
// If source has BOM, returns true as second return value.
func NewReader(r io.Reader) (*csv.Reader, bool) {
	var (
		br  *bufio.Reader
		ok  bool
		bom bool
	)
	if br, ok = r.(*bufio.Reader); !ok {
		br = bufio.NewReader(r)
	}

	l := len(utf8bom)
	bs, err := br.Peek(l)
	if err != nil {
		return csv.NewReader(br), false
	}
	if isSameBytes(bs[:l], UTF8BOM()) {
		br.Discard(l)
		bom = true
	}

	return csv.NewReader(br), bom
}

// NewWriter returns *csv.Writer for UTF8.
// If true as bom, Writer writes BOM at the top.
func NewWriter(w io.Writer, bom bool) *csv.Writer {
	var (
		bw *bufio.Writer
		ok bool
	)
	if bw, ok = w.(*bufio.Writer); !ok {
		bw = bufio.NewWriter(w)
	}

	if bom {
		bw.Write(UTF8BOM())
	}
	return csv.NewWriter(bw)
}

// NewReaderWithEnc returns *csv.Reader for given encoding.
func NewReaderWithEnc(r io.Reader, e encoding.Encoding) *csv.Reader {
	return csv.NewReader(transform.NewReader(r, e.NewDecoder()))
}

// NewWriterWithEnc returns *csv.Writer for given encoding.
func NewWriterWithEnc(w io.Writer, e encoding.Encoding) *csv.Writer {
	return csv.NewWriter(transform.NewWriter(w, e.NewEncoder()))
}

func isSameBytes(bs1 []byte, bs2 []byte) bool {
	for i, b := range bs1 {
		if b != bs2[i] {
			return false
		}
	}
	return true
}
