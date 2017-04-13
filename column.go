package csvutil

import (
	"strconv"

	"github.com/pkg/errors"
)

type column struct {
	symbol string
	index  int
	err    error
}

type columns []*column

func (cs columns) err() error {
	for _, c := range cs {
		if c.err != nil {
			return c.err
		}
	}
	return nil
}

func (c *column) findIndex(hdr []string) error {
	if c.symbol == "" {
		return nil
	}
	if isDigit(c.symbol) {
		i, _ := strconv.Atoi(c.symbol)
		c.index = i
		return nil
	}
	if hdr == nil {
		return errors.New("not number column symbol")
	}
	for i, h := range hdr {
		if h == c.symbol {
			c.index = i
			return nil
		}
	}
	return errors.Errorf("column %s not found", c.symbol)
}

func newColumnWithIndex(sym string, hdr []string) *column {
	col := &column{
		symbol: sym,
		index:  -1,
	}
	err := col.findIndex(hdr)
	if err != nil {
		col.err = err
	}
	return col
}

func newColumnsWithIndexes(syms []string, hdr []string) columns {
	cols := make([]*column, len(syms))
	for i, sym := range syms {
		cols[i] = newColumnWithIndex(sym, hdr)
	}
	return cols
}

func newUniqueColumns(syms []string, hdr []string) columns {
	cols := newColumnsWithIndexes(syms, hdr)
	return uniqColumns(cols)
}

func uniqColumns(cols columns) columns {
	var newCols []*column
	for _, col := range cols {
		exists := false
		for _, newCol := range newCols {
			if col.index != -1 && newCol.index == col.index {
				exists = true
				break
			}
		}
		if !exists {
			newCols = append(newCols, col)
		}
	}
	return newCols
}
