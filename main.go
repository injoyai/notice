package main

import (
	"errors"
	"github.com/injoyai/conv/cfg"
	"github.com/injoyai/conv/codec"
	"github.com/injoyai/logs"
	"github.com/injoyai/notice/input/forbidden"
	"github.com/injoyai/notice/input/http"
	"github.com/injoyai/notice/output"
	"github.com/injoyai/notice/output/desktop"
	"github.com/injoyai/notice/output/gotify"
	"github.com/injoyai/notice/output/plugin"
	"github.com/injoyai/notice/output/pushplus"
	"github.com/injoyai/notice/output/script"
	"github.com/injoyai/notice/output/sms"
	"github.com/injoyai/notice/output/webhook"
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

	//加载短信
	_sms, err := sms.NewAliyun(&sms.AliyunConfig{})
	logs.PanicErr(err)

	//加载微信通知
	_wechat, err := wechat.New(DataDir)
	logs.PanicErr(err)

	//gotify
	_gotify := gotify.New(
		cfg.GetString("gotify.host"),
		cfg.GetString("gotify.token"),
		cfg.GetInt("gotify.priority", 0),
	)

	//加载桌面端通知
	_desktop, err := desktop.New(cfg.GetInt("desktop.port", 8427))
	logs.PanicErr(err)

	//加载违禁词规则
	f := forbidden.New(cfg.GetStrings("forbidden.words")...)

	//webhook
	_webhook := webhook.New(nil)

	//pushplus
	_pushplus := pushplus.New(cfg.GetString("pushplus.token"))

	//插件
	_plugin := plugin.New()

	//脚本
	_script := script.New(20, nil)

	//消息中间件
	output.Manager.Use(func(u *user.User, msg *output.Message) error {
		//校验权限
		limit := u.LimitMap()
		if len(limit) > 0 {
			if _, ok := limit[output.TypeAll]; !ok {
				if _, ok2 := limit[msg.Method]; !ok2 {
					return errors.New("无权限")
				}
			}
		}
		//校验违禁词
		return f.Check(msg.Title, msg.Content)
	})

	//注册pusher
	output.Manager.Register(
		_sms,
		_wechat,
		_gotify,
		_desktop,
		_webhook,
		_pushplus,
		_plugin,
		_script,
	)

	//加载用户
	logs.PanicErr(user.Init(DataDir))

	//加载http服务
	http.Init(cfg.GetInt("http.port", 8426))
}
