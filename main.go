package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
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
			// Transport: &debug.DebugRequestTransport{
			// 	RequestHeader:  true,
			// 	RequestBody:    true,
			// 	ResponseHeader: true,
			// 	ResponseBody:   false,
			// },
		},
	})
}

type SelfListener struct {
}

func (l *SelfListener) ProgressChangedCallback(event *cos.ProgressEvent) {
	switch event.EventType {
	case cos.ProgressStartedEvent:
	case cos.ProgressDataEvent:
		percent := event.ConsumedBytes * 100 / event.TotalBytes
		fmt.Printf("%3d%%\033[4D", percent)
	case cos.ProgressCompletedEvent:
		fmt.Print("[\033[32m ✓ \033[0m]")
		fmt.Printf("%12.2f\n", B2MB(event.TotalBytes))
	}
}

func B2MB(val int64) float64 {
	return float64(val) / (1 << 10)
}

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

func GetFile(remote string, local string, flag bool) {
	if local[len(local)-1] == '/' {
		local = filepath.Join(local, filepath.Base(remote))
	} // local is dir
	if !flag && DirExist(local) {
		var ch string
		fmt.Printf("\033[32m> %s 已存在，是否覆盖(y or n)?\033[0m ", local)
		fmt.Scanf("%s", &ch)
		if ch != "y" {
			return
		}
		fmt.Print("\033[1A")
	}
	if !IsDir(local) {
		mkpath := filepath.Dir(local)
		err := os.MkdirAll(mkpath, os.ModePerm)
		if err != nil {
			fmt.Printf("\033[31m> %s \033[0m\n", err)
			return
		}
		fmt.Printf("\033[31m> %s is nor found\033[0m\n", local)
		fmt.Printf("\033[32m> %s is Created\033[0m\n", mkpath)
		fmt.Print("\033[2A")
	}

	fmt.Printf("%-8s%-30s%s%-30s%-3s", "--", remote, "=> ", local, " ")

	opt := &cos.ObjectGetOptions{
		Listener: &SelfListener{},
	}
	_, err := c.Object.GetToFile(context.Background(), remote, local, opt)
	if err != nil {
		fmt.Println("[\033[31m ✕ \033[0m]")
		fmt.Println("\033[32mError...", err, "\033[0m")
		return
	}

}

func GetDir(remote string, local string) {
	var flag bool
	if IsDir(local) {
		var ch string
		fmt.Printf("\033[32m> %s 已存在，是否全部覆盖(a / n / e)?\033[0m ", local)
		fmt.Scanf("%s", &ch)
		if ch == "n" || ch == "a" {
			flag = (ch == "a")
			fmt.Print("\033[1A")
		} else {
			return
		}
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
		fmt.Println("\033[32mError...", err, "\033[0m")
		return
	}
	isTruncated := true
	for isTruncated {
		opt.Marker = marker
		v, _, err := c.Bucket.Get(context.Background(), opt)
		if err != nil {
			fmt.Println(err)
			break
		}
		for _, content := range v.Contents {
			gremote := content.Key
			glocal := filepath.Join(local, content.Key)
			GetFile(gremote, glocal, flag)
		}
		isTruncated = v.IsTruncated // 是否还有数据
		marker = v.NextMarker       // 设置下次请求的起始 key
	}
}

func ListDir(remote string) {
	if remote == "." {
		remote = ""
	}
	var marker string
	opt := &cos.BucketGetOptions{
		Prefix:    remote, // prefix表示要查询的文件夹
		Delimiter: "",     // deliter表示分隔符, 设置为/表示列出当前目录下的object, 设置为空表示列出所有的object
		MaxKeys:   100,    // 设置最大遍历出多少个对象, 一次listobject最大支持1000
	}
	_, _, err := c.Bucket.Get(context.Background(), opt)
	if err != nil {
		fmt.Println("[\033[31m ✕ \033[0m]")
		fmt.Println("\033[31mError...", err, "\033[0m")
		return
	}
	isTruncated := true
	for isTruncated {
		opt.Marker = marker
		v, _, err := c.Bucket.Get(context.Background(), opt)
		if err != nil {
			fmt.Println(err)
			break
		}
		for _, content := range v.Contents {
			fmt.Printf("\033[0m%-8s%-30s%s%-30s%-8s%12.3f\033[0m\n",
				"--", content.Key, "-- ", "-", "\033[3C[\033[32m ✓ \033[0m]", B2MB(content.Size))
		}
		// common prefix表示表示被delimiter截断的路径, 如delimter设置为/, common prefix则表示所有子目录的路径
		for _, commonPrefix := range v.CommonPrefixes {
			fmt.Printf("\033[0m%-8s%-30s%s%-30s%-8s%12s\033[0m\n",
				"Dir =>", commonPrefix, "-- ", "-", "\033[3C[\033[32m ✓ \033[0m]", "---")
		}
		isTruncated = v.IsTruncated // 是否还有数据
		marker = v.NextMarker       // 设置下次请求的起始 key
	}
}

func ArgError() {
	fmt.Println("Usage: cosput [-f file] [-d dir] [path]")
	fmt.Println("Default: path = .")
}

func main() {
	cnt := len(os.Args)
	if cnt == 3 || cnt == 4 {
		remote := os.Args[2]
		local := "./"

		if cnt == 4 {
			local = os.Args[3]
		}
		fmt.Printf("\033[1;7m%-8s%-30s%s%-30s%8s%12s\033[0m\n",
			"Type", "Remote", "=> ", "Local", "State", "Size(KB)")
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

	ArgError()
}
