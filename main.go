package main

import (
	"github.com/injoyai/conv/cfg"
	"github.com/injoyai/conv/codec"
	"github.com/injoyai/logs"
	"github.com/injoyai/notice/input/forbidden"
	"github.com/injoyai/notice/input/http"
	in_tcp "github.com/injoyai/notice/input/tcp"
	"github.com/injoyai/notice/output/desktop"
	"github.com/injoyai/notice/output/sms"
	"github.com/injoyai/notice/output/tcp"
	"github.com/injoyai/notice/output/wechat"
	"github.com/injoyai/notice/user"
	"os"
	"path/filepath"
)

var DataDir = "./"

func init() {
	switch {
	case len(os.Args) > 1:
		DataDir = os.Args[1]
	default:
		DataDir = "./"
	}
	cfg.Init(filepath.Join(DataDir, "/config/config.yaml"), codec.Yaml)
}

func main() {

	//加载违禁词规则
	forbidden.Init()

	//加载短信
	sms.Init()

	//加载用户
	logs.PanicErr(user.Init(DataDir))

	//加载微信通知
	logs.PanicErr(wechat.Init(DataDir))

	//加载桌面端通知
	desktop.Init()

	//加载tcp服务
	go tcp.Init(cfg.GetInt("tcp.port", 8427), in_tcp.DealMessage)

	//加载http服务
	http.Init(cfg.GetInt("http.port", 8426))
}
