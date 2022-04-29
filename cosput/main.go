package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/tencentyun/cos-go-sdk-v5"
)

var c *cos.Client

func init() {
	//将<bucket>和<region>修改为真实的信息
	//bucket的命名规则为{name}-{appid} ，此处填写的存储桶名称必须为此格式
	u, _ := url.Parse("https://qi-1300340355.cos.ap-chengdu.myqcloud.com")
	b := &cos.BaseURL{BucketURL: u}
	c = cos.NewClient(b, &http.Client{
		//设置超时时间
		Timeout: 100 * time.Second,
		Transport: &cos.AuthorizationTransport{
			//如实填写账号和密钥，也可以设置为环境变量
			SecretID:  "AKIDhIiwzrwbnpfLp8rzXu66vY0gDYKuen5E",
			SecretKey: "A9w7B9wYyGRS1PcgYIvMbeWDvWwwkMn3",
		},
	})
}

func ArgError() {
	fmt.Println("Usage: cosput [-f file] [-d dir] [path]")
	fmt.Println("Default: path = .")
}

func FileSize(filename string) int64 {
	obj, err := os.Open(filename)
	if err != nil {
		panic("File not Found")
	}
	fileInfo, _ := obj.Stat()
	return fileInfo.Size()
}

func PutFile(source string, target string) {
	size := float64(FileSize(source)) / (1 << 20)
	if size > (1 << 20) {
		fmt.Println("\033[31mFile too big, Max size = 200MB.\033[0m")
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

func main() {
	cnt := len(os.Args)
	if cnt == 3 || cnt == 4 {
		source := os.Args[2]

		target := "./"
		if cnt == 4 {
			target = os.Args[3]
		}
		fmt.Printf("\033[1;7m%-8s%-12s%-30s%s%-30s%8s\033[0m\n",
			"Type", "Size(MB)", "Source", "=> ", "Target", "State")
		switch os.Args[1][1] {
		case 'f': // file
			target = filepath.Base(source)
			PutFile(source, target)
		case 'd': // dir
			PutDir(source, target)
		default:
			ArgError()
		}
		return
	}
	ArgError()
}
