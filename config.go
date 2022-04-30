package main

import (
	"fmt"
	"io/ioutil"
	"os/user"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Config struct {
	UserID string `yaml:"UserID"`
	Region string `yaml:"Region"`
	Bucket string `yaml:"Bucket"`

	SecretID     string `yaml:"SecretID"`
	SecretKey    string `yaml:"SecretKey"`
	SessionToken string `yaml:"Token"`
}

var conf Config
var default_config_path = ".costool.yaml"

func init_config() {
	// 读取默认路径
	user, err := user.Current()
	if nil == err {
		default_config_path = filepath.Join(user.HomeDir, default_config_path)
		fmt.Println("\033[32m> load config in ", default_config_path, "\033[0m")
	} else {
		panic("default_config_path")
	}

	config(default_config_path) // 加载config
	// 检测是否设定bucket
	if conf.Bucket == "" {
		SelectBucket()
	}
	fmt.Println("\033[32m> 选中", conf.Bucket, "\033[0m")
}

// 加载config
func config(path string) {
	bytes, _ := ioutil.ReadFile(path)
	err := yaml.Unmarshal(bytes, &conf)
	if err != nil {
		panic("config parser fail")
	}
}

func SaveConfig() {
	data2, _ := yaml.Marshal(conf)
	ioutil.WriteFile(default_config_path, data2, 0660)
}
