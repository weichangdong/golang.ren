package main

import (
	"flag"
	"fmt"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
)

var opType string

const OpTypeAsubB = "a-b"
const OpTypeBsubA = "b-a"
const OpTypeAaddB = "a+b"
const OpTypeAintersectionB = "ab"

// a文件名称 a文件用来操作的列
var aFileName string
var aFileColName string

var bFileName string
var bFileColName string

var outputFileName string

var excelSheetName string
var debug bool

const Desc = "用a文件的某n列 和 b文件的某n列去做比较"

func init() {
	flag.BoolVar(&debug, "debug", false, "是否debug debug输出一些中间日志")
	flag.StringVar(&outputFileName, "output-file", "", "输出类型文件名称 为空只输出到屏幕")
	flag.StringVar(&excelSheetName, "sheet", "Sheet1", "excel文件工作表名称")
	flag.StringVar(&opType, "type", "", "操作类型 [a-b] a有的b没有 [b-a] b有的a没有 [a+b] a和b的并集(去重) [ab] a和b的都有的,即交集")
	flag.StringVar(&aFileName, "a", "", "a excel文件")
	flag.StringVar(&bFileName, "b", "", "b excel文件")
	flag.StringVar(&aFileColName, "acol", "", "a excel文件需要比较的列的头名称,多个以逗号隔开")
	flag.StringVar(&bFileColName, "bcol", "", "b excel文件需要比较的列的头名称,多个以逗号隔开")
	aFileName = strings.TrimSpace(aFileName)
	bFileName = strings.TrimSpace(bFileName)
	aFileColName = strings.TrimSpace(aFileColName)
	bFileColName = strings.TrimSpace(bFileColName)
	flag.Parse()
}

func checkPara() bool {
	if !InStringSlice(opType, []string{OpTypeAsubB, OpTypeBsubA, OpTypeAaddB, OpTypeAintersectionB}) {
		fmt.Println("type 操作类型错误 opType:", opType)
		return false
	}
	if aFileName == "" || bFileName == "" {
		fmt.Println("a或者b 文件地址为空")
		return false
	}
	if aFileColName == "" || bFileColName == "" {
		fmt.Println("a或者b 要比较的列为空")
		return false
	}
	return true
}

func main() {
	check := checkPara()
	if !check {
		flag.Usage()
		return
	}

	fmt.Println("aFileName:", aFileName, " aFileColName:", aFileColName)
	fmt.Println("bFileName:", bFileName, " bFileColName:", bFileColName)
	err, contentsA, headerA := ReadExcelFile(aFileName, aFileColName)
	if err != nil || len(contentsA) == 0 {
		fmt.Println("读取excel错误 fileName:", aFileName)
		return
	}
	err, contentsB, headerB := ReadExcelFile(bFileName, bFileColName)
	if err != nil || len(contentsB) == 0 {
		fmt.Println("读取excel错误 fileName:", bFileName)
		return
	}
	data := map[string][]string{}
	switch opType {
	case OpTypeAsubB:
		data = AsubB(contentsA, contentsB)
	case OpTypeBsubA:
		data = AsubB(contentsB, contentsA)
	case OpTypeAaddB:
		data = AaddB(contentsA, contentsB)
	case OpTypeAintersectionB:
		data = AintersectionB(contentsA, contentsB)
	}
	if debug {
		fmt.Println("data:", data)
	}
	outputData(opType, aFileName, bFileName, data, outputFileName, headerA, headerB)
}

func outputData(opType string, aFileName string, bFileName string, data map[string][]string, outputFileName string, headerA, headerB []string) {
	outputDesc := ""
	header := ""
	switch opType {
	case OpTypeAsubB:
		outputDesc = aFileName + " [有] " + bFileName + " [没有] 数据如下:"
		header = strings.Join(headerA, "\t")
	case OpTypeBsubA:
		outputDesc = bFileName + " [有] " + aFileName + " [没有] 数据如下:"
		header = strings.Join(headerB, "\t")
	case OpTypeAaddB:
		outputDesc = aFileName + "  " + bFileName + " [并集] 数据如下:"
		header = strings.Join(headerB, "\t")
	case OpTypeAintersectionB:
		outputDesc = aFileName + "  " + bFileName + " [交集] 数据如下:"
		header = strings.Join(headerB, "\t")
	}
	fmt.Println(outputDesc)
	fmt.Println(header)

	showData := [][]string{}
	for _, one := range data {
		/*
			timeStr := one[6]
			timeFloat, _ := Str2Float64(timeStr)
			oktime := TimeFromExcelTime(timeFloat, true)
			timeDisplay := outputTimeByDate(oktime)
			fmt.Println(timeDisplay)
		*/
		fmt.Println(strings.Join(one, "\t"))
		showData = append(showData, one)
	}
	if outputFileName != "" {
		WriteExcelFile(headerA, showData, outputFileName)
	}

}

func outputTimeByUnix(timestamp int64) string {
	local, _ := time.LoadLocation("Asia/Shanghai")
	tm := time.Unix(timestamp, 0)
	return tm.In(local).Format("2006-01-02 15:04:05")
}

func outputTimeByDate(dt time.Time) string {
	return Int2Str(dt.Hour()) + ":" + Int2Str(dt.Minute()) + ":" + Int2Str(dt.Second())
}

// a有的b没有
func AsubB(a, b map[string][]string) (sub map[string][]string) {
	sub = make(map[string][]string, 0)
	for key, oneA := range a {
		if _, ok := b[key]; !ok {
			sub[key] = oneA
		}
	}
	return
}

// a和b的并集
func AaddB(a, b map[string][]string) (add map[string][]string) {
	add = make(map[string][]string, 0)
	for key, oneA := range a {
		add[key] = oneA
	}
	for key, oneB := range b {
		add[key] = oneB
	}
	return
}

// a和b的都有的,即交集
func AintersectionB(a, b map[string][]string) (intersection map[string][]string) {
	intersection = make(map[string][]string, 0)
	for key, oneA := range a {
		if _, ok := b[key]; ok {
			intersection[key] = oneA
		}
	}
	for key, oneB := range b {
		if _, ok := a[key]; ok {
			intersection[key] = oneB
		}
	}
	return
}

func ReadExcelFile(fileName, fileColName string) (err error, allContents map[string][]string, excelHeader []string) {
	var file *excelize.File
	file, err = excelize.OpenFile(fileName)
	if err != nil {
		fmt.Println("ReadExcelFile error err:", err)
		return
	}

	fileColNames := strings.Split(fileColName, ",")
	allContents = make(map[string][]string, 0)
	rowsNum := 0
	fileColNum := []int{}
	fileColExists := false
	rows := file.GetRows(excelSheetName)
	for k2, row := range rows {
		if k2 == 0 {
			for k1, colCell := range row {
				excelHeader = append(excelHeader, colCell)
				if InStringSlice(colCell, fileColNames) {
					fileColNum = append(fileColNum, k1)
					fileColExists = true
				}
			}
		}
		if !fileColExists {
			fmt.Println("excel文件需要比较的列的头名称 不存在这个 fileColName:", fileColName)
			break
		}
		if k2 == 0 {
			continue
		}
		fileColContent := []string{}
		for k1, colCell := range row {
			if InIntSlice(k1, fileColNum) {
				fileColContent = append(fileColContent, colCell)
			}
		}
		fileColContentStr := strings.Join(fileColContent, "-")
		allContents[fileColContentStr] = row
		rowsNum++
	}
	if debug {
		fmt.Println("rowsNum:", rowsNum, " excelHeader:", excelHeader, " aFileColContent:", allContents)
	}
	return
}

func WriteExcelFile(header []string, data [][]string, fileName string) {
	xlsx := excelize.NewFile()
	sheetName := "Sheet1"
	for k, name := range header {
		xlsx.SetCellValue(sheetName, LetterMap(k)+"1", name)
	}
	for k, one := range data {
		kn := k + 2
		for i, o := range one {
			xlsx.SetCellValue(sheetName, LetterMap(i)+Int2Str(kn), o)
		}
	}
	xlsx.SaveAs(fileName)
}

func FileName(a, b string) string {
	a1 := filepath.Base(a)
	b1 := filepath.Base(b)
	ext := path.Ext(a)
	return a1 + b1 + "." + ext
}

func Int2Str(number int) string {
	return strconv.FormatInt(int64(number), 10)
}

func Str2Float64(s string) (f float64, err error) {
	f, err = strconv.ParseFloat(s, 64)

	return
}

func InStringSlice(v string, sl []string) bool {
	for _, vv := range sl {
		if strings.ToLower(strings.TrimSpace(vv)) == strings.ToLower(strings.TrimSpace(v)) {
			return true
		}
	}
	return false
}

func InIntSlice(v int, sl []int) bool {
	for _, vv := range sl {
		if vv == v {
			return true
		}
	}
	return false
}

func LetterMap(i int) string {
	base := 65
	ok := fmt.Sprintf("%s", string(base+i))
	if i > 52 {
		ok = "B" + ok
	} else if i > 26 {
		ok = "A" + ok
	}
	return ok
}
