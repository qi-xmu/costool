package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/tencentyun/cos-go-sdk-v5"
)

var c *cos.Client

var put = flag.Bool("p", false, "上传文件. -p [local]path [remote]key")
var get = flag.Bool("g", false, "下载文件. -g [remote]Prekey [local]Dir")
var buc = flag.Bool("b", false, "更改bucket. -b")
var lis = flag.Bool("l", false, "列出所有文件 -l [remote]key*")
var ove = flag.Bool("y", false, "默认覆盖文件")

func init() {
	init_config()
	conf_url := fmt.Sprintf("https://%s.cos.%s.myqcloud.com", conf.Bucket, conf.Region)
	u, _ := url.Parse(conf_url)
	b := &cos.BaseURL{BucketURL: u}
	c = cos.NewClient(b, &http.Client{
		Timeout: 100 * time.Second,
		Transport: &cos.AuthorizationTransport{
			SecretID:     conf.SecretID,
			SecretKey:    conf.SecretKey,
			SessionToken: conf.SessionToken,
		},
	})
	fmt.Print("\033[1A\033[K\033[1A\033[K")
}

func main() {
	flag.Parse()
	if *buc {
		SelectBucket()
	}
	if *put {
		PutFunc()
	}
	if *get {
		GetFunc()
	}
	if *lis {
		ListFunc()
	}
	if *buc && !*put && !*lis && !*get {
		ArgError()
	}
}
func ArgError() {
	flag.PrintDefaults()
}
