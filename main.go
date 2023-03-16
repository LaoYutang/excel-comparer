package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/xuri/excelize/v2"
)

var dirPath = "./"

func main() {
	files, err := os.ReadDir(dirPath + "/old")
	if err != nil {
		fmt.Println("文件夹读取失败！")
	}

	os.RemoveAll(dirPath + "/diffs")
	os.MkdirAll(dirPath+"/diffs", os.ModePerm)

	for _, file := range files {
		fmt.Printf(">>>>>> 开始对比文件 %v \n", file.Name())
		compareTwoFiles(file.Name())
	}

	fmt.Println("按回车键关闭")
	fmt.Scanln()
}

func compareTwoFiles(fileName string) {
	path1 := "./old/" + fileName
	path2 := "./new/" + fileName
	file1, err1 := excelize.OpenFile(path1)
	if err1 != nil {
		fmt.Printf("文件打开失败 %v \n", path1)
	}
	file2, err2 := excelize.OpenFile(path2)
	if err2 != nil {
		fmt.Printf("文件打开失败 %v \n", path2)
	}

	// 遍历所有工作表
	for _, sheetName := range file1.GetSheetMap() {
		// 获取两张工作表数据
		data1, sheetErr1 := file1.GetRows(sheetName)
		if sheetErr1 != nil {
			panic(sheetErr1)
		}
		data2, sheetErr2 := file2.GetRows(sheetName)
		if sheetErr2 != nil {
			panic(sheetErr2)
		}

		// 对比数据并输出结果
		diffs, comparerErr := comparer(data1, data2)
		var diffFile *excelize.File
		if (comparerErr != nil) || (len(diffs) > 0) {
			if diffFile == nil {
				diffFile = excelize.NewFile()
			}
			diffFile.NewSheet(sheetName)
			defer diffFile.SaveAs(dirPath + "/diffs/" + fileName)
		}

		if comparerErr != nil {
			diffFile.SetCellValue(sheetName, "A1", "行数不一致")
			fmt.Printf("工作表 %v 行数不相同 \n", sheetName)
		} else if len(diffs) > 0 {
			fmt.Printf("工作表 %v 存在差异, 差异行如下: \n", sheetName)
			for index, diff := range diffs {
				fmt.Println(diff)
				diffFile.SetSheetRow(sheetName, "A"+strconv.Itoa(index+1), diff)
			}
		} else {
			fmt.Printf("工作表 %v 无差异 \n", sheetName)
		}
	}
}
