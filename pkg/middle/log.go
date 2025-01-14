package middle

import (
	"github.com/fatih/color"
	"github.com/injoyai/conv"
	"github.com/injoyai/logs"
	"github.com/injoyai/notice/pkg/push"
)

func NewLog(enable ...bool) *Log {
	return &Log{
		enable: len(enable) == 0 || enable[0],
		Entity: logs.New("日志").SetColor(color.FgYellow).SetFormatter(logs.TimeFormatter),
	}
}

type Log struct {
	enable bool
	*logs.Entity
}

func (this *Log) Handler(u push.User, msg *push.Message, next func() error) (err error) {
	defer func() {
		if this.enable {
			this.Printf("用户[%s] 方式[%s] 消息[%s] 结果: %s\n", u.GetName(), msg.Method, msg.Content, conv.New(err).String("成功"))
		}
	}()
	return next()
}
