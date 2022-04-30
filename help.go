package main

import "fmt"

func ArgError() {
	fmt.Println("> Usage: cosget [-f remote_file] [-d remote_dir] [local_path]")
	fmt.Println("> Hint: 当 remote_dir = . 时 匹配所有文件")
}
