package main

import (
	_ "embed"
	"fmt"
	"github.com/injoyai/base/g"
	"github.com/injoyai/conv"
	"github.com/injoyai/goutil/cache"
	"github.com/injoyai/goutil/oss"
	"github.com/injoyai/goutil/oss/tray"
	"github.com/injoyai/lorca"
	"github.com/injoyai/notice/pkg/push"
	"path/filepath"
	"time"
)

//go:embed index.html
var html string

var (
	dir  = oss.UserInjoyDir("notice/desktop/")
	TCP  = NewTCP(filepath.Join(dir, "login"), time.Second*10)
	Push = cache.NewFile(filepath.Join(dir, "push"))
)

func main() {
	go TCP.Rerun.Run(TCP)
	tray.Run(
		tray.WithShow(func(m *tray.Menu) { ui() }, tray.Icon(IcoShow)),
		tray.WithStartup(),
		tray.WithExit(),
		tray.WithIco(Ico),
		tray.WithHint("消息通知"),
	)
}

func ui() {
	lorca.Run(&lorca.Config{
		Width:  480,
		Height: 490,
		Html:   html,
	}, func(app lorca.APP) error {

		TCP.onLogin = func() {
			app.Eval(fmt.Sprintf("showPush(true,'%s','%s','%s','%s','%s')",
				Push.GetString("method"),
				Push.GetString("target"),
				Push.GetString("type"),
				Push.GetString("title"),
				Push.GetString("content"),
			))
		}
		TCP.onClose = func(err error) {
			app.Eval(fmt.Sprintf("showLogin(true,'%s','%s','%s')",
				TCP.Cache.GetString("address"),
				TCP.Cache.GetString("username"),
				TCP.Cache.GetString("password"),
			))
		}

		app.Bind("init", func() {
			if TCP.login {
				app.Eval(fmt.Sprintf("showPush(true,'%s','%s','%s','%s','%s')",
					Push.GetString("method"),
					Push.GetString("target"),
					Push.GetString("type"),
					Push.GetString("title"),
					Push.GetString("content"),
				))
			} else {
				app.Eval(fmt.Sprintf("showLogin(true,'%s','%s','%s')",
					TCP.Cache.GetString("address"),
					TCP.Cache.GetString("username"),
					TCP.Cache.GetString("password"),
				))
			}
		})

		app.Bind("fnLogin", func(address, username, password string) {
			app.Eval("loginBefore()")
			err := TCP.Update(address, username, password)
			app.Eval(fmt.Sprintf("loginAfter('%v')", conv.String(err)))
			if err == nil {
				app.Eval(fmt.Sprintf("showPush(true,'%s','%s','%s','%s','%s')",
					Push.GetString("method"),
					Push.GetString("target"),
					Push.GetString("type"),
					Push.GetString("title"),
					Push.GetString("content"),
				))
			}
		})

		app.Bind("fnPush", func(method, target, Type, title, content string) {
			Push.Set("method", method)
			Push.Set("target", target)
			Push.Set("type", Type)
			Push.Set("title", title)
			Push.Set("content", content)
			Push.Save()

			app.Eval("pushBefore()")
			id := g.RandString(16)
			err := TCP.WriteAny(push.Message{
				ID:      id,
				Method:  method,
				Target:  target,
				Type:    Type,
				Title:   title,
				Content: content,
				Time:    time.Now().Unix(),
			})
			if err == nil {
				_, err = TCP.wait.Wait(id)
			}
			app.Eval(fmt.Sprintf("pushAfter('%v')", conv.String(err)))
		})

		app.Bind("fnClose", func() {
			TCP.Close()
			app.Eval(fmt.Sprintf("showLogin(true,'%s','%s','%s')",
				TCP.Cache.GetString("address"),
				TCP.Cache.GetString("username"),
				TCP.Cache.GetString("password"),
			))
		})

		return nil
	})
}
