package main

import (
	_ "embed"
	"fmt"
	"github.com/getlantern/systray"
	"github.com/injoyai/base/g"
	"github.com/injoyai/conv"
	"github.com/injoyai/goutil/cache"
	"github.com/injoyai/goutil/oss"
	"github.com/injoyai/lorca"
	"github.com/injoyai/notice/output"
	"os"
	"path/filepath"
	"time"
)

//go:embed index.html
var html string

var (
	App  lorca.APP
	TCP  = NewTCP()
	Push = cache.NewFile(oss.UserInjoyDir("notice/cache/push"))
)

func main() {
	////提示已经有进程在运行
	//result, err := shell.Exec("tasklist /FI 'imagename eq notice.exe'")
	//if err == nil && strings.Contains(result.String(), "notice.exe"){
	//
	//}

	//关闭老的进程
	//shell.Stop("notice.exe")

	go TCP.Rerun.Run(TCP)

	systray.Run(onReady, onExit)
}

func openUI() {
	lorca.Run(&lorca.Config{
		Width:  480,
		Height: 430,
		Html:   html,
	}, func(app lorca.APP) error {
		App = app

		TCP.onLogin = func() {
			app.Eval(fmt.Sprintf("showPush(true,'%s','%s','%s','%s')",
				Push.GetString("method"),
				Push.GetString("target"),
				Push.GetString("type"),
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
				app.Eval(fmt.Sprintf("showPush(true,'%s','%s','%s','%s')",
					Push.GetString("method"),
					Push.GetString("target"),
					Push.GetString("type"),
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
				app.Eval(fmt.Sprintf("showPush(true,'%s','%s','%s','%s')",
					Push.GetString("method"),
					Push.GetString("target"),
					Push.GetString("type"),
					"", //Push.GetString("content"),
				))
			}
		})

		app.Bind("fnPush", func(method, target, Type, content string) {
			Push.Set("method", method)
			Push.Set("target", target)
			Push.Set("type", Type)
			Push.Set("content", content)
			Push.Save()

			app.Eval("pushBefore()")
			id := g.RandString(16)
			err := TCP.WriteAny(output.Message{
				ID:      id,
				Output:  []string{method + ":" + target},
				Type:    Type,
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

func onReady() {

	systray.SetIcon(IcoNotice)
	systray.SetTooltip("消息通知")

	//显示菜单,这个库不能区分左键和右键,固设置了该菜单
	mShow := systray.AddMenuItem("显示", "显示界面")
	mShow.SetIcon(IcoMenuShow)
	go func() {
		for {
			<-mShow.ClickedCh
			//show会阻塞,多次点击无效
			openUI()
		}
	}()

	filename := oss.ExecName()
	name := filepath.Base(filename)
	startLnk := oss.UserStartupDir(name + ".lnk")
	startup := oss.Exists(startLnk)
	mStartup := systray.AddMenuItemCheckbox("自启", "开机自启", startup)
	go func() {
		for {
			<-mStartup.ClickedCh
			if mStartup.Checked() {
				os.Remove(startLnk)
			} else {
				Shortcut(oss.UserStartupDir(name+".lnk"), filename)
			}
			if oss.Exists(startLnk) {
				mStartup.Check()
			} else {
				mStartup.Uncheck()
			}
		}
	}()

	//退出菜单
	mQuit := systray.AddMenuItem("退出", "退出程序")
	mQuit.SetIcon(IcoMenuQuit)
	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
	}()

}

func onExit() {
	if App != nil {
		App.Close()
	}
}
