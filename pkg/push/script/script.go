package script

import (
	"errors"
	"github.com/injoyai/goutil/script"
	"github.com/injoyai/goutil/script/js"
	"github.com/injoyai/notice/pkg/push"
)

func New(pool int, m map[string]string) *Script {
	if m == nil {
		m = make(map[string]string)
	}
	return &Script{
		js: js.NewPool(pool),
		m:  m,
	}
}

type Script struct {
	js *js.Pool
	m  map[string]string
}

func (this *Script) Name() string {
	return "脚本"
}

func (this *Script) Types() []string {
	return []string{push.TypeScript}
}

func (this *Script) Push(msg *push.Message) error {
	s, ok := this.m[msg.Target]
	if !ok {
		return errors.New("脚本不存在: " + msg.Target)
	}
	_, err := this.js.Exec(s, func(i script.Client) {
		i.Set("title", msg.Title)
		i.Set("content", msg.Content)
	})
	return err
}
