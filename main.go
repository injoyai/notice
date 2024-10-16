package main

import (
	"github.com/injoyai/conv/cfg"
	"github.com/injoyai/conv/codec"
	"github.com/injoyai/logs"
	"github.com/injoyai/notice/input/http"
	"github.com/injoyai/notice/output/wechat"
)

func init() {
	cfg.Init("./config/config.yaml", codec.Yaml)
}

func main() {

	//加载微信通知
	logs.PanicErr(wechat.Init())

	//加载http服务
	http.Init(cfg.GetInt("http.port", 8426))
}
