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

func newColumns(sym string) []*column {
	syms := split(sym)
	cols := make([]*column, len(syms))
	for i, sym := range syms {
		cols[i] = newColumn(sym)
	}
	return cols
}
