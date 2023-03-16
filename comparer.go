package main

import (
	"errors"
)

// 数组对比方法
func comparer(data1 [][]string, data2 [][]string) ([][]string, error) {
	// 行数不一致返回
	if len(data1) != len(data2) {
		return nil, errors.New("表格行数不一致")
	}

	var diffs [][]string
	// 循环对比
	for rowIndex, row := range data1 {
		for cellIndex, cell := range row {
			if cell != data2[rowIndex][cellIndex] {
				diffs = append(diffs, row, data2[rowIndex], []string{})
			}
		}
	}

	return diffs, nil
}
