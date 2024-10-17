package main

import (
	"github.com/injoyai/conv/cfg"
	"github.com/injoyai/conv/codec"
	"github.com/injoyai/logs"
	"github.com/injoyai/notice/input/http"
	in_tcp "github.com/injoyai/notice/input/tcp"
	"github.com/injoyai/notice/output/desktop"
	"github.com/injoyai/notice/output/tcp"
	"github.com/injoyai/notice/output/wechat"
	"github.com/injoyai/notice/user"
)

func init() {
	cfg.Init("./config/config.yaml", codec.Yaml)
}

func main() {

	//加载用户
	logs.PanicErr(user.Init())

	//加载微信通知
	logs.PanicErr(wechat.Init())

	//加载桌面端通知
	desktop.Init()

	//加载tcp服务
	go tcp.Init(cfg.GetInt("tcp.port", 8427), in_tcp.DealMessage)

	//加载http服务
	http.Init(cfg.GetInt("http.port", 8426))
}
