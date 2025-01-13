package middle

import (
	"github.com/injoyai/conv"
	"github.com/injoyai/logs"
	"github.com/injoyai/notice/pkg/push"
)

func NewLog() *Log {
	return &Log{}
}

type Log struct{}

func (this *Log) Handler(u push.User, msg *push.Message, next func() error) error {
	err := next()
	logs.Debugf("用户[%s]推送消息[%s][%s], 结果: %s", u.GetName(), msg.Method, msg.Title, conv.New(err).String("成功"))
	return err
}
