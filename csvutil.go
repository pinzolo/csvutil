package csvutil

import (
	"bufio"
	"encoding/csv"
	"io"

	"golang.org/x/text/encoding"
	"golang.org/x/text/transform"
)

// UTF8BOM is bytes for byte order mark of UTF8.
var UTF8BOM = [3]byte{0xEF, 0xBB, 0xBF}

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

	bs, err := br.Peek(len(UTF8BOM))
	if err != nil {
		return csv.NewReader(br), false
	}
	if bs[0] == UTF8BOM[0] && bs[1] == UTF8BOM[1] && bs[2] == UTF8BOM[2] {
		br.Discard(len(UTF8BOM))
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
		bw.Write(UTF8BOM[:])
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
