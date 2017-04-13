package csvutil

import (
	"bytes"
	"strings"
	"testing"
	"unicode"
)

func TestNameWithoutTargetColumns(t *testing.T) {
	r := &bytes.Buffer{}
	w := &bytes.Buffer{}
	o := NameOption{}
	if err := Name(r, w, o); err == nil {
		t.Error("Name without size should raise error.")
	}
}

func TestNameWithNoHeaderAndNotDigitName(t *testing.T) {
	s := `1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := NameOption{
		NoHeader: true,
		Name:     "aaa",
	}

	if err := Name(r, w, o); err == nil {
		t.Error("Name with no header and not digit name column symbol should raise error.")
	}
}

func TestNameWithNoHeaderAndNotDigitFirstName(t *testing.T) {
	s := `1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := NameOption{
		NoHeader:  true,
		FirstName: "aaa",
	}

	if err := Name(r, w, o); err == nil {
		t.Error("Name with no header and not digit first name column symbol should raise error.")
	}
}

func TestNameWithNoHeaderAndNotDigitLastName(t *testing.T) {
	s := `1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := NameOption{
		NoHeader: true,
		LastName: "aaa",
	}

	if err := Name(r, w, o); err == nil {
		t.Error("Name with no header and not digit last name column symbol should raise error.")
	}
}

func TestNameWithNoHeaderAndNotDigitKana(t *testing.T) {
	s := `1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := NameOption{
		NoHeader: true,
		Kana:     "aaa",
	}

	if err := Name(r, w, o); err == nil {
		t.Error("Kana with no header and not digit kana column symbol should raise error.")
	}
}

func TestNameWithNoHeaderAndNotDigitFirstKana(t *testing.T) {
	s := `1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := NameOption{
		NoHeader:  true,
		FirstKana: "aaa",
	}

	if err := Name(r, w, o); err == nil {
		t.Error("Name with no header and not digit first kana column symbol should raise error.")
	}
}

func TestNameWithNoHeaderAndNotDigitLastKana(t *testing.T) {
	s := `1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := NameOption{
		NoHeader: true,
		LastKana: "aaa",
	}

	if err := Name(r, w, o); err == nil {
		t.Error("Name with no header and not digit last kana column symbol should raise error.")
	}
}

func TestNameWithNoHeaderAndNotDigitGender(t *testing.T) {
	s := `1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := NameOption{
		NoHeader: true,
		Gender:   "aaa",
	}

	if err := Name(r, w, o); err == nil {
		t.Error("Name with no header and not digit gender column symbol should raise error.")
	}
}

func TestNameWithNoHeaderAndNotDigitReference(t *testing.T) {
	s := `1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := NameOption{
		NoHeader:  true,
		Name:      "0",
		Reference: "aaa",
	}

	if err := Name(r, w, o); err == nil {
		t.Error("Name with no header and not digit reference column symbol should raise error.")
	}
}

func TestNameWithUnsupportedSpaceWidth(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := NameOption{
		Name:       "0",
		SpaceWidth: -1,
	}

	if err := Name(r, w, o); err == nil {
		t.Error("Name with negative space width should raise error.")
	}
}

func TestNameWithUnsupportedGenderFormat(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := NameOption{
		Gender:       "0",
		GenderFormat: "x",
	}

	if err := Name(r, w, o); err == nil {
		t.Error("Name with unknown gender format should raise error.")
	}
}

func TestNameOnNameNotFound(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := NameOption{
		Name: "ddd",
	}

	if err := Name(r, w, o); err == nil {
		t.Error("Name with unknown name column symbol should raise error.")
	}
}

func TestNameOnFirstNameNotFound(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := NameOption{
		FirstName: "ddd",
	}

	if err := Name(r, w, o); err == nil {
		t.Error("Name with unknown first name column symbol should raise error.")
	}
}

func TestNameOnLastNameNotFound(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := NameOption{
		LastName: "ddd",
	}

	if err := Name(r, w, o); err == nil {
		t.Error("Name with unknown last name column symbol should raise error.")
	}
}

func TestNameOnKanaNotFound(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := NameOption{
		Kana: "ddd",
	}

	if err := Name(r, w, o); err == nil {
		t.Error("Name with unknown kana column symbol should raise error.")
	}
}

func TestNameOnFirstKanaNotFound(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := NameOption{
		FirstKana: "ddd",
	}

	if err := Name(r, w, o); err == nil {
		t.Error("Name with unknown first kana column symbol should raise error.")
	}
}

func TestNameOnLastKanaNotFound(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := NameOption{
		LastKana: "ddd",
	}

	if err := Name(r, w, o); err == nil {
		t.Error("Name with unknown last kana column symbol should raise error.")
	}
}

func TestNameOnLastGenderNotFound(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := NameOption{
		Gender:       "ddd",
		GenderFormat: "jp_short",
	}

	if err := Name(r, w, o); err == nil {
		t.Error("Name with unknown gender column symbol should raise error.")
	}
}

func TestNameOnLastReferenceNotFound(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := NameOption{
		Name:      "aaa",
		Reference: "ddd",
	}

	if err := Name(r, w, o); err == nil {
		t.Error("Name with unknown reference column symbol should raise error.")
	}
}

func TestNameWithBrokenCSV(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := NameOption{
		Name: "aaa",
	}

	if err := Name(r, w, o); err == nil {
		t.Error("Name with broken csv should raise error.")
	}
}

func TestNameWithName(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := NameOption{
		Name: "aaa",
	}

	if err := Name(r, w, o); err != nil {
		t.Error(err)
	}

	actual := readCSV(w.String())
	for i, rec := range actual {
		if i == 0 {
			continue
		}
		s := rec[0]
		if isHiraganaOrSpace(s) {
			t.Errorf("name %s is hiragana", s)
		}
		if isKatakanaOrSpace(s) {
			t.Errorf("name %s is katakana", s)
		}
		if !isMultibyte(s) {
			t.Errorf("name %s is not multibyte", s)
		}
	}
}

func TestNameWithNameAndSpaceWidth1(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := NameOption{
		Name:       "aaa",
		SpaceWidth: 1,
	}

	if err := Name(r, w, o); err != nil {
		t.Error(err)
	}

	actual := readCSV(w.String())
	for i, rec := range actual {
		if i == 0 {
			continue
		}
		s := rec[0]
		if isHiraganaOrSpace(s) {
			t.Errorf("name %s is hiragana", s)
		}
		if isKatakanaOrSpace(s) {
			t.Errorf("name %s is katakana", s)
		}
		if !isMultibyte(s) {
			t.Errorf("name %s is not multibyte", s)
		}
		if !strings.Contains(s, " ") {
			t.Errorf("name %s does not have half space", s)
		}
	}
}

func TestNameWithNameAndSpaceWidth2(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := NameOption{
		Name:       "aaa",
		SpaceWidth: 2,
	}

	if err := Name(r, w, o); err != nil {
		t.Error(err)
	}

	actual := readCSV(w.String())
	for i, rec := range actual {
		if i == 0 {
			continue
		}
		s := rec[0]
		if isHiraganaOrSpace(s) {
			t.Errorf("name %s is hiragana", s)
		}
		if isKatakanaOrSpace(s) {
			t.Errorf("name %s is katakana", s)
		}
		if !isMultibyte(s) {
			t.Errorf("name %s is not multibyte", s)
		}
		if !strings.Contains(s, "　") {
			t.Errorf("name %s does not have full space", s)
		}
	}
}

func TestNameWithFirstName(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := NameOption{
		FirstName:  "aaa",
		SpaceWidth: 1,
	}

	if err := Name(r, w, o); err != nil {
		t.Error(err)
	}

	actual := readCSV(w.String())
	for i, rec := range actual {
		if i == 0 {
			continue
		}
		s := rec[0]
		if isKatakanaOrSpace(s) {
			t.Errorf("first name %s is katakana", s)
		}
		if !isMultibyte(s) {
			t.Errorf("first name %s is not multibyte", s)
		}
		if strings.Contains(s, " ") {
			t.Errorf("first name %s should not have space", s)
		}
	}
}

func TestNameWithLastName(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := NameOption{
		LastName:   "aaa",
		SpaceWidth: 1,
	}

	if err := Name(r, w, o); err != nil {
		t.Error(err)
	}

	actual := readCSV(w.String())
	for i, rec := range actual {
		if i == 0 {
			continue
		}
		s := rec[0]
		if isHiraganaOrSpace(s) {
			t.Errorf("last name %s is hiragana", s)
		}
		if isKatakanaOrSpace(s) {
			t.Errorf("last name %s is katakana", s)
		}
		if !isMultibyte(s) {
			t.Errorf("last name %s is not multibyte", s)
		}
		if strings.Contains(s, " ") {
			t.Errorf("last name %s should not have space", s)
		}
	}
}

func TestNameWithKanaAndSpaceWidth1(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := NameOption{
		Kana:       "aaa",
		SpaceWidth: 1,
	}

	if err := Name(r, w, o); err != nil {
		t.Error(err)
	}

	actual := readCSV(w.String())
	for i, rec := range actual {
		if i == 0 {
			continue
		}
		s := rec[0]
		if isHiraganaOrSpace(s) {
			t.Errorf("kana %s is hiragana", s)
		}
		if !isKatakanaOrSpace(s) {
			t.Errorf("kana %s should be katakana", s)
		}
		if !strings.Contains(s, " ") {
			t.Errorf("kana %s does not have half space", s)
		}
	}
}

func TestNameWithKanaAndSpaceWidth2(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := NameOption{
		Kana:       "aaa",
		SpaceWidth: 2,
	}

	if err := Name(r, w, o); err != nil {
		t.Error(err)
	}

	actual := readCSV(w.String())
	for i, rec := range actual {
		if i == 0 {
			continue
		}
		s := rec[0]
		if isHiraganaOrSpace(s) {
			t.Errorf("kana %s is hiragana", s)
		}
		if !isKatakanaOrSpace(s) {
			t.Errorf("kana %s should be katakana", s)
		}
		if !strings.Contains(s, "　") {
			t.Errorf("kana %s does not have half space", s)
		}
	}
}

func TestNameWithFirstKana(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := NameOption{
		FirstKana:  "aaa",
		SpaceWidth: 1,
	}

	if err := Name(r, w, o); err != nil {
		t.Error(err)
	}

	actual := readCSV(w.String())
	for i, rec := range actual {
		if i == 0 {
			continue
		}
		s := rec[0]
		if isHiraganaOrSpace(s) {
			t.Errorf("first kana %s is hiragana", s)
		}
		if !isKatakanaOrSpace(s) {
			t.Errorf("first kana %s should be katakana", s)
		}
		if strings.Contains(s, " ") {
			t.Errorf("first kana %s should not have space", s)
		}
	}
}

func TestNameWithLastKana(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := NameOption{
		LastKana:   "aaa",
		SpaceWidth: 1,
	}

	if err := Name(r, w, o); err != nil {
		t.Error(err)
	}

	actual := readCSV(w.String())
	for i, rec := range actual {
		if i == 0 {
			continue
		}
		s := rec[0]
		if isHiraganaOrSpace(s) {
			t.Errorf("last kana %s is hiragana", s)
		}
		if !isKatakanaOrSpace(s) {
			t.Errorf("last kana %s should be katakana", s)
		}
		if strings.Contains(s, " ") {
			t.Errorf("last kana %s should not have space", s)
		}
	}
}

func TestNameWithKanaAndHiraganaAndSpaceWidth1(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := NameOption{
		Kana:       "aaa",
		Hiragana:   true,
		SpaceWidth: 1,
	}

	if err := Name(r, w, o); err != nil {
		t.Error(err)
	}

	actual := readCSV(w.String())
	for i, rec := range actual {
		if i == 0 {
			continue
		}
		s := rec[0]
		if !isHiraganaOrSpace(s) {
			t.Errorf("kana %s should be hiragana", s)
		}
		if isKatakanaOrSpace(s) {
			t.Errorf("kana %s is katakana", s)
		}
		if !strings.Contains(s, " ") {
			t.Errorf("kana %s does not have half space", s)
		}
	}
}

func TestNameWithKanaAndHiraganaAndSpaceWidth2(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := NameOption{
		Kana:       "aaa",
		Hiragana:   true,
		SpaceWidth: 2,
	}

	if err := Name(r, w, o); err != nil {
		t.Error(err)
	}

	actual := readCSV(w.String())
	for i, rec := range actual {
		if i == 0 {
			continue
		}
		s := rec[0]
		if !isHiraganaOrSpace(s) {
			t.Errorf("kana %s should be hiragana", s)
		}
		if isKatakanaOrSpace(s) {
			t.Errorf("kana %s is katakana", s)
		}
		if !strings.Contains(s, "　") {
			t.Errorf("kana %s does not have half space", s)
		}
	}
}

func TestNameWithFirstKanaAndHiragana(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := NameOption{
		FirstKana:  "aaa",
		Hiragana:   true,
		SpaceWidth: 1,
	}

	if err := Name(r, w, o); err != nil {
		t.Error(err)
	}

	actual := readCSV(w.String())
	for i, rec := range actual {
		if i == 0 {
			continue
		}
		s := rec[0]
		if !isHiraganaOrSpace(s) {
			t.Errorf("first kana %s should be hiragana", s)
		}
		if isKatakanaOrSpace(s) {
			t.Errorf("first kana %s is katakana", s)
		}
		if strings.Contains(s, " ") {
			t.Errorf("first kana %s should not have space", s)
		}
	}
}

func TestNameWithLastKanaAndHiragana(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := NameOption{
		LastKana:   "aaa",
		Hiragana:   true,
		SpaceWidth: 1,
	}

	if err := Name(r, w, o); err != nil {
		t.Error(err)
	}

	actual := readCSV(w.String())
	for i, rec := range actual {
		if i == 0 {
			continue
		}
		s := rec[0]
		if !isHiraganaOrSpace(s) {
			t.Errorf("last kana %s should be hiragana", s)
		}
		if isKatakanaOrSpace(s) {
			t.Errorf("last kana %s is katakana", s)
		}
		if strings.Contains(s, " ") {
			t.Errorf("last kana %s should not have space", s)
		}
	}
}

func TestNameWithGender(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	for i, gf := range supportedGenderFormats {
		r := bytes.NewBufferString(s)
		w := &bytes.Buffer{}
		o := NameOption{
			Gender:       "aaa",
			GenderFormat: gf,
			SpaceWidth:   1,
		}

		if err := Name(r, w, o); err != nil {
			t.Error(err)
		}

		actual := readCSV(w.String())
		for j, rec := range actual {
			if j == 0 {
				continue
			}
			s := rec[0]
			genders := []string{maleGenders[i], femaleGenders[i]}
			if !containsString(genders, s) {
				t.Errorf("gender %s is invalid for %s format", s, gf)
			}
		}
	}
}

func TestNameWithGenderMaleRate100(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := NameOption{
		Gender:       "aaa",
		GenderFormat: "en_short",
		MaleRate:     100,
		SpaceWidth:   1,
	}

	if err := Name(r, w, o); err != nil {
		t.Error(err)
	}

	actual := readCSV(w.String())
	for j, rec := range actual {
		if j == 0 {
			continue
		}
		s := rec[0]
		if s == "F" {
			t.Errorf("gender %s should be male", s)
		}
	}
}

func TestNameWithGenderMaleRate0(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := NameOption{
		Gender:       "aaa",
		GenderFormat: "en_short",
		MaleRate:     0,
		SpaceWidth:   1,
	}

	if err := Name(r, w, o); err != nil {
		t.Error(err)
	}

	actual := readCSV(w.String())
	for j, rec := range actual {
		if j == 0 {
			continue
		}
		s := rec[0]
		if s == "M" {
			t.Errorf("gender %s should be female", s)
		}
	}
}

func TestNameWithReference(t *testing.T) {
	s := `aaa,bbb,ccc
1,タナカ アキラ,3
4,たなか あきら,6
7,田中 明,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := NameOption{
		Name:       "aaa",
		Reference:  "bbb",
		SpaceWidth: 1,
	}

	if err := Name(r, w, o); err != nil {
		t.Error(err)
	}

	actual := readCSV(w.String())
	for j, rec := range actual {
		if j == 0 {
			continue
		}
		s := rec[0]
		if !strings.HasPrefix(s, "田中") {
			t.Errorf("not named reference: %s", s)
		}
	}
}

func TestNameWithReferenceFullSpace(t *testing.T) {
	s := `aaa,bbb,ccc
1,タナカ　アキラ,3
4,たなか　あきら,6
7,田中　明,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := NameOption{
		Name:       "aaa",
		Reference:  "bbb",
		SpaceWidth: 1,
	}

	if err := Name(r, w, o); err != nil {
		t.Error(err)
	}

	actual := readCSV(w.String())
	for j, rec := range actual {
		if j == 0 {
			continue
		}
		s := rec[0]
		if !strings.HasPrefix(s, "田中") {
			t.Errorf("not named reference: %s", s)
		}
	}
}

func TestNameWithReferenceLastNameOnly(t *testing.T) {
	s := `aaa,bbb,ccc
1,タナカ,3
4,たなか,6
7,田中,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := NameOption{
		Name:       "aaa",
		Reference:  "bbb",
		SpaceWidth: 1,
	}

	if err := Name(r, w, o); err != nil {
		t.Error(err)
	}

	actual := readCSV(w.String())
	for j, rec := range actual {
		if j == 0 {
			continue
		}
		s := rec[0]
		if !strings.HasPrefix(s, "田中") {
			t.Errorf("not named reference: %s", s)
		}
	}
}

func TestNameWithReferenceAndMaleRate100(t *testing.T) {
	s := `aaa,bbb,ccc
1,タナカ アキラ,3
4,たなか あきら,6
7,田中 明,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := NameOption{
		Name:         "aaa",
		Reference:    "bbb",
		Gender:       "ccc",
		GenderFormat: "en_short",
		MaleRate:     100,
		SpaceWidth:   1,
	}

	if err := Name(r, w, o); err != nil {
		t.Error(err)
	}

	actual := readCSV(w.String())
	for j, rec := range actual {
		if j == 0 {
			continue
		}
		s := rec[0]
		if !strings.HasPrefix(s, "田中") {
			t.Errorf("not named reference: %s", s)
		}
		if rec[2] == "F" {
			t.Errorf("gender %s should be male", rec[2])
		}
	}
}

func TestNameWithReferenceAndMaleRate0(t *testing.T) {
	s := `aaa,bbb,ccc
1,タナカ アキラ,3
4,たなか あきら,6
7,田中 明,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := NameOption{
		Name:         "aaa",
		Reference:    "bbb",
		Gender:       "ccc",
		GenderFormat: "en_short",
		MaleRate:     0,
		SpaceWidth:   1,
	}

	if err := Name(r, w, o); err != nil {
		t.Error(err)
	}

	actual := readCSV(w.String())
	for j, rec := range actual {
		if j == 0 {
			continue
		}
		s := rec[0]
		if !strings.HasPrefix(s, "田中") {
			t.Errorf("not named reference: %s", s)
		}
		if rec[2] == "M" {
			t.Errorf("gender %s should be female", rec[2])
		}
	}
}

func TestNameWithReferenceOnReferenceFail(t *testing.T) {
	s := `aaa,bbb,ccc
1,Andrew,3
4,Justin,6
7,Peter,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := NameOption{
		Name:       "aaa",
		Reference:  "bbb",
		SpaceWidth: 1,
	}

	if err := Name(r, w, o); err != nil {
		t.Error(err)
	}

	actual := readCSV(w.String())
	for j, rec := range actual {
		if j == 0 {
			continue
		}
		s := rec[0]
		if s != "" {
			t.Errorf("empty should be output on reference failed: %s", s)
		}
	}
}

func TestNameWithRistrictReferenceOnReferenceFail(t *testing.T) {
	s := `aaa,bbb,ccc
1,Andrew,3
4,Justin,6
7,Peter,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := NameOption{
		Name:              "aaa",
		Reference:         "bbb",
		RistrictReference: true,
		SpaceWidth:        1,
	}

	if err := Name(r, w, o); err == nil {
		t.Error("Name with ristrict reference should raise error on reference failed.")
	}
}

func isHiraganaOrSpace(s string) bool {
	rs := []rune(s)
	for _, r := range rs {
		if r != ' ' && r != '　' && !unicode.In(r, unicode.Hiragana) {
			return false
		}
	}
	return true
}

func isKatakanaOrSpace(s string) bool {
	rs := []rune(s)
	for _, r := range rs {
		if r != ' ' && r != '　' && !unicode.In(r, unicode.Katakana) {
			return false
		}
	}
	return true
}

func isMultibyte(s string) bool {
	return len(s) != len([]rune(s))
}
