package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/tencentyun/cos-go-sdk-v5"
)

// main get func
func GetFunc() {
	narg := flag.NArg()
	if narg != 1 && narg != 2 {
		ArgError()
		return
	}
	remote := flag.Args()[0]
	local := "./"
	if flag.NArg() == 2 {
		local = flag.Args()[1]
	}
	fmt.Printf("\033[1;7m%-8s%-30s%s%-30s%8s%12s\033[0m\n",
		"Type", "Remote", "=> ", "Local", "State ", "Size(KB)")
	GetKey(remote, local)
	os.Exit(0)
}

type SelfListener struct {
}

func (l *SelfListener) ProgressChangedCallback(event *cos.ProgressEvent) {
	switch event.EventType {
	case cos.ProgressStartedEvent:
	case cos.ProgressFailedEvent:
	case cos.ProgressDataEvent:
		percent := event.ConsumedBytes * 100 / event.TotalBytes
		fmt.Printf("%3d%%\033[4D", percent)
	case cos.ProgressCompletedEvent:

	}
}

func GetFile(remote string, local string, flag bool) {
	// 如果文件夹存在 或者结尾为/ 识别为 文件夹
	if DirString(local) || IsDir(local) {
		// 如果文件夹不存在
		fmt.Println(local)
		if !DirExist(local) {
			MkDir(local)
		}
		local = filepath.Join(local, filepath.Base(remote))
	} else if !DirExist(filepath.Dir(local)) {
		MkDir(local)
	}
	// 其他识别为 文件 flag 为覆盖之类
	if !flag && FileExist(local) {
		var ch string
		fmt.Printf("\033[32m> %s 已存在，是否覆盖(y or n)?\033[0m ", local)
		fmt.Scanf("%s", &ch)
		fmt.Print("\033[1A")
		if ch != "y" {
			GetPrint(remote, local)
			fmt.Print("[\033[31m ✕ \033[0m]\n")
			return
		}
	}
	GetPrint(remote, local)
	opt := &cos.ObjectGetOptions{
		Listener: &SelfListener{},
	}
	Resp, err := c.Object.GetToFile(context.Background(), remote, local, opt)
	if err != nil {
		fmt.Print("[\033[31m ✕ \033[0m]")
		ErrorP(err)
		return
	}
	fmt.Print("[\033[32m ✓ \033[0m]")
	fmt.Printf("%12.3f\n", B2KB(Resp.ContentLength))
}

func GetKey(remote string, local string) {
	var flag bool = *ove
	if IsDir(local) && !flag {
		var ch string
		fmt.Printf("\033[32m> %s 已存在，是否全部覆盖(a / n / e)?\033[0m ", local)
		fmt.Scanf("%s", &ch)
		if ch != "n" && ch != "a" {
			return
		}
		flag = (ch == "a")
		fmt.Print("\033[1A")
	}
	if remote == "." {
		remote = ""
	}
	var marker string
	opt := &cos.BucketGetOptions{
		Prefix:    remote, // prefix表示要查询的文件夹
		Delimiter: "",     // deliter表示分隔符, 设置为/表示列出当前目录下的object, 设置为空表示列出所有的object
		MaxKeys:   100,    // 设置最大遍历出多少个对象, 一次listobject最大支持1000
	}
	// 获取目录
	_, _, err := c.Bucket.Get(context.Background(), opt)
	if err != nil {
		fmt.Println("[\033[31m ✕ \033[0m]")
		ErrorP(err)
		return
	}
	isTruncated := true
	for isTruncated {
		opt.Marker = marker
		v, _, err := c.Bucket.Get(context.Background(), opt)
		if err != nil {
			ErrorP(err)
			break
		}
		for _, content := range v.Contents {
			gremote := content.Key
			rremote := strings.Replace(gremote, remote, "", -1) // 相对的路径
			glocal := filepath.Join(local, rremote)             // dir --> dir
			if DirString(local) {
				glocal += gremote[strings.LastIndexByte(gremote, '/'):]
			} // file --> dir
			GetFile(gremote, glocal, flag)
		}
		isTruncated = v.IsTruncated // 是否还有数据
		marker = v.NextMarker       // 设置下次请求的起始 key
	}
}
