package main

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
)

/*
用于把一张表的一条数据,输出insert的sql
*/

var tableName = "repay_plan"
var sqlBuild map[string]string = make(map[string]string, 0)
var isIdUseNull = false
var sqlBuildBetter = []map[string]string{}
var (
	okSql       string
	okSqlBetter string
)

var sqlStr = `
                           id: 150762293
                     order_id: 19112043
                       amount: 4945
                 amount_payed: 0
`

func main() {
	sr := strings.NewReader(sqlStr)
	buf := bufio.NewReader(sr)

	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			} else {
				fmt.Println("Read file error!", err)
				return
			}
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		firstPos := strings.Index(line, ":")
		field := strings.TrimSpace(line[0:firstPos])
		value := strings.TrimSpace(line[firstPos+1:])
		//fmt.Println(field, value)
		sqlBuild[field] = value
		sqlBuildBetter = append(sqlBuildBetter, map[string]string{
			field: value,
		})
	}
	//makeSql(sqlBuild)
	makeSqlBetter(sqlBuildBetter)
}
func makeSql(sqlBuild map[string]string) {
	fieldOk, valueOk := "", ""
	totalNum := len(sqlBuild)
	okNum := 0
	for k, v := range sqlBuild {
		if totalNum-1 == okNum {
			fieldOk += "`" + k + "`"
		} else {
			fieldOk += "`" + k + "`,"
		}

		if k == "id" && isIdUseNull {
			if okNum == 0 {
				valueOk += "NULL"
			} else {
				valueOk += "," + "NULL"
			}
		} else {
			isAllNum, _ := regexp.Match("^[0-9]+$", []byte(v))
			if isAllNum {
				if okNum == 0 {
					valueOk += v
				} else {
					valueOk += "," + v
				}
			} else {
				if okNum == 0 {
					valueOk += "'" + v + "'"
				} else {
					valueOk += ",'" + v + "'"
				}
			}

		}
		okNum++
	}
	okSql += "insert into " + tableName + "(" + fieldOk + ") values (" + valueOk + ");"
	fmt.Println("\n")
	fmt.Println(okSql)
	fmt.Println("\n")
}

func makeSqlBetter(sqlBuildBetter []map[string]string) {
	fieldOk, valueOk := "", ""
	totalNum := len(sqlBuildBetter)
	okNum := 0
	for _, tmpMap := range sqlBuildBetter {
		for k, v := range tmpMap {
			if totalNum-1 == okNum {
				fieldOk += "`" + k + "`"
			} else {
				fieldOk += "`" + k + "`,"
			}

			if k == "id" && isIdUseNull {
				if okNum == 0 {
					valueOk += "NULL"
				} else {
					valueOk += "," + "NULL"
				}
			} else {
				isAllNum, _ := regexp.Match("^[0-9]+$", []byte(v))
				if isAllNum {
					if okNum == 0 {
						valueOk += v
					} else {
						valueOk += "," + v
					}
				} else {
					if okNum == 0 {
						valueOk += "'" + v + "'"
					} else {
						valueOk += ",'" + v + "'"
					}
				}

			}
			okNum++
		}
	}
	okSqlBetter += "insert into " + tableName + "(" + fieldOk + ") values (" + valueOk + ");"
	fmt.Println("\n")
	fmt.Println(okSqlBetter)
	fmt.Println("\n")
}
