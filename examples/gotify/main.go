package main

import (
	"github.com/injoyai/conv/cfg"
	"github.com/injoyai/logs"
	"github.com/injoyai/notice/pkg/push"
	"github.com/injoyai/notice/pkg/push/gotify"
)

func main() {
	cfg.Init(cfg.WithFile("./config/config_real.yaml"))

	p := gotify.New(
		cfg.GetString("gotify.host"),
		cfg.GetString("gotify.token"),
		cfg.GetInt("gotify.priority"),
	)

	push.Manager.Register(p)

	err := push.Manager.Push(&push.Message{
		Method:  push.TypeGotify,
		Title:   "标题",
		Content: "内容",
	}, nil)

	logs.PrintErr(err)
}
