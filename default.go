package notice

import (
	"errors"
	cfg2 "github.com/injoyai/conv/cfg/v2"
	"github.com/injoyai/conv/codec"
	"github.com/injoyai/logs"
	"github.com/injoyai/notice/pkg/api"
	"github.com/injoyai/notice/pkg/middle/forbidden"
	"github.com/injoyai/notice/pkg/push"
	"github.com/injoyai/notice/pkg/push/desktop"
	"github.com/injoyai/notice/pkg/push/gotify"
	"github.com/injoyai/notice/pkg/push/plugin"
	"github.com/injoyai/notice/pkg/push/pushplus"
	"github.com/injoyai/notice/pkg/push/script"
	"github.com/injoyai/notice/pkg/push/serverchan"
	"github.com/injoyai/notice/pkg/push/sms"
	"github.com/injoyai/notice/pkg/push/webhook"
	"github.com/injoyai/notice/pkg/push/wechat"
	"github.com/injoyai/notice/pkg/user"
	"path/filepath"
)

func Default(dataDir string) {

	cfg := cfg2.New(cfg2.WithFile(filepath.Join(dataDir, "/config/config.yaml"), codec.Yaml))

	//加载短信
	_sms, err := sms.NewAliyun(&sms.AliyunConfig{})
	logs.PanicErr(err)

	//加载微信通知
	_wechat, err := wechat.New(dataDir)
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

	//server酱
	_serverchan := serverchan.New(cfg.GetString("serverchan.sendkey"))

	//插件
	_plugin := plugin.New()

	//脚本
	_script := script.New(20, nil)

	//消息中间件
	push.Manager.Use(func(u *user.User, msg *push.Message) error {
		//校验权限
		limit := u.LimitMap()
		if len(limit) > 0 {
			if _, ok := limit[push.TypeAll]; !ok {
				if _, ok2 := limit[msg.Method]; !ok2 {
					return errors.New("无权限")
				}
			}
		}
		//校验违禁词
		return f.Check(msg.Title, msg.Content)
	})

	//注册pusher
	push.Manager.Register(
		_sms,
		_wechat,
		_gotify,
		_desktop,
		_webhook,
		_pushplus,
		_serverchan,
		_plugin,
		_script,
	)

	//加载用户
	logs.PanicErr(user.Init(dataDir))

	//加载http服务
	api.Init(cfg.GetInt("http.port", 8426))
}
