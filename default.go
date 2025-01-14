package notice

import (
	"github.com/injoyai/conv"
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
	user "github.com/injoyai/notice/pkg/user/default"
	"github.com/redis/go-redis/v9"
	"path/filepath"
	"time"
)

const (
	DefaultDesktopPort = 8427
	DefaultHTTPPort    = 8426
)

func Default(dataDir string) {

	cfg := cfg2.New(cfg2.WithFile(filepath.Join(dataDir, "/config/config_real.yaml"), codec.Yaml))

	//加载短信
	_alisms, _ := sms.NewAliyun(&sms.AliyunConfig{})

	//桌面端通知
	_desktop, _ := desktop.New(
		cfg.GetInt("desktop.port", DefaultDesktopPort),
		cfg.GetBool("desktop.enable", true),
	)

	//telegram
	_telegram, _ := telegram.New(cfg.GetString("telegram.token"))

	//注册pusher
	push.Manager.Register(
		_alisms,      //阿里短信
		_desktop,     //桌面端
		_telegram,    //telegram
		plugin.New(), //插件
		local.New(),  //本机

		//webhook
		webhook.New(func() map[string]*webhook.Config {
			m := make(map[string]*webhook.Config)
			for k, v := range cfg.GetGMap("webhook") {
				x := new(webhook.Config)
				conv.Unmarshal(v, x)
				m[k] = x
			}
			return m
		}()),

		//脚本
		script.New(
			cfg.GetInt("script.pool", 10),
			cfg.GetSMap("script.content"),
		),

		//pushplus
		pushplus.New(
			cfg.GetString("pushPlus.token"),
			NewClient(cfg.GetDuration("pushPlus.timeout")),
		),

		//server酱
		serverchan.New(
			cfg.GetString("serverChan.sendKey"),
			NewClient(cfg.GetDuration("serverChan.timeout")),
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
		//消息队列
		middle.NewQueue(
			cfg.GetInt("middle.queue.limit", 10),
			cfg.GetInt("middle.queue.cap", 20),
			cfg.GetDuration("middle.queue.timeout", time.Second*10),
		),
		middle.NewForbidden(cfg.GetStrings("middle.forbidden.words")...), //校验违禁词
		middle.NewAuth(cfg.GetBool("middle.auth.enable")),                //校验权限
		middle.NewLog(cfg.GetBool("middle.log.enable")),                  //日志
	)

	//加载用户
	logs.PanicErr(user.Init(&user.Config{
		Type: cfg.GetString("user.database.type", "sqlite"),
		DSN:  cfg.GetString("user.database.dsn", filepath.Join(dataDir, user.Filename)),
		Auth: user.AuthConfig{
			Enable: cfg.GetBool("user.auth.enable"),
			Type:   cfg.GetString("user.auth.type", user.Memory),
			Redis: &redis.Options{
				Addr:     cfg.GetString("user.auth.redis.addr", "127.0.0.1:6379"),
				Password: cfg.GetString("user.auth.redis.password"),
				DB:       cfg.GetInt("user.auth.redis.db"),
			},
			SuperToken: cfg.GetStrings("user.auth.super"),
		},
	}))

	//加载http服务
	api.Init(cfg.GetInt("http.port", DefaultHTTPPort))
}
