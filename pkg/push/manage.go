package push

import (
	"fmt"
	"github.com/injoyai/base/g"
	"github.com/injoyai/conv"
	"io"
	"net/http"
)

var Manager = NewManage()

func NewManage() *Manage {
	return &Manage{
		pusher: map[string][]Pusher{},
		middle: []Middle{},
	}
}

type Manage struct {
	pusher map[string][]Pusher
	middle []Middle
}

func (this *Manage) Handler(r *http.Request, u User) error {
	defer r.Body.Close()
	bs, _ := io.ReadAll(r.Body)
	msg := &Message{}
	err := conv.Unmarshal(bs, msg)
	if err != nil {
		return err
	}
	return this.Push(u, msg)
}

// Use 中间件,越后面添加的越先执行,类似洋葱,一层层包起来
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

// Push 推送消息
func (this *Manage) Push(u User, msg *Message) (err error) {
	defer g.Recover(&err)

	pushs, ok := this.pusher[msg.Method]
	if !ok {
		return fmt.Errorf("推送方式[%s]未注册", msg.Method)
	}

	for _, p := range pushs {
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
