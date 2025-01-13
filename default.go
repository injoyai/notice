package notice

import (
	cfg2 "github.com/injoyai/conv/cfg/v2"
	"github.com/injoyai/conv/codec"
	"github.com/injoyai/logs"
	"github.com/injoyai/notice/pkg/api"
	"github.com/injoyai/notice/pkg/middle"
	"github.com/injoyai/notice/pkg/push"
	"github.com/injoyai/notice/pkg/push/desktop"
	"github.com/injoyai/notice/pkg/push/dingtalk"
	"github.com/injoyai/notice/pkg/push/gotify"
	"github.com/injoyai/notice/pkg/push/local"
	"github.com/injoyai/notice/pkg/push/plugin"
	"github.com/injoyai/notice/pkg/push/pushplus"
	"github.com/injoyai/notice/pkg/push/script"
	"github.com/injoyai/notice/pkg/push/serverchan"
	"github.com/injoyai/notice/pkg/push/sms"
	"github.com/injoyai/notice/pkg/push/telegram"
	"github.com/injoyai/notice/pkg/push/webhook"
	"github.com/injoyai/notice/pkg/user"
	"path/filepath"
)

const (
	DefaultDesktopPort = 8427
	DefaultHTTPPort    = 8428
)

func Default(dataDir string) {

	cfg := cfg2.New(cfg2.WithFile(filepath.Join(dataDir, "/config/config.yaml"), codec.Yaml))

	//加载短信
	_alisms, _ := sms.NewAliyun(&sms.AliyunConfig{})

	//加载桌面端通知
	_desktop, _ := desktop.New(cfg.GetInt("desktop.port", DefaultDesktopPort))

	//telegram
	_telegram, _ := telegram.New(cfg.GetString("telegram.token"))

	//注册pusher
	push.Manager.Register(
		_alisms,
		_desktop,
		_telegram,

		webhook.New(nil), //webhook

		plugin.New(),        //插件
		script.New(20, nil), //脚本
		local.New(),         //本机

		//pushplus
		pushplus.New(
			cfg.GetString("pushplus.token"),
			NewClient(cfg.GetDuration("pushplus.timeout")),
		),

		//server酱
		serverchan.New(
			cfg.GetString("serverchan.sendkey"),
			NewClient(cfg.GetDuration("serverchan.timeout")),
		),

		//gotify
		gotify.New(
			cfg.GetString("gotify.host"),
			cfg.GetString("gotify.token"),
			cfg.GetInt("gotify.priority", 0),
			NewClient(cfg.GetDuration("gotify.timeout")),
		),

		//钉钉
		dingtalk.New(
			cfg.GetString("dingtalk.url"),
			cfg.GetString("dingtalk.secret"),
			NewClient(cfg.GetDuration("dingtalk.timeout")),
		),
	)

	//消息中间件
	push.Manager.Use(
		middle.NewLog(),  //日志
		middle.NewAuth(), //校验权限
		middle.NewForbidden(cfg.GetStrings("forbidden.words")...), //校验违禁词
		middle.NewQueue(10, cfg.GetDuration("queue.timeout")),     //消息队列
	)

	//加载用户
	logs.PanicErr(user.Init(dataDir))

	//加载http服务
	api.Init(cfg.GetInt("http.port", DefaultHTTPPort))
}
