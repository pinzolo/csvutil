package csvutil

import (
	"fmt"
	"io"
	"math/rand"
	"strconv"

	gimei "github.com/mattn/go-gimei"
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
	if o.ZipCode != "" {
		return true
	}
	if o.Prefecture != "" {
		return true
	}
	if o.City != "" {
		return true
	}
	if o.Town != "" {
		return true
	}
	return false
}

func (o AddressOption) validate() error {
	if !o.hasTargetColumn() {
		return errors.New("no column")
	}
	if o.NoHeader {
		if !isDigit(o.ZipCode) {
			return errors.New("not number zip code column symbol")
		}
		if !isDigit(o.Prefecture) {
			return errors.New("not number prefecture column symbol")
		}
		if !isDigit(o.City) {
			return errors.New("not number city column symbol")
		}
		if !isDigit(o.Town) {
			return errors.New("not number town symbol")
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

func (c addrCols) indexes() []int {
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
	var hdr []string
	for {
		rec, err := cr.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return errors.Wrap(err, "cannot read csv line")
		}
		if hdr == nil && !o.NoHeader {
			hdr = rec
			cw.Write(rec)
			continue
		}
		if cols == nil {
			cols, err = setupAddressCols(o, hdr)
			if err != nil {
				return errors.Wrap(err, "column not found")
			}
		}
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
		cw.Write(newRec)
	}

	return nil
}

func setupAddressCols(o AddressOption, hdr []string) (*addrCols, error) {
	cols := &addrCols{}
	var err error
	cols.zipCode, err = newColumnWithIndex(o.ZipCode, hdr)
	if err != nil {
		return nil, err
	}
	cols.prefecture, err = newColumnWithIndex(o.Prefecture, hdr)
	if err != nil {
		return nil, err
	}
	cols.town, err = newColumnWithIndex(o.Town, hdr)
	if err != nil {
		return nil, err
	}
	cols.city, err = newColumnWithIndex(o.City, hdr)
	if err != nil {
		return nil, err
	}
	return cols, nil
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
