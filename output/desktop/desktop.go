package desktop

import (
	"context"
	"errors"
	"github.com/injoyai/logs"
	"github.com/injoyai/notice/output"
	"github.com/injoyai/notice/output/tcp"
)

func Init() {

	f := func(name string, msg *output.Message, Type string) error {
		//给桌面端发送消息
		logs.Tracef("给桌面端[%s]发送消息[%s]\n", name, msg.Content)

		c := tcp.Server.GetClient(name)
		if c == nil {
			logs.Warnf("给桌面端[%s]发送消息错误： 客户端不在线\n", name)
			return errors.New("客户端不在线")
			//return
		}
		d := msg.Details()
		d.Type = Type
		return c.WriteAny(d)
	}

	output.Trunk.Subscribe(func(ctx context.Context, msg *output.Message) {
		msg.Listen(map[string]func(name string, msg *output.Message) error{
			output.TypeDesktopNotice: func(name string, msg *output.Message) error {
				//给桌面端发送消息
				return f(name, msg, output.WinTypeNotice)
			},
			output.TypeDesktopVoice: func(name string, msg *output.Message) error {
				//给桌面端发送语音
				return f(name, msg, output.WinTypeVoice)
			},
			output.TypeDesktopPopup: func(name string, msg *output.Message) error {
				//给桌面端发送弹窗
				return f(name, msg, output.WinTypePopup)
			},
		})
	})

}
