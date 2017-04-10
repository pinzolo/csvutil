package csvutil

import (
	"strconv"

	"github.com/pkg/errors"
)

type column struct {
	symbol string
	index  int
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

func newColumnWithIndex(sym string, hdr []string) (*column, error) {
	col := newColumn(sym)
	err := col.findIndex(hdr)
	if err != nil {
		return nil, err
	}
	return col, nil
}

func newColumnsWithIndexes(syms []string, hdr []string) ([]*column, error) {
	cols := newColumns(syms)
	for _, col := range cols {
		err := col.findIndex(hdr)
		if err != nil {
			return nil, err
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

func uniqColumns(cols []*column) []*column {
	var newCols []*column
	for _, col := range cols {
		exists := false
		for _, newCol := range newCols {
			if newCol.index == col.index {
				exists = true
				break
			}
		}
		if col.index != -1 && !exists {
			newCols = append(newCols, col)
		}
	}
	return newCols
}
