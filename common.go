package main

import (
	"fmt"
	"os"
)

func B2MB(val int64) float64 {
	return float64(val) / (1 << 10)
}

// 判断是否是文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

func DirExist(filename string) bool {
	_, err := os.Open(filename)
	return err == nil
}

func FileExist(filename string) int64 {
	obj, err := os.Open(filename)
	if err != nil {
		ErrorP("File or Dir not Found")
		os.Exit(1)
	}
	fileInfo, _ := obj.Stat()
	return fileInfo.Size()
}

func ErrorP(a ...interface{}) {
	fmt.Println("\033[31m!", a, "\033[0m")
}
