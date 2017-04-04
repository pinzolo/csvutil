package csvutil

import (
	"strconv"

	"github.com/pkg/errors"
)

type column struct {
	symbol string
	index  int
	found  bool
}

func (c *column) findIndex(hdr []string) error {
	if isDigit(c.symbol) {
		i, err := strconv.Atoi(c.symbol)
		if err != nil {
			return errors.Wrap(err, "Cannot parse index")
		}
		c.index = i
		c.found = true
		return nil
	}
	if hdr == nil {
		return errors.New("Column symbol require index only in no header.")
	}
	for i, h := range hdr {
		if h == c.symbol {
			c.index = i
			c.found = true
			return nil
		}
	}
	return errors.Errorf("Column %s is not found", c.symbol)
}

func newColumn(sym string) *column {
	return &column{
		symbol: sym,
		index:  -1,
	}
}

func newColumns(syms []string) []*column {
	cols := make([]*column, len(syms))
	for i, sym := range syms {
		cols[i] = newColumn(sym)
	}
	return cols
}

func newColumnsWithIndexes(syms []string, hdr []string) ([]*column, error) {
	cols := newColumns(syms)
	for _, col := range cols {
		err := col.findIndex(hdr)
		if err != nil {
			return nil, errors.Wrap(err, "Cannot find index")
		}
	}
	return cols, nil
}

func newUniqueColumns(syms []string, hdr []string) ([]*column, error) {
	cols, err := newColumnsWithIndexes(syms, hdr)
	if err != nil {
		return nil, err
	}
	return uniqColumns(cols), nil
}
