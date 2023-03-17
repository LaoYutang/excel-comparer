package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"sync"

	"github.com/xuri/excelize/v2"
)

var dirPath = "./"
var fileGroup sync.WaitGroup

func main() {
	// 读取old文件夹中的文件
	files, err := os.ReadDir(dirPath + "/old")
	if err != nil {
		panic("文件夹 ./old 读取失败！")
	}

	// 重置diffs文件夹
	os.RemoveAll(dirPath + "/diffs")
	os.MkdirAll(dirPath+"/diffs", os.ModePerm)

	for _, file := range files {
		fmt.Printf(">>>>>> 开始对比文件 %v \n", file.Name())
		fileGroup.Add(1)
		go compareTwoFiles(file.Name())
	}
	fileGroup.Wait()

	fmt.Println(">>>>>> 对比完成，按回车键关闭")
	fmt.Scanln()
}

// 文件级别对比
func compareTwoFiles(fileName string) {
	// 打开两个文件，打开失败直接返回
	path1 := "./old/" + fileName
	path2 := "./new/" + fileName
	file1, err1 := excelize.OpenFile(path1)
	if err1 != nil {
		fmt.Printf("[error]文件打开失败 %v \n", path1)
		return
	}
	file2, err2 := excelize.OpenFile(path2)
	if err2 != nil {
		fmt.Printf("[error]文件打开失败 %v \n", path2)
		return
	}

	var sheetGroup sync.WaitGroup
	// 遍历所有工作表
	for _, sheetName := range file1.GetSheetMap() {
		sheetGroup.Add(1)
		go compareTwoSheet(file1, file2, fileName, sheetName, &sheetGroup)
	}
	sheetGroup.Wait()

	fileGroup.Done()
}

// 工作表级比较
func compareTwoSheet(file1 *excelize.File, file2 *excelize.File, fileName string, sheetName string, sheetGroup *sync.WaitGroup) {
	// 获取两张工作表数据
	data1Chan := make(chan [][]string)
	data2Chan := make(chan [][]string)
	go func() {
		data1, sheetErr1 := file1.GetRows(sheetName)
		if sheetErr1 != nil {
			fmt.Printf("[error]工作表获取失败 %v old %v", fileName, sheetName)
			return
		}
		data1Chan <- data1
	}()
	go func() {
		data2, sheetErr2 := file2.GetRows(sheetName)
		if sheetErr2 != nil {
			fmt.Printf("[error]工作表获取失败 %v old %v", fileName, sheetName)
			return
		}
		data2Chan <- data2
	}()
	data1 := <-data1Chan
	data2 := <-data2Chan
	close(data1Chan)
	close(data2Chan)

	// 差异切片
	var diffs [][]string

	if len(data1) != len(data2) {
		// 如果行数不一致
		diffs = append(diffs, []string{"表格行数不一致"})
	} else {
		const maxLength float64 = 2000              // 每个协程最大行数
		length := float64(len(data1))               // 总行数
		count := int(math.Ceil(length / maxLength)) // 循环总数
		resChan := make(chan [][]string, 20)        // 结果输出管道
		resNum := count                             // 获取结果计数

		// 开启count个协程处理获取差异列
		for i := 0; i < count; i++ {
			var subData1, subData2 [][]string
			if i == count-1 {
				// 最后一个
				subData1 = data1[i*int(maxLength):]
				subData2 = data2[i*int(maxLength):]
			} else {
				subData1 = data1[i*int(maxLength) : (i+1)*int(maxLength)]
				subData2 = data2[i*int(maxLength) : (i+1)*int(maxLength)]
			}

			go comparer(subData1, subData2, resChan)
		}

		// 获取协程处理数据
		for {
			diffs = append(diffs, <-resChan...)
			if resNum--; resNum == 0 {
				close(resChan)
				break
			}
		}
	}

	// 打印到stdout并生成deffs文件
	var diffFile *excelize.File
	if len(diffs) > 0 {
		// 如果存在差异
		if diffFile == nil {
			diffFile = excelize.NewFile()
		}
		diffFile.NewSheet(sheetName)
		defer diffFile.SaveAs(dirPath + "/diffs/" + fileName)
		fmt.Printf("工作表 %v --> %v 存在差异\n", fileName, sheetName)
		for index, diff := range diffs {
			fmt.Println(diff)
			diffFile.SetSheetRow(sheetName, "A"+strconv.Itoa(index+1), &diff)
		}
	} else {
		// 无差异
		fmt.Printf("工作表 %v --> %v 无差异 \n", fileName, sheetName)
	}

	sheetGroup.Done()
}
