package main

import (
	"github.com/injoyai/conv/cfg/v2"
	"github.com/injoyai/logs"
	"github.com/injoyai/notice/pkg/push"
	"github.com/injoyai/notice/pkg/push/mail"
)

func main() {

	cfg.Init(cfg.WithFile("./config/config_real.yaml"))

	w := mail.New(&mail.Config{
		Host:     cfg.GetString("mail.host"),
		Port:     cfg.GetInt("mail.port", 25),
		Username: cfg.GetString("mail.username"),
		Password: cfg.GetString("mail.password"),
	})

	push.Manager.Register(w)

	err := push.Manager.Push(&push.Message{
		Method:  push.TypeMail,
		Target:  cfg.GetString("mail.username"),
		Title:   "标题",
		Content: "内容",
	}, nil)

	logs.PrintErr(err)

}
