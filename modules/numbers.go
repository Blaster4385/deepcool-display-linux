package modules

import (
	"errors"
)

type Pattern [][]bool

var DigitPatterns = map[int]Pattern{
	0: {
		{true, true, true},
		{true, false, true},
		{true, false, true},
		{true, false, true},
		{true, true, true},
	},
	1: {
		{false, true, false},
		{true, true, false},
		{false, true, false},
		{false, true, false},
		{true, true, true},
	},
	2: {
		{true, true, true},
		{false, false, true},
		{true, true, true},
		{true, false, false},
		{true, true, true},
	},
	3: {
		{true, true, true},
		{false, false, true},
		{true, true, true},
		{false, false, true},
		{true, true, true},
	},
	4: {
		{true, false, true},
		{true, false, true},
		{true, true, true},
		{false, false, true},
		{false, false, true},
	},
	5: {
		{true, true, true},
		{true, false, false},
		{true, true, true},
		{false, false, true},
		{true, true, true},
	},
	6: {
		{true, true, true},
		{true, false, false},
		{true, true, true},
		{true, false, true},
		{true, true, true},
	},
	7: {
		{true, true, true},
		{false, false, true},
		{false, true, false},
		{false, true, false},
		{false, true, false},
	},
	8: {
		{true, true, true},
		{true, false, true},
		{true, true, true},
		{true, false, true},
		{true, true, true},
	},
	9: {
		{true, true, true},
		{true, false, true},
		{true, true, true},
		{false, false, true},
		{true, true, true},
	},
}

var SymbolPatterns = map[string]Pattern{
	"celsius": {
		{true, false, false, false, false},
		{false, false, true, true, false},
		{false, true, false, false, false},
		{false, true, false, false, false},
		{false, false, true, true, false},
	},
	"fahrenheit": {
		{true, false, true, true, false},
		{false, false, true, false, false},
		{false, false, true, true, false},
		{false, false, true, false, false},
		{false, false, true, false, false},
	},
	"percent": {
		{false, false, false, false, false},
		{false, true, false, false, true},
		{false, false, false, true, false},
		{false, false, true, false, false},
		{false, true, false, false, true},
	},
}

func InsertPattern(grid [][]bool, pattern Pattern, row, col int) {
	for i, rowPattern := range pattern {
		for j, val := range rowPattern {
			if row+i < len(grid) && col+j < len(grid[0]) {
				grid[row+i][col+j] = val
			}
		}
	}
}

func CreateNumberGrid(value int, symbol string, row int) ([][]bool, error) {
	if value < 0 || value >= 1000 {
		return nil, errors.New("value must be between 0 and 999")
	}
	if _, ok := SymbolPatterns[symbol]; !ok {
		return nil, errors.New("unsupported symbol")
	}

	grid := make([][]bool, 14)
	for i := range grid {
		grid[i] = make([]bool, 14)
	}

	var (
		digits    []int
		symbolCol int
	)

	if value < 100 {
		digits = []int{value / 10, value % 10}
		symbolCol = 9
	} else {
		digits = []int{value / 100, (value % 100) / 10, value % 10}
		symbolCol = 13
	}

	col := 1
	for _, digit := range digits {
		InsertPattern(grid, DigitPatterns[digit], row, col)
		col += 4
	}
	InsertPattern(grid, SymbolPatterns[symbol], row, symbolCol)

	return grid, nil
}
