package csvutil

import (
	"fmt"
	"io"
	"math/rand"
	"strconv"

	gimei "github.com/pinzolo/go-gimei"
	"github.com/pkg/errors"
)

var prefs = []string{
	"北海道",
	"青森県",
	"岩手県",
	"宮城県",
	"秋田県",
	"山形県",
	"福島県",
	"茨城県",
	"栃木県",
	"群馬県",
	"埼玉県",
	"千葉県",
	"東京都",
	"神奈川県",
	"新潟県",
	"富山県",
	"石川県",
	"福井県",
	"山梨県",
	"長野県",
	"岐阜県",
	"静岡県",
	"愛知県",
	"三重県",
	"滋賀県",
	"京都府",
	"大阪府",
	"兵庫県",
	"奈良県",
	"和歌山県",
	"鳥取県",
	"島根県",
	"岡山県",
	"広島県",
	"山口県",
	"徳島県",
	"香川県",
	"愛媛県",
	"高知県",
	"福岡県",
	"佐賀県",
	"長崎県",
	"熊本県",
	"大分県",
	"宮崎県",
	"鹿児島県",
	"沖縄県",
}

// AddressOption is option holder for Address.
type AddressOption struct {
	// Source file does not have header line. (default false)
	NoHeader bool
	// Encoding of source file. (default utf8)
	Encoding string
	// Encoding for output.
	OutputEncoding string
	// ZipCode column symbol.
	ZipCode string
	// Prefecture column symbol.
	Prefecture string
	// Output prefecture code.
	PrefectureCode bool
	// City column symbol.
	City string
	// Town column symbol.
	Town string
	// BlockNumber output flag.
	BlockNumber bool
	// BlockNumber width(1 or 2)
	NumberWidth int
}

func (o AddressOption) hasTargetColumn() bool {
	syms := []string{
		o.ZipCode,
		o.Prefecture,
		o.City,
		o.Town,
	}
	for _, s := range syms {
		if s != "" {
			return true
		}
	}
	return false
}

func (o AddressOption) validate() error {
	if !o.hasTargetColumn() {
		return errors.New("no column")
	}
	if o.NoHeader {
		if !isEmptyOrDigit(o.ZipCode) {
			return errors.New("not number zip code column symbol")
		}
		if !isEmptyOrDigit(o.Prefecture) {
			return errors.New("not number prefecture column symbol")
		}
		if !isEmptyOrDigit(o.City) {
			return errors.New("not number city column symbol")
		}
		if !isEmptyOrDigit(o.Town) {
			return errors.New("not number town column symbol")
		}
	}
	if o.NumberWidth != 1 && o.NumberWidth != 2 {
		return errors.New("invalid number width (1 or 2)")
	}

	return nil
}

func (o AddressOption) isFullWidthBlockNumber() bool {
	return o.NumberWidth == 2
}

func (o AddressOption) outputEncoding() string {
	if o.OutputEncoding != "" {
		return o.OutputEncoding
	}
	return o.Encoding
}

type addrCols struct {
	zipCode    *column
	prefecture *column
	city       *column
	town       *column
}

func (c *addrCols) err() error {
	cols := []*column{
		c.zipCode,
		c.prefecture,
		c.city,
		c.town,
	}
	for _, c := range cols {
		if c.err != nil {
			return c.err
		}
	}
	return nil
}

func (c *addrCols) indexes() []int {
	return []int{
		c.zipCode.index,
		c.prefecture.index,
		c.city.index,
		c.town.index,
	}
}

// Address overwrite value of given column by dummy address.
func Address(r io.Reader, w io.Writer, o AddressOption) error {
	if err := o.validate(); err != nil {
		return errors.Wrap(err, "invalid option")
	}

	cr, bom := reader(r, o.Encoding)
	cw := writer(w, bom, o.outputEncoding())
	defer cw.Flush()

	var cols *addrCols
	csvp := NewCSVProcessor(cr, cw)
	if o.NoHeader {
		csvp.SetPreBodyRead(func() error {
			cols = setupAddressCols(o, nil)
			return cols.err()
		})
	} else {
		csvp.SetHeaderHanlder(func(hdr []string) ([]string, error) {
			cols = setupAddressCols(o, hdr)
			return hdr, cols.err()
		})
	}
	csvp.SetRecordHandler(func(rec []string) ([]string, error) {
		newRec := make([]string, len(rec))
		for i, s := range rec {
			if !containsInt(cols.indexes(), i) {
				newRec[i] = s
				continue
			}

			addr := gimei.NewAddress()
			if i == cols.zipCode.index {
				newRec[i] += fakeZipCode()
			}
			if i == cols.prefecture.index {
				if o.PrefectureCode {
					newRec[i] += prefCode(addr.Prefecture.Kanji())
				} else {
					newRec[i] += addr.Prefecture.Kanji()
				}
			}
			if i == cols.city.index {
				newRec[i] += addr.City.Kanji()
			}
			if i == cols.town.index {
				newRec[i] += addr.Town.Kanji()
				if o.BlockNumber {
					newRec[i] += fakeBlockNumber(o.isFullWidthBlockNumber())
				}
			}
		}
		return newRec, nil
	})

	return csvp.Process()
}

func setupAddressCols(o AddressOption, hdr []string) *addrCols {
	cols := &addrCols{}
	cols.zipCode = newColumnWithIndex(o.ZipCode, hdr)
	cols.prefecture = newColumnWithIndex(o.Prefecture, hdr)
	cols.town = newColumnWithIndex(o.Town, hdr)
	cols.city = newColumnWithIndex(o.City, hdr)
	return cols
}

func fakeZipCode() string {
	return fmt.Sprintf("%03d-%04d", rand.Intn(1000), rand.Intn(10000))
}

func fakeBlockNumber(fullWidth bool) string {
	bn := strconv.Itoa(rand.Intn(70) + 1)
	if lot(50) {
		bn += "-" + strconv.Itoa(rand.Intn(30)+1)
	}
	if lot(25) {
		bn += "-" + strconv.Itoa(rand.Intn(10)+1)
	}
	if !fullWidth {
		return bn
	}
	return toFullWidthNum(bn)
}

func prefCode(s string) string {
	for i, p := range prefs {
		if p == s {
			return fmt.Sprintf("%02d", i+1)
		}
	}
	return ""
}
