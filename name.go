package csvutil

import (
	"io"
	"strings"

	gimei "github.com/pinzolo/go-gimei"
	"github.com/pkg/errors"
)

var supportedGenderFormats = []string{"code", "en_short", "en_long", "jp_short", "jp_long", "symbol"}
var maleGenders = []string{"1", "M", "Male", "男", "男性", "♂"}
var femaleGenders = []string{"2", "F", "Female", "女", "女性", "♀"}

// NameOption is option holder for Name.
type NameOption struct {
	// Source file does not have header line. (default false)
	NoHeader bool
	// Encoding of source file. (default utf8)
	Encoding string
	// Encoding for output
	OutputEncoding string
	// ZipCode column symbol
	ZipCode string
	// Full name column symbol
	Name string
	// First name column symbol
	FirstName string
	// Last name column symbol
	LastName string
	// Kana of full name column symbol
	Kana string
	// Kana of first name column symbol
	FirstKana string
	// Kana of last name column symbol
	LastKana string
	// Output hiragana as kana
	Hiragana bool
	// Gender column symbol
	Gender string
	// Gender format
	GenderFormat string
	// Rate of male output
	MaleRate int
	// Reference column symbol
	Reference string
	// Ignore reference error
	RistrictReference bool
	// Delimiter space width
	SpaceWidth int
}

func (o NameOption) hasTargetColumn() bool {
	if o.Name != "" {
		return true
	}
	if o.FirstName != "" {
		return true
	}
	if o.LastName != "" {
		return true
	}
	if o.Kana != "" {
		return true
	}
	if o.FirstKana != "" {
		return true
	}
	if o.LastKana != "" {
		return true
	}
	if o.Gender != "" {
		return true
	}
	return false
}

func (o NameOption) validate() error {
	if !o.hasTargetColumn() {
		return errors.New("no column")
	}
	if o.NoHeader {
		if !isEmptyOrDigit(o.Name) {
			return errors.New("not number name column symbol")
		}
		if !isEmptyOrDigit(o.FirstName) {
			return errors.New("not number first name column symbol")
		}
		if !isEmptyOrDigit(o.LastName) {
			return errors.New("not number last name column symbol")
		}
		if !isEmptyOrDigit(o.Kana) {
			return errors.New("not number kana column symbol")
		}
		if !isEmptyOrDigit(o.FirstKana) {
			return errors.New("not number first kana column symbol")
		}
		if !isEmptyOrDigit(o.LastKana) {
			return errors.New("not number last kana column symbol")
		}
		if !isEmptyOrDigit(o.Gender) {
			return errors.New("not number gender column symbol")
		}
		if !isEmptyOrDigit(o.Reference) {
			return errors.New("not number reference column symbol")
		}
	}
	if o.SpaceWidth < 0 || 2 < o.SpaceWidth {
		return errors.New("invalid space width (0 or 1 or 2)")
	}
	if o.Gender != "" && !containsString(supportedGenderFormats, o.GenderFormat) {
		return errors.Errorf("unsupported gender format: %s", o.GenderFormat)
	}

	return nil
}

func (o NameOption) space() string {
	var sp string
	if o.SpaceWidth == 1 {
		return " "
	} else if o.SpaceWidth == 2 {
		return "　"
	}
	return sp
}

func (o NameOption) outputEncoding() string {
	if o.OutputEncoding != "" {
		return o.OutputEncoding
	}
	return o.Encoding
}

type nameCols struct {
	name      *column
	firstName *column
	lastName  *column
	kana      *column
	firstKana *column
	lastKana  *column
	gender    *column
	reference *column
}

func (c *nameCols) err() error {
	cols := []*column{
		c.name,
		c.firstName,
		c.lastName,
		c.kana,
		c.firstKana,
		c.lastKana,
		c.gender,
		c.reference,
	}
	for _, col := range cols {
		if col.err != nil {
			return col.err
		}
	}
	return nil
}

func (c *nameCols) indexes() []int {
	return []int{
		c.name.index,
		c.firstName.index,
		c.lastName.index,
		c.kana.index,
		c.firstKana.index,
		c.lastKana.index,
		c.gender.index,
	}
}

// Name overwrite value of given column by dummy name.
func Name(r io.Reader, w io.Writer, o NameOption) error {
	if err := o.validate(); err != nil {
		return errors.Wrap(err, "invalid option")
	}

	cr, bom := reader(r, o.Encoding)
	cw := writer(w, bom, o.outputEncoding())
	defer cw.Flush()

	var cols *nameCols
	csvp := NewCSVProcessor(cr, cw)
	if o.NoHeader {
		csvp.SetPreBodyRead(func() error {
			cols = setupNameCols(o, nil)
			return cols.err()
		})
	} else {
		csvp.SetHeaderHanlder(func(hdr []string) ([]string, error) {
			cols = setupNameCols(o, hdr)
			return hdr, cols.err()
		})
	}
	csvp.SetRecordHandler(func(rec []string) ([]string, error) {
		name, err := fakeName(rec, o, cols)
		if err != nil {
			if o.RistrictReference {
				return nil, err
			}
		}
		newRec := make([]string, len(rec))
		for i, s := range rec {
			if !containsInt(cols.indexes(), i) {
				newRec[i] = s
				continue
			}

			if name == nil {
				continue
			}

			if i == cols.name.index {
				newRec[i] = name.Last.Kanji() + o.space() + name.First.Kanji()
			} else if i == cols.firstName.index {
				newRec[i] = name.First.Kanji()
			} else if i == cols.lastName.index {
				newRec[i] = name.Last.Kanji()
			} else if i == cols.kana.index {
				if o.Hiragana {
					newRec[i] = name.Last.Hiragana() + o.space() + name.First.Hiragana()
				} else {
					newRec[i] = name.Last.Katakana() + o.space() + name.First.Katakana()
				}
			} else if i == cols.firstKana.index {
				if o.Hiragana {
					newRec[i] = name.First.Hiragana()
				} else {
					newRec[i] = name.First.Katakana()
				}
			} else if i == cols.lastKana.index {
				if o.Hiragana {
					newRec[i] = name.Last.Hiragana()
				} else {
					newRec[i] = name.Last.Katakana()
				}
			} else if i == cols.gender.index {
				genders := maleGenders
				if name.IsFemale() {
					genders = femaleGenders
				}
				for fi, f := range supportedGenderFormats {
					if f == o.GenderFormat {
						newRec[i] = genders[fi]
					}
				}
			}
		}
		return newRec, nil
	})

	return csvp.Process()
}

func setupNameCols(o NameOption, hdr []string) *nameCols {
	cols := &nameCols{}
	cols.name = newColumnWithIndex(o.Name, hdr)
	cols.firstName = newColumnWithIndex(o.FirstName, hdr)
	cols.lastName = newColumnWithIndex(o.LastName, hdr)
	cols.kana = newColumnWithIndex(o.Kana, hdr)
	cols.firstKana = newColumnWithIndex(o.FirstKana, hdr)
	cols.lastKana = newColumnWithIndex(o.LastKana, hdr)
	cols.gender = newColumnWithIndex(o.Gender, hdr)
	cols.reference = newColumnWithIndex(o.Reference, hdr)
	return cols
}

func fakeName(rec []string, o NameOption, cols *nameCols) (*gimei.Name, error) {
	if lot(o.MaleRate) {
		if o.Reference == "" {
			return gimei.NewMale(), nil
		}
		return gimei.NewMaleByLastName(getReferenceLastName(rec[cols.reference.index]))
	}
	if o.Reference == "" {
		return gimei.NewFemale(), nil
	}
	return gimei.NewFemaleByLastName(getReferenceLastName(rec[cols.reference.index]))
}

func getReferenceLastName(n string) string {
	if strings.Contains(n, " ") {
		return strings.Split(n, " ")[0]
	} else if strings.Contains(n, "　") {
		return strings.Split(n, "　")[0]
	}
	return n
}
