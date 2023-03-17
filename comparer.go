package main

// 数组对比方法
func comparer(data1 [][]string, data2 [][]string, resChan chan [][]string) {
	var diffs [][]string

	// 循环对比
	for rowIndex, row := range data1 {
		for cellIndex, cell := range row {
			if cell != data2[rowIndex][cellIndex] {
				diffs = append(diffs, row, data2[rowIndex], []string{})
			}
		}
	}

	resChan <- diffs
}
