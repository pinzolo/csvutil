package csvutil

import (
	"bytes"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

func TestAddressWithoutTargetColumns(t *testing.T) {
	r := &bytes.Buffer{}
	w := &bytes.Buffer{}
	o := AddressOption{}
	if err := Address(r, w, o); err == nil {
		t.Error("Address without size should raise error.")
	}
}

func TestAddressWithUnsupportedNumberWidth(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := AddressOption{
		ZipCode:     "0",
		NumberWidth: -1,
	}
	if err := Address(r, w, o); err == nil {
		t.Error("Address with negative size should raise error.")
	}
}

func TestAddressWithNoHeaderAndNotDigitZipCode(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := AddressOption{
		NoHeader:    true,
		ZipCode:     "aaa",
		NumberWidth: 1,
	}

	err := Address(r, w, o)
	if err == nil {
		t.Error("Address with no header and not digit zip code column symbol should raise error.")
	}
}

func TestAddressWithNoHeaderAndNotDigitPrefecture(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := AddressOption{
		NoHeader:    true,
		Prefecture:  "aaa",
		NumberWidth: 1,
	}

	err := Address(r, w, o)
	if err == nil {
		t.Error("Address with no header and not digit prefecture column symbol should raise error.")
	}
}

func TestAddressWithNoHeaderAndNotDigitCity(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := AddressOption{
		NoHeader:    true,
		City:        "aaa",
		NumberWidth: 1,
	}

	err := Address(r, w, o)
	if err == nil {
		t.Error("Address with no header and not digit city column symbol should raise error.")
	}
}

func TestAddressWithNoHeaderAndNotDigitTown(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := AddressOption{
		NoHeader:    true,
		Town:        "aaa",
		NumberWidth: 1,
	}

	err := Address(r, w, o)
	if err == nil {
		t.Error("Address with no header and not digit town column symbol should raise error.")
	}
}

func TestAddressOnZipCodeNotFound(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := AddressOption{
		ZipCode:     "ddd",
		NumberWidth: 1,
	}

	err := Address(r, w, o)
	if err == nil {
		t.Error("Address with unknown zip code column symbol should raise error.")
	}
}

func TestAddressOnPrefectureNotFound(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := AddressOption{
		Prefecture:  "ddd",
		NumberWidth: 1,
	}

	err := Address(r, w, o)
	if err == nil {
		t.Error("Address with unknown prefecture column symbol should raise error.")
	}
}

func TestAddressOnCityNotFound(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := AddressOption{
		City:        "ddd",
		NumberWidth: 1,
	}

	err := Address(r, w, o)
	if err == nil {
		t.Error("Address with unknown city column symbol should raise error.")
	}
}

func TestAddressOnTownNotFound(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := AddressOption{
		Town:        "ddd",
		NumberWidth: 1,
	}

	err := Address(r, w, o)
	if err == nil {
		t.Error("Address with unknown town column symbol should raise error.")
	}
}

func TestAddressWithZipCode(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := AddressOption{
		ZipCode:     "aaa",
		NumberWidth: 1,
	}

	if err := Address(r, w, o); err != nil {
		t.Error(err)
	}

	rgx := regexp.MustCompile(`\d{3}-\d{4}`)
	actual := readCSV(w.String())
	for i, rec := range actual {
		if i == 0 {
			continue
		}
		if !rgx.MatchString(rec[0]) {
			t.Errorf("Zip code not found: %s, line: %d", rec[0], i)
		}
	}
}

func TestAddressWithPrefecture(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := AddressOption{
		Prefecture:  "aaa",
		NumberWidth: 1,
	}

	if err := Address(r, w, o); err != nil {
		t.Error(err)
	}

	actual := readCSV(w.String())
	for i, rec := range actual {
		if i == 0 {
			continue
		}
		if !containsString(prefs, rec[0]) {
			t.Errorf("Prefecture not found: %s, line: %d", rec[0], i)
		}
	}
}

func TestAddressWithPrefectureCode(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := AddressOption{
		Prefecture:     "aaa",
		PrefectureCode: true,
		NumberWidth:    1,
	}

	if err := Address(r, w, o); err != nil {
		t.Error(err)
	}

	actual := readCSV(w.String())
	for i, rec := range actual {
		if i == 0 {
			continue
		}
		c, err := strconv.Atoi(rec[0])
		if err != nil {
			t.Errorf("Prefecture should output as code: %s, line: %d", rec[0], i)
		}
		if c <= 0 && 48 <= c {
			t.Errorf("Invalid prefecture code: %s, line: %d", rec[0], i)
		}
	}
}

func TestAddressWithCity(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := AddressOption{
		City:        "aaa",
		NumberWidth: 1,
	}

	if err := Address(r, w, o); err != nil {
		t.Error(err)
	}

	actual := readCSV(w.String())
	for i, rec := range actual {
		if i == 0 {
			continue
		}
		if _, err := strconv.Atoi(rec[0]); err == nil {
			t.Errorf("City not found: %s, line: %d", rec[0], i)
		}
	}
}

func TestAddressWithTown(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := AddressOption{
		Town:        "aaa",
		NumberWidth: 1,
	}

	if err := Address(r, w, o); err != nil {
		t.Error(err)
	}

	actual := readCSV(w.String())
	for i, rec := range actual {
		if i == 0 {
			continue
		}
		if _, err := strconv.Atoi(rec[0]); err == nil {
			t.Errorf("Town not found: %s, line: %d", rec[0], i)
		}
		for _, n := range halfWidthNums {
			if strings.Contains(rec[0], n) {
				t.Errorf("Town should not have block number: %s, line: %d", rec[0], i)
			}
		}
	}
}

func TestAddressWithBlockNumber(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := AddressOption{
		Town:        "aaa",
		BlockNumber: true,
		NumberWidth: 1,
	}

	if err := Address(r, w, o); err != nil {
		t.Error(err)
	}

	actual := readCSV(w.String())
	for i, rec := range actual {
		if i == 0 {
			continue
		}
		if _, err := strconv.Atoi(rec[0]); err == nil {
			t.Errorf("Town not found: %s, line: %d", rec[0], i)
		}

		ok := false
		for _, n := range halfWidthNums {
			if strings.Contains(rec[0], n) {
				ok = true
			}
		}
		if !ok {
			t.Errorf("Town should have half width block number: %s", rec[0])
		}
	}
}

func TestAddressWithFullWidthBlockNumber(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := AddressOption{
		Town:        "aaa",
		BlockNumber: true,
		NumberWidth: 2,
	}

	if err := Address(r, w, o); err != nil {
		t.Error(err)
	}

	actual := readCSV(w.String())
	for i, rec := range actual {
		if i == 0 {
			continue
		}
		if _, err := strconv.Atoi(rec[0]); err == nil {
			t.Errorf("Town not found: %s, line: %d", rec[0], i)
		}

		ok := false
		for _, n := range fullWidthNums {
			if strings.Contains(rec[0], n) {
				ok = true
			}
		}
		if !ok {
			t.Errorf("Town should have full width block number: %s", rec[0])
		}
	}
}

func TestAddressWithAppending(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := AddressOption{
		ZipCode:     "0",
		Prefecture:  "aaa",
		City:        "0",
		Town:        "aaa",
		BlockNumber: true,
		NumberWidth: 2,
	}

	if err := Address(r, w, o); err != nil {
		t.Error(err)
	}

	actual := readCSV(w.String())
	for i, rec := range actual {
		if i == 0 {
			continue
		}
		if _, err := strconv.Atoi(rec[0]); err == nil {
			t.Errorf("Invalid town: %s, line: %d", rec[0], i)
		}
		if !regexp.MustCompile(`^\d{3}-\d{4}`).MatchString(rec[0]) {
			t.Errorf("Zip code not found: %s", rec[0])
		}

		prefOK := false
		for _, p := range prefs {
			if strings.Contains(rec[0], p) {
				prefOK = true
			}
		}
		if !prefOK {
			t.Errorf("Prefecture not found: %s", rec[0])
		}

		bnOK := false
		for _, n := range fullWidthNums {
			if strings.Contains(rec[0], n) {
				bnOK = true
			}
		}
		if !bnOK {
			t.Errorf("Town should have block number: %s", rec[0])
		}
	}

}
