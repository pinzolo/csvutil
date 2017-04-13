package csvutil

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestBuildingWithoutColumn(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := BuildingOption{
		NumberWidth: 1,
	}

	if err := Building(r, w, o); err == nil {
		t.Error("Building without column symbol should raise error.")
	}
}

func TestBuildingWithNoHeaderButColumnNotNumber(t *testing.T) {
	s := `1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := BuildingOption{
		NoHeader:    true,
		Column:      "foo",
		NumberWidth: 1,
	}

	if err := Building(r, w, o); err == nil {
		t.Error("Building with not number column symbol for no header CSV should raise error.")
	}

}

func TestBuildingWithNegativeOfficeRate(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := BuildingOption{
		Column:      "aaa",
		OfficeRate:  -1,
		NumberWidth: 1,
	}

	if err := Building(r, w, o); err == nil {
		t.Error("Building with negative office rate should raise error.")
	}
}

func TestBuildingWithOver100OfficeRate(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := BuildingOption{
		Column:      "aaa",
		OfficeRate:  101,
		NumberWidth: 1,
	}

	if err := Building(r, w, o); err == nil {
		t.Error("Building with over 100 office rate should raise error.")
	}
}

func TestBuildingWithUnknownColumn(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := BuildingOption{
		Column:      "ddd",
		NumberWidth: 1,
	}

	if err := Building(r, w, o); err == nil {
		t.Error("Building with unknown column should raise error.")
	}
}

func TestBuildingWithInvalidNumberWidth(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := BuildingOption{
		Column:      "aaa",
		NumberWidth: 0,
	}

	if err := Building(r, w, o); err == nil {
		t.Error("Building with invalid office rate should raise error.")
	}
}

func TestBuildingOnApartment(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := BuildingOption{
		Column:      "aaa",
		NumberWidth: 1,
	}

	if err := Building(r, w, o); err != nil {
		t.Error(err)
	}

	data := readCSV(w.String())
	if ok := allOK(data, 0, isHalfApartment); !ok {
		t.Errorf("Building failed updating on apartment. %+v", data)
	}
}

func TestBuildingWithNoHeader(t *testing.T) {
	s := `1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := BuildingOption{
		Column:      "0",
		NumberWidth: 1,
		NoHeader:    true,
	}

	if err := Building(r, w, o); err != nil {
		t.Error(err)
	}

	data := readCSV(w.String())
	if ok := allOKNoHeader(data, 0, isHalfApartment); !ok {
		t.Errorf("Building failed updating on apartment. %+v", data)
	}
}

func TestBuildingWithBrokenCSV(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := BuildingOption{
		Column:      "aaa",
		NumberWidth: 1,
	}

	if err := Building(r, w, o); err == nil {
		t.Error("Building with broken csv should raise error.")
	}
}

func TestBuildingOnApartmentFull(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := BuildingOption{
		Column:      "aaa",
		NumberWidth: 2,
	}

	if err := Building(r, w, o); err != nil {
		t.Error(err)
	}

	data := readCSV(w.String())
	if ok := allOK(data, 0, isFullApartment); !ok {
		t.Errorf("Building failed updating on apartment. %+v", data)
	}
}

func TestBuildingOnOffice(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := BuildingOption{
		Column:      "aaa",
		OfficeRate:  100,
		NumberWidth: 1,
	}

	if err := Building(r, w, o); err != nil {
		t.Error(err)
	}

	data := readCSV(w.String())
	if ok := allOK(data, 0, isHalfOffice); !ok {
		t.Errorf("Building failed updating on office. %+v", data)
	}
}

func TestBuildingOnOfficeFull(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := BuildingOption{
		Column:      "aaa",
		OfficeRate:  100,
		NumberWidth: 2,
	}

	if err := Building(r, w, o); err != nil {
		t.Error(err)
	}

	data := readCSV(w.String())
	if ok := allOK(data, 0, isFullOffice); !ok {
		t.Errorf("Building failed updating on office. %+v", data)
	}
}

func TestBuildingWithoutAppend(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := BuildingOption{
		Column:      "aaa",
		NumberWidth: 1,
	}

	if err := Building(r, w, o); err != nil {
		t.Error(err)
	}

	data := readCSV(w.String())
	if strings.HasPrefix(data[1][0], "1") {
		t.Errorf("Building should not append after source value. %+v", data)
	}
	if strings.HasPrefix(data[2][0], "4") {
		t.Errorf("Building should not append after source value. %+v", data)
	}
	if strings.HasPrefix(data[3][0], "7") {
		t.Errorf("Building should not append after source value. %+v", data)
	}
}

func TestBuildingWithAppend(t *testing.T) {
	s := `aaa,bbb,ccc
1,2,3
4,5,6
7,8,9
`
	r := bytes.NewBufferString(s)
	w := &bytes.Buffer{}
	o := BuildingOption{
		Column:      "aaa",
		NumberWidth: 1,
		Append:      true,
	}

	if err := Building(r, w, o); err != nil {
		t.Error(err)
	}

	data := readCSV(w.String())
	if !strings.HasPrefix(data[1][0], "1") {
		t.Errorf("Building should append after source value. %+v", data)
	}
	if !strings.HasPrefix(data[2][0], "4") {
		t.Errorf("Building should append after source value. %+v", data)
	}
	if !strings.HasPrefix(data[3][0], "7") {
		t.Errorf("Building should append after source value. %+v", data)
	}
}

func isApartment(s string) bool {
	for _, fn := range firstApartmentNames {
		if strings.Contains(s, fn) {
			return true
		}
	}
	for _, sn := range singleApartmentNames {
		if strings.Contains(s, sn) {
			return true
		}
	}
	return false
}

func isHalfApartment(s string) bool {
	if !isApartment(s) {
		return false
	}
	for _, n := range fullWidthNums {
		if strings.Contains(s, n) {
			fmt.Println(s, n)
			return false
		}
	}
	for _, n := range halfWidthNums {
		if strings.Contains(s, n) {
			return true
		}
	}
	return false
}

func isFullApartment(s string) bool {
	if !isApartment(s) {
		return false
	}
	for _, n := range halfWidthNums {
		if strings.Contains(s, n) {
			return false
		}
	}
	for _, n := range fullWidthNums {
		if strings.Contains(s, n) {
			return true
		}
	}
	return false
}

func isOffice(s string) bool {
	return strings.Contains(s, "ビル")
}

func isHalfOffice(s string) bool {
	if !isOffice(s) {
		return false
	}
	for _, n := range fullWidthNums {
		if strings.Contains(s, n) {
			return false
		}
	}
	for _, n := range halfWidthNums {
		if strings.Contains(s, n) {
			return true
		}
	}
	return false
}

func isFullOffice(s string) bool {
	if !isOffice(s) {
		return false
	}
	for _, n := range halfWidthNums {
		if strings.Contains(s, n) {
			return false
		}
	}
	for _, n := range fullWidthNums {
		if strings.Contains(s, n) {
			return true
		}
	}
	return false
}
