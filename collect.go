package csvutil

import (
	"io"
	"sort"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

var supportedSortKeys = []string{"value", "count"}

type collectedItem struct {
	value string
	count int
}

// CollectOption is option holder for Collect.
type CollectOption struct {
	// Source file does not have header line. (default false)
	NoHeader bool
	// Encoding of source file. (default utf8)
	Encoding string
	// Encoding for output.
	OutputEncoding string
	// Column symbol of target column
	Column string
	// AllowEmpty as value
	AllowEmpty bool
	// Print count
	PrintCount bool
	// Sort flag
	Sort bool
	// Sort key (value or count)
	SortKey string
	// Sort in descending order
	Descending bool
}

func (o CollectOption) validate() error {
	if o.Column == "" {
		return errors.New("no column")
	}
	if o.NoHeader {
		if !isEmptyOrDigit(o.Column) {
			return errors.New("not number column symbol")
		}
	}
	if o.SortKey != "" && !containsString(supportedSortKeys, o.SortKey) {
		return errors.Errorf("unsupported sort key: %s", o.SortKey)
	}
	return nil
}

func (o CollectOption) outputEncoding() string {
	if o.OutputEncoding != "" {
		return o.OutputEncoding
	}
	return o.Encoding
}

// Collect print collected CSV values.
func Collect(r io.Reader, w io.Writer, o CollectOption) error {
	if err := o.validate(); err != nil {
		return errors.Wrap(err, "invalid option")
	}

	cr, bom := reader(r, o.Encoding)

	var col *column
	csvp := NewReadOnlyCSVProcessor(cr)
	if o.NoHeader {
		csvp.SetPreBodyRead(func() error {
			col = newColumnWithIndex(o.Column, nil)
			return col.err
		})
	} else {
		csvp.SetHeaderHanlder(func(hdr []string) ([]string, error) {
			col = newColumnWithIndex(o.Column, hdr)
			return hdr, col.err
		})
	}
	var items []*collectedItem
	csvp.SetRecordHandler(func(rec []string) ([]string, error) {
		s := rec[col.index]
		if !o.AllowEmpty && s == "" {
			return nil, nil
		}
		for _, item := range items {
			if item.value == s {
				item.count++
				return nil, nil
			}
		}
		items = append(items, &collectedItem{value: s, count: 1})
		return nil, nil
	})
	if err := csvp.Process(); err != nil {
		return err
	}

	if o.Sort {
		items = sortCollectedItems(items, o)
	}

	cw := writer(w, bom, o.outputEncoding())
	cw.Comma = '\t'
	defer cw.Flush()
	for _, item := range items {
		if o.PrintCount {
			cw.Write([]string{strconv.Itoa(item.count), item.value})
			continue
		}
		cw.Write([]string{item.value})
	}

	return nil
}

func sortCollectedItems(items []*collectedItem, o CollectOption) []*collectedItem {
	if o.SortKey == "count" {
		sort.Slice(items, func(i, j int) bool {
			if o.Descending {
				return items[i].count > items[j].count
			}
			return items[i].count < items[j].count
		})
	} else {
		sort.Slice(items, func(i, j int) bool {
			if o.Descending {
				return strings.Compare(items[i].value, items[j].value) > 0

			}
			return strings.Compare(items[i].value, items[j].value) < 0
		})
	}
	return items
}
