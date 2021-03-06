package csvutil

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

func TestNewReaderWithoutBOM(t *testing.T) {
	f, err := os.Open("testdata/utf8.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	r, bom := NewReader(f)
	if bom {
		t.Fatalf("Cannot check BOM. Expect: false, but got %v", bom)
	}
	record, err := r.Read()
	if err != nil {
		t.Fatal(err)
	}
	if record[0] != "名前" {
		t.Fatalf("Cannot read csv. Expect: 名前, but got %v", record[0])
	}
}

func TestNewReaderWithBOM(t *testing.T) {
	f, err := os.Open("testdata/utf8_with_bom.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	r, bom := NewReader(f)
	if !bom {
		t.Fatalf("Cannot check BOM. Expect: true, but got %v", bom)
	}
	record, err := r.Read()
	if err != nil {
		t.Fatal(err)
	}
	if record[0] != "名前" {
		t.Fatalf("Cannot read csv. Expect: 名前, but got %v", record[0])
	}
}

func TestNewReaderWithShortCSV(t *testing.T) {
	f, err := os.Open("testdata/short.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	r, bom := NewReader(f)
	if bom {
		t.Fatalf("Cannot check BOM. Expect: false, but got %v", bom)
	}
	record, err := r.Read()
	if err != nil {
		t.Fatal(err)
	}
	if record[0] != "a" {
		t.Fatalf("Cannot read csv. Expect: a, but got %v", record[0])
	}
}

func TestNewWriterWithoutBOM(t *testing.T) {
	b := &bytes.Buffer{}
	w := NewWriter(b, false)
	w.Write([]string{"名前", "個数"})
	w.Flush()
	if b.Bytes()[0] == utf8bom[0] {
		t.Fatal("Expect: No BOM, but got BOM in the top.")
	}
	if !strings.Contains(b.String(), "名前,個数") {
		t.Fatal("Failed writing CSV.", b.String())
	}
}

func TestNewWriterWithBOM(t *testing.T) {
	b := &bytes.Buffer{}
	w := NewWriter(b, true)
	w.Write([]string{"名前", "個数"})
	w.Flush()
	bs := b.Bytes()
	if bs[0] != utf8bom[0] && bs[1] != utf8bom[1] && bs[2] != utf8bom[2] {
		t.Fatal("Expect: With BOM, but colud not get BOM.")
	}
	if !strings.Contains(b.String(), "名前,個数") {
		t.Fatal("Failed writing CSV.", b.String())
	}
}

func TestNewReaderWithShiftJIS(t *testing.T) {
	f, err := os.Open("testdata/sjis.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	r := NewReaderWithEnc(f, japanese.ShiftJIS)
	record, err := r.Read()
	if err != nil {
		t.Fatal(err)
	}
	if record[0] != "名前" {
		t.Fatalf("Cannot read csv. Expect: 名前, but got %v", record[0])
	}
}

func TestNewReaderWithEUCJP(t *testing.T) {
	f, err := os.Open("testdata/eucjp.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	r := NewReaderWithEnc(f, japanese.EUCJP)
	record, err := r.Read()
	if err != nil {
		t.Fatal(err)
	}
	if record[0] != "名前" {
		t.Fatalf("Cannot read csv. Expect: 名前, but got %v", record[0])
	}
}

func TestNewWriterWithEnc(t *testing.T) {
	for _, e := range []encoding.Encoding{japanese.ShiftJIS, japanese.EUCJP} {
		b := &bytes.Buffer{}
		w := NewWriterWithEnc(b, e)
		w.Write([]string{"名前", "個数"})
		w.Flush()
		s, err := toUTF8(b.Bytes(), e)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(s, "名前,個数") {
			t.Fatal("Failed writing CSV.", s)
		}
	}
}

func toUTF8(bs []byte, e encoding.Encoding) (string, error) {
	b, err := ioutil.ReadAll(transform.NewReader(bytes.NewReader(bs), e.NewDecoder()))
	if err != nil {
		return "", err
	}
	return string(b), err
}
