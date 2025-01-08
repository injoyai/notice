package output

import (
	"errors"
	"github.com/injoyai/notice/user"
)

var Manager = &Manage{}

type Manage struct {
	use    []func(from *user.User, msg *Message) error
	pusher map[string][]Interface
}

func (this *Manage) Use(i ...func(from *user.User, msg *Message) error) {
	this.use = append(this.use, i...)
}

func (this *Manage) Register(i ...Interface) {
	for _, v := range i {
		for _, _type := range v.Types() {
			this.pusher[_type] = append(this.pusher[_type], v)
		}
	}
}

func (this *Manage) Push(from *user.User, msg *Message) error {

	pushs, ok := this.pusher[msg.Method]
	if !ok {
		return errors.New("消息推送失败,推送方式不存在")
	}

	//中间件
	for _, v := range this.use {
		if err := v(from, msg); err != nil {
			return err
		}
	}

	//推送
	for _, p := range pushs {
		return p.Push(msg)
	}

	return nil
}
