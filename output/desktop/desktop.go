package desktop

import (
	"context"
	"github.com/injoyai/logs"
	"github.com/injoyai/notice/output"
	"github.com/injoyai/notice/output/tcp"
	"strings"
)

func Init() {
	output.Trunk.Subscribe(func(ctx context.Context, data interface{}) {
		msg := data.(*output.Message)
		for _, out := range msg.Output {
			if strings.HasPrefix(out, output.TypeTCPDesktop+":") {
				name, _ := strings.CutPrefix(out, output.TypeTCP+":")

				//给桌面端发送消息
				logs.Tracef("给桌面端[%s]发送消息[%s]\n", name, msg.Content)

				c := tcp.Server.GetClient(name)
				if c == nil {
					logs.Warnf("给桌面端[%s]发送消息错误： 客户端不在线\n", name)
					return
				}
				c.WriteAny(msg.Details())
			}
		}
	})

}
