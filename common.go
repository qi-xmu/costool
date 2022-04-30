package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func B2KB(val int64) float64 {
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

func DirString(path string) bool {
	return path[len(path)-1] == '/'
}

func DirExist(path string) bool {
	_, err := os.Open(path)
	return err == nil
}

func FileExist(path string) bool {
	return DirExist(path)
}

func FileSize(filename string) int64 {
	obj, err := os.Open(filename)
	if err != nil {
		ErrorP("File or Dir not Found")
		os.Exit(1)
	}
	fileInfo, _ := obj.Stat()
	return fileInfo.Size()
}

func MkDir(local string) {
	mkpath := filepath.Dir(local)
	err := os.MkdirAll(mkpath, os.ModePerm)
	if err != nil {
		fmt.Printf("\033[31m> %s \033[0m\n", err)
		return
	}
	fmt.Printf("\033[31m> %s is not found. \033[0m\n", local)
	fmt.Printf("\033[32m> %s is Created. \033[0m\n", mkpath)
	fmt.Print("\033[2A")
}
func ErrorP(a ...interface{}) {
	fmt.Println("\033[31m!", a, "\033[0m")
}

func slice(val string, n int) string {
	l := len(val)
	res := val
	if l > n {
		res = val[:n-3] + ".."
	}
	return res
}

func GetPrint(remote, local string) {
	fmt.Printf("%-8s%-30s%s%-30s%-2s", "--", slice(remote, 30), "=> ", slice(local, 30), "")
}
func PutPrint(local, remote string) {
	fmt.Printf("%-8s%-30s%s%-30s%-2s", "FIlE", slice(local, 30), "=> ", slice(remote, 30), "")
}
