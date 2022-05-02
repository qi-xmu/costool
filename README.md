# costool

## 使用截图

![img](https://secure2.wostatic.cn/static/3Q1hrmDFJRzUics4N44Hz/image.png?auth_key=1651510137-kVeK3wqydMBUzsSo3tM4iB-0-bd8260adde51e42c20c027a2f0ff6c7f)

## 使用方法

![img](https://secure2.wostatic.cn/static/c2hLcTwWXN5SB7HYQwvYcu/costool.png?auth_key=1651510736-deRfVk1bRfdYnDfsfFFs5G-0-4d3ff2392293907b0f0491163d0657b3) 

配置文件，需要放在用户目录下的`.costool.yaml`文件中格式如下：

```yaml
UserID: "xxxx"
SecretID: xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
SecretKey: xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
Token: ""
```

注意，键值之间存在空格。

## 安装方法

1. 下载预编译包

   [第一个完整版本 · @奇/costool - Gitee.com](https://gitee.com/qi_xmu/costool/releases/v1.0) 

2. 自行编译，需要安装go环境

  ```bash
  git clone https://gitee.com/qi_xmu/costool.git
  cd costool
  make && make install
  ```

## 其他

遇到bug或者问问题，可以提出issuse反馈。
