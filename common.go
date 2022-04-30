package main

import "os"

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

func FileSize(filename string) int64 {
	obj, err := os.Open(filename)
	if err != nil {
		panic("File not Found")
	}
	fileInfo, _ := obj.Stat()
	return fileInfo.Size()
}
