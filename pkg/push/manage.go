package push

import (
	"errors"
	"github.com/injoyai/base/g"
)

var Manager = NewManage()

func NewManage() *Manage {
	m := &Manage{
		pusher: map[string][]Pusher{},
		middle: []Middle{},
	}
	m.Use(m)
	return m
}

type Manage struct {
	pusher map[string][]Pusher
	middle []Middle
}

// Use 中间件
func (this *Manage) Use(i ...Middle) *Manage {
	this.middle = append(this.middle, i...)
	return this
}

// Register 注册推送
func (this *Manage) Register(i ...Pusher) *Manage {
	for _, v := range i {
		for _, _type := range v.Types() {
			this.pusher[_type] = append(this.pusher[_type], v)
		}
	}
	return this
}

// Handler 实现中间件接口,校验推送方式是否存在
func (this *Manage) Handler(u User, msg *Message, f func() error) error {
	_, ok := this.pusher[msg.Method]
	if !ok {
		return errors.New("推送方式不存在")
	}
	return f()
}

// Push 推送消息
func (this *Manage) Push(u User, msg *Message) (err error) {
	defer g.Recover(&err)
	for _, p := range this.pusher[msg.Method] {
		if p == nil {
			continue
		}
		h := func(u User, msg *Message) error { return p.Push(msg) }
		err := this.doMiddle(h, 0)(u, msg)
		if err != nil {
			return err
		}
	}
	return
}

func (this *Manage) doMiddle(f Handler, index int) Handler {
	if index >= len(this.middle) {
		return func(u User, msg *Message) error {
			return f(u, msg)
		}
	}
	return this.doMiddle(func(u User, msg *Message) error {
		h := this.middle[index]
		if h == nil {
			return f(u, msg)
		}
		return h.Handler(u, msg, func() error {
			return f(u, msg)
		})
	}, index+1)
}
