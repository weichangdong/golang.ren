package main

import (
	"bufio"
	"fmt"
	"io"
	"sort"
	"strings"
)

const WhichFuncWillUse = "desc" //desc

// 比较两个 `desc table`的结果的不同
// desc一个表,然后吧字段组合成可insert
var diffTest = `

`

var diffOnline = `

`

var descStr = `
| id                            | bigint(20) unsigned | NO   | PRI | 0       |       |
| order_id                      | bigint(20) unsigned | NO   | MUL | 0       |       |
| amount                        | bigint(20) unsigned | NO   |     | 0       |       |
| amount_payed                  | bigint(20) unsigned | NO   |     | 0       |       |
`

func main() {
	if WhichFuncWillUse == "diff" {
		diff()
	} else if WhichFuncWillUse == "desc" {
		desc()
	}
}

func desc() {
	str := parseStr(descStr)
	fmt.Println("`" + strings.Join(str, "`,`") + "`")
}

func diff() {
	afterDiff1 := parseStr(diffTest)
	afterDiff2 := parseStr(diffOnline)
	fmt.Println("Len1:", len(afterDiff1), " Len2:", len(afterDiff2))
	theDiff := findDiff(afterDiff1, afterDiff2)
	fmt.Println(theDiff)
}

func InArray(f string, all []string) bool {
	for _, o := range all {
		if f == o {
			return true
		}
	}
	return false
}

func findDiff(d1, d2 []string) []string {
	re := []string{}
	for _, d := range d1 {
		if InArray(d, d2) == false {
			re = append(re, d)
		}
	}

	for _, d := range d2 {
		if InArray(d, d1) == false {
			re = append(re, d)
		}
	}

	return re
}
func parseStr(diffStr string) []string {
	sr := strings.NewReader(diffStr)
	buf := bufio.NewReader(sr)

	allFields := []string{}
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			} else {
				fmt.Println("Read file error!", err)
				return allFields
			}
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		one := strings.Split(line, "|")
		if len(one) < 2 {
			continue
		}
		field := strings.TrimSpace(one[1])
		if field == "" {
			continue
		}
		allFields = append(allFields, field)
	}
	if WhichFuncWillUse == "diff" {
		sort.Strings(allFields)
	}
	return allFields
}
