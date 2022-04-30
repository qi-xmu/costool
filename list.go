package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/tencentyun/cos-go-sdk-v5"
)

// main list func
func ListFunc() {
	narg := flag.NArg()
	if narg != 1 && narg != 0 {
		ArgError()
		return
	}
	remote := ""
	if narg == 1 {
		remote = flag.Args()[0]
		fmt.Printf("Remote %s\n", remote)
	}
	fmt.Printf("\033[1;7m%-36s%-36s%8s%12s\033[0m\n",
		"Prefix", "Name", "State ", "Size(KB)")
	ListDir(remote)
}

func ListDir(remote string) {
	var marker string
	opt := &cos.BucketGetOptions{
		Prefix:    remote, // prefix表示要查询的文件夹
		Delimiter: "",     // deliter表示分隔符, 设置为/表示列出当前目录下的object, 设置为空表示列出所有的object
		MaxKeys:   100,    // 设置最大遍历出多少个对象, 一次listobject最大支持1000
	}
	_, _, err := c.Bucket.Get(context.Background(), opt)
	if err != nil {
		fmt.Println("[\033[31m ✕ \033[0m]")
		ErrorP(err)
		os.Exit(1)
	}
	isTruncated := true
	for isTruncated {
		opt.Marker = marker
		v, _, err := c.Bucket.Get(context.Background(), opt)
		if err != nil {
			ErrorP(err)
			os.Exit(1)
		}
		for _, content := range v.Contents {
			p := strings.LastIndexByte(content.Key, '/') + 1
			Prefix := content.Key[:p] // 前缀
			name := content.Key[p:]   // name
			fmt.Printf("\033[0m%-36s%-35s%-8s%12.3f\033[0m\n",
				Prefix, name, "\033[3C[\033[32m ✓ \033[0m]", B2KB(content.Size))
		}
		isTruncated = v.IsTruncated // 是否还有数据
		marker = v.NextMarker       // 设置下次请求的起始 key
	}
}
