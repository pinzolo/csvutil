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

func TestAddressWithUnsupportedBlockNumberWidth(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := AddressOption{
		ZipCode:          "0",
		BlockNumberWidth: -1,
	}
	if err := Address(r, w, o); err == nil {
		t.Error("Address with negative size should raise error.")
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
		ZipCode:          "aaa",
		BlockNumberWidth: 1,
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
		Prefecture:       "aaa",
		BlockNumberWidth: 1,
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
		Prefecture:       "aaa",
		PrefectureCode:   true,
		BlockNumberWidth: 1,
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
		City:             "aaa",
		BlockNumberWidth: 1,
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
		Town:             "aaa",
		BlockNumberWidth: 1,
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
		Town:             "aaa",
		BlockNumber:      true,
		BlockNumberWidth: 1,
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
		Town:             "aaa",
		BlockNumber:      true,
		BlockNumberWidth: 2,
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
		ZipCode:          "0",
		Prefecture:       "aaa",
		City:             "0",
		Town:             "aaa",
		BlockNumber:      true,
		BlockNumberWidth: 2,
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
