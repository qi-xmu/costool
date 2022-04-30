package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/tencentyun/cos-go-sdk-v5"
)

// main put func

func PutFunc() {
	narg := flag.NArg()
	if narg != 1 && narg != 2 {
		ArgError()
		return
	}
	local := flag.Args()[0]
	remote := flag.Args()[0]
	if flag.NArg() == 2 {
		remote = flag.Args()[1]
	}
	fmt.Printf("\033[1;7m%-8s%-30s%s%-30s%8s%12s\033[0m\n",
		"Type", "Local", "=> ", "Remote", "State ", "Size(KB)")
	if IsDir(local) {
		remote = remote + "/" + filepath.Base(local)
		PutDir(local, remote)
	} else {
		PutFile(local, remote)
	}
}

func PutFile(local string, remote string) {
	size := float64(FileSize(local)) / (1 << 20)
	if size > (1 << 5) {
		fmt.Println("\033[31mFile too big, Max size = 32MB.\033[0m")
		return
	}
	remote = strings.Replace(remote, "\\", "/", -1)
	if DirString(remote) {
		remote = filepath.Join(remote, filepath.Base(local))
	}

	opt := &cos.ObjectPutOptions{
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			Listener: &SelfListener{},
		},
	}
	PutPrint(local, remote)
	_, err := c.Object.PutFromFile(context.Background(), remote, local, opt)

	if err != nil {
		fmt.Print("[\033[31m ✕ \033[0m]")
		ErrorP(err)
		return
	}
}

func PutDir(local string, remote string) {
	fmt.Printf("\033[1;4;32m%-8s%-70s\033[0m\n", "> DIR", local)

	files, _ := ioutil.ReadDir(local)
	for _, fi := range files {
		source := filepath.Join(local, fi.Name())
		target := filepath.Join(remote, fi.Name())
		if fi.IsDir() {
			if fi.Name()[0] == '.' {
				continue
			}
			PutDir(source, target) // 递归所有文件
		} else {
			PutFile(source, target) // 出口
		}
	}
	fmt.Printf("\033[12C\033[1;4;32m%-8s%70s\033[0m\n", "< END", local)
}
