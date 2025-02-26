package notice

import (
	"embed"
	_ "embed"
	"github.com/injoyai/conv"
	cfg2 "github.com/injoyai/conv/cfg/v2"
	"github.com/injoyai/conv/codec"
	"github.com/injoyai/goutil/frame/in/v3"
	"github.com/injoyai/goutil/frame/mux"
	"github.com/injoyai/goutil/oss"
	"github.com/injoyai/logs"
	"github.com/injoyai/notice/pkg/middle"
	"github.com/injoyai/notice/pkg/push"
	"github.com/injoyai/notice/pkg/push/desktop"
	"github.com/injoyai/notice/pkg/push/dingtalk"
	"github.com/injoyai/notice/pkg/push/gotify"
	"github.com/injoyai/notice/pkg/push/local"
	"github.com/injoyai/notice/pkg/push/mail"
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
	Version            = "0.0.4"
	DefaultDesktopPort = 8427
	DefaultHTTPPort    = 8426
)

//go:embed config/config_example.yaml
var exampleConfig []byte

//go:embed dist
var web embed.FS

func Default(dataDir string) {

	logs.SetFormatter(logs.TimeFormatter)

	logs.Info("版本:", Version)

	filename := filepath.Join(dataDir, "config/config.yaml")

	oss.NewNotExist(filename, exampleConfig)

	cfg := cfg2.New(cfg2.WithFile(filename, codec.Yaml))

	//_wechat, _ := wechat.New(dataDir)

	//加载短信
	_alisms, _ := sms.NewAliyun(&sms.AliyunConfig{})

	//桌面端通知
	_desktop, _ := desktop.New(
		cfg.GetInt("desktop.port", DefaultDesktopPort),
		cfg.GetBool("desktop.enable", true),
	)

	//telegram
	_telegram, err := telegram.New(
		cfg.GetString("telegram.token"),
		cfg.GetString("telegram.proxy"),
		cfg.GetString("telegram.defaultChatID"),
	)
	logs.PrintErr(err)

	//注册pusher
	push.Manager.Register(
		_alisms,      //阿里短信
		_desktop,     //桌面端
		_telegram,    //telegram
		plugin.New(), //插件
		local.New(),  //本机

		//邮箱
		mail.New(&mail.Config{
			Host:     cfg.GetString("mail.host"),
			Port:     cfg.GetInt("mail.port", 25),
			Username: cfg.GetString("mail.username"),
			Password: cfg.GetString("mail.password"),
		}),

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
		middle.NewRetry(cfg.GetInt("middle.retry")), //重试
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
	err = HTTP(cfg.GetInt("http.port", DefaultHTTPPort))
	logs.Err(err)
}

func HTTP(port int) error {

	s := mux.New().SetPort(port)
	s.Group("/api", func(g *mux.Grouper) {

		//登录
		g.POST("/login", func(r *mux.Request) {
			req := &user.LoginReq{}
			r.Parse(req)
			token, err := user.Login(req)
			in.CheckErr(err)
			in.Succ(map[string]any{"token": token})
		})

		//校验权限
		g.Middle(func(r *mux.Request) {
			token := r.GetHeader("Authorization")
			if len(token) == 0 {
				token = r.GetQueryVar("token").String()
			}
			u, valid, err := user.CheckToken(token)
			in.CheckErr(err)
			if !valid {
				in.DefaultClient.Forbidden()
			}
			r.SetCache("user", u)
		})

		//发送消息
		g.ALL("/notice", func(r *mux.Request) {
			u := r.GetCache("user").Val().(*user.User)
			msg := &push.Message{}
			r.Parse(msg)
			//加入发送队列
			err := push.Manager.Push(msg, u)
			in.Err(err)
		})

		//查询用户列表
		g.GET("/user/all", func(r *mux.Request) {
			data, err := user.GetAll()
			in.CheckErr(err)
			in.Succ(data)
		})

		//添加/修改用户
		g.POST("/user", func(r *mux.Request) {
			req := new(user.User)
			r.Parse(req)
			err := user.Create(req)
			in.Err(err)
		})

		//删除用户
		g.DELETE("/user", func(r *mux.Request) {
			username := r.GetString("username")
			err := user.Del(username)
			in.Err(err)
		})

	})
	s.StaticEmbed("/", web, "dist")

	return s.Run()
}
