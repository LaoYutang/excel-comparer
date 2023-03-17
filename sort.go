package main

import (
	"sort"
	"strings"
)

func sortRows(rows [][]string) {
	sort.Slice(rows, func(i int, j int) bool {
		str1 := strings.Join(rows[i], ",")
		str2 := strings.Join(rows[j], ",")
		return str1 < str2
	})
}
