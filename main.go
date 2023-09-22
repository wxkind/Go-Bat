package main

import (
	"Go-Bat/message"
)

func main() {

	//创建机器人对象
	GoBat := message.NewGoBat()
	//异步启动监听消息
	GoBat.Start()
	//	netstat -tunlp | grep 5700 5700端口占用进程n
	//去掉调试信息，减小可执行文件体积
	//	go build -ldflags="-s -w" -a -v -o xxx main.go

}
