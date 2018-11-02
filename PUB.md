# 部署说明

## 环境

将代码放在 GOPATH 目录


## 编译

服务器分为游戏服务器和登录服务器

游戏服务器编译命令：

go build -o scmjgsvr

cd player/login
go build -o scmjlsvr 

跨平台编译 加前缀： CGO_ENABLED=0 GOOS=linux GOARCH=amd64 
如 : CGO_ENABLED=0 GOOS=linux GOARCH=amd64  go build -o scmjgsvr


会在当前目录下生成 scmjgsvr (配合当前目录下的 config.toml 配置)
在 player/login 目录生成 scmjlsvr (配合其目录下生成 config.toml 配置)

将 scmjgsvr 和 scmjlsvr 上传到服务器 （包括其配置文件）
程序和其配置放在同一目录下，不同程序放不同目录


服务器运行 (在所在目录)：

./scmjgsvr &
./scmjlsvr &


查看服务器：

ps -ef | grep scmj

关闭服务器： 使用 kill 关闭通过 ps 查看的进程即可



配置文件说明：

主要核查游戏端口号
