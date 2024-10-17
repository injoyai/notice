package main

import (
	_ "embed"
	"fmt"
	"github.com/getlantern/systray"
	"github.com/injoyai/goutil/oss"
	"github.com/injoyai/lorca"
	"github.com/injoyai/notice/output"
	"os"
	"path/filepath"
)

//go:embed index.html
var html string

var (
	App lorca.APP
	TCP = NewTCP()
)

func main() {
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

		TCP.Rerun.OnDial(func(index, retry int, err error) {
			app.Eval(fmt.Sprintf("showLogin(%v,'%s','%s','%s')",
				err != nil,
				TCP.Cache.GetString("address"),
				TCP.Cache.GetString("username"),
				TCP.Cache.GetString("password"),
			))
		})

		app.Bind("init", func() {
			app.Eval(fmt.Sprintf("showLogin(%v,'%s','%s','%s')",
				TCP.Closed(),
				TCP.Cache.GetString("address"),
				TCP.Cache.GetString("username"),
				TCP.Cache.GetString("password"),
			))
		})

		app.Bind("fnLogin", func(address, username, password string) {
			if err := TCP.Update(address, username, password); err != nil {
				app.Eval(fmt.Sprintf("notice('%s')", err.Error()))
				return
			}
			app.Eval("notice('成功')")
			app.Eval("showPush(true)")
		})

		app.Bind("fnPush", func(method, target, Type, content string) {
			err := TCP.WriteAny(output.Message{
				Output:  []string{method + ":" + target},
				Type:    Type,
				Content: content,
			})
			if err != nil {
				app.Eval(fmt.Sprintf("notice('%s')", err.Error()))
				return
			}
			app.Eval("notice('成功')")
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

	openUI()

}

func onExit() {
	if App != nil {
		App.Close()
	}
}
