package desktop

import (
	"context"
	"github.com/injoyai/logs"
	"github.com/injoyai/notice/output"
	"github.com/injoyai/notice/output/tcp"
)

func Init() {

	f := func(name string, msg *output.Message, Type string) {
		//给桌面端发送消息
		logs.Tracef("给桌面端[%s]发送消息[%s]\n", name, msg.Content)

		c := tcp.Server.GetClient(name)
		if c == nil {
			logs.Warnf("给桌面端[%s]发送消息错误： 客户端不在线\n", name)
			return
		}
		d := msg.Details()
		d.Type = Type
		c.WriteAny(d)
	}

	output.Trunk.Subscribe(func(ctx context.Context, msg *output.Message) {
		msg.Listen(map[string]func(name string, msg *output.Message){
			output.TypeDesktopNotice: func(name string, msg *output.Message) {
				//给桌面端发送消息
				f(name, msg, output.WinTypeNotice)
			},
			output.TypeDesktopVoice: func(name string, msg *output.Message) {
				//给桌面端发送语音
				f(name, msg, output.WinTypeVoice)
			},
			output.TypeDesktopPopup: func(name string, msg *output.Message) {
				//给桌面端发送弹窗
				f(name, msg, output.WinTypePopup)
			},
		})
	})

}
