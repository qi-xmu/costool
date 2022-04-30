package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/tencentyun/cos-go-sdk-v5"
)

var c *cos.Client

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
}

func main() {
	cnt := len(os.Args) // 参数长度

	if cnt == 3 || cnt == 4 {
		remote := os.Args[2]
		local := "./"
		if cnt == 4 {
			local = os.Args[3]
		}
		fmt.Printf("\033[1;7m%-8s%-30s%s%-30s%8s%12s\033[0m\n",
			"No.", "Remote", "=> ", "Local", "State", "Size(KB)")
		switch os.Args[1][1] {
		case 'f': // file
			GetFile(remote, local, false)
		case 'd': // dir
			GetDir(remote, local)
		case 'l':
			ListDir(remote)
		default:
			ArgError()
		}
		return
	}
	if cnt == 2 {
		switch os.Args[1][1] {
		case 'b':
			SelectBucket()
		default:
			ArgError()
		}
		return
	}
	ArgError()
}
