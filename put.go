package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

func PutFile(source string, target string) {
	size := float64(FileSize(source)) / (1 << 20)
	if size > (1 << 5) {
		fmt.Println("\033[31mFile too big, Max size = 32MB.\033[0m")
		return
	}
	source_path, _ := filepath.Abs(source)
	fmt.Printf("%-8s%-12.2f%-30s%s%-30s%-3s", "FIlE", size, source, "=> ", target, " ")

	if target[len(target)-1] == '/' {
		target = filepath.Join(target, filepath.Base(source_path))
	}
	target = strings.Replace(target, "\\", "/", -1)
	// _, err := c.Object.PutFromFile(context.Background(), target, source_path, nil)
	_, _, err := c.Object.Upload(context.Background(), target, source_path, nil)

	if err != nil {
		fmt.Println("[\033[31m ✕ \033[0m]")
		fmt.Println("\033[32mError...", err, "\033[0m")
		return
	}
	fmt.Println("[\033[32m ✓ \033[0m]")

}

func PutDir(source string, target string) {
	source_path := source
	fmt.Printf("\033[1;32m%-8s%-30s\033[0m\n", "DIR =>", source_path)

	files, _ := ioutil.ReadDir(source_path)
	for _, fi := range files {
		if fi.IsDir() {
			if fi.Name()[0] == '.' {
				continue
			}
			tsource := filepath.Join(source, fi.Name())
			ttarget := filepath.Join(target, fi.Name())
			PutDir(tsource, ttarget)
		} else {
			fsource := filepath.Join(source, fi.Name())
			ftarget := filepath.Join(target, fi.Name())
			PutFile(fsource, ftarget)
		}
	}
}
