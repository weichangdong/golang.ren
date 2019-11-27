package file

import (
	"bufio"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

// 一次性读取整个文件
func ReadAll(path string) {
	fi, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	fd, _ := ioutil.ReadAll(fi)
	return string(fd)
}

// 从字符串一次一行读取
func ReadLineFromStr(str string) []string {
	sr := strings.NewReader(str)
	buf := bufio.NewReader(sr)

	allFields := []string{}
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return allFields
			}
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		allFields = append(allFields, line)
	}
	return allFields
}

// 从文件一次一行读取
func ReadLineFromStr(path string) []string {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	buf := bufio.NewReader(f)

	allFields := []string{}
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return allFields
			}
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		allFields = append(allFields, line)
	}
	return allFields
}
