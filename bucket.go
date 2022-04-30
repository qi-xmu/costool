package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/tencentyun/cos-go-sdk-v5"
)

func bucket() []cos.Bucket {
	c = cos.NewClient(nil, &http.Client{
		Timeout: 100 * time.Second,
		Transport: &cos.AuthorizationTransport{
			SecretID:     conf.SecretID,
			SecretKey:    conf.SecretKey,
			SessionToken: conf.SessionToken,
		},
	})
	Resp, _, err := c.Service.Get(context.Background())
	if err != nil {
		panic(err)
	}
	// fmt.Printf("\033[7m> %-3s%-20s%-20s\033[0m\n", "No", "Bucket Name", "Bucket Region")
	// for i, b := range Resp.Buckets {
	// 	fmt.Printf("  %-3d%-20s%-20s\n", i, b.Name, b.Region)
	// }
	return Resp.Buckets
}

func SelectBucket() {

	buck := bucket()
	blen := len(buck)
	for {
		var tmp string
		fmt.Printf("\033[7m> %-3s%-20s%-20s\033[0m\n", "No", "Bucket Name", "Bucket Region")
		for i, b := range buck {
			fmt.Printf("  %-3d%-20s%-20s\033[0m\n", i, b.Name, b.Region)
		}
		fmt.Printf("\033[32m> 选择你使用的存储桶(0, 1, ... or e): \033[0m")
		fmt.Scan(&tmp)

		if tmp == "e" {
			os.Exit(0)
		}
		num, err := strconv.Atoi(tmp)
		if num < blen && err == nil {
			fmt.Printf("%25s\n", "")
			fmt.Printf("\033[32m> 选中存储桶 %d, %s\033[0m\n", num, buck[num].Name)
			conf.Bucket = buck[num].Name
			conf.Region = buck[num].Region

			SaveConfig()
			return
		}
	}
}

// func test() {
// 	buck := bucket()
// 	blen := len(buck)
// 	for

// }
