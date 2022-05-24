package main

import (
	"gin-memos/conf"
	"gin-memos/route"
)

func main() {
	conf.Init() // MySQL初始化
	r := route.NewRouter()
	r.Run(conf.HttpPort) // 配置文件中给定了是8080端口
}
