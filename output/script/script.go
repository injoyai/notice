package script

import (
	"errors"
	"github.com/injoyai/goutil/script/js"
	"github.com/injoyai/notice/output"
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

func (this *Script) Types() []string {
	return []string{output.TypeScript}
}

func (this *Script) Push(msg *output.Message) error {
	s, ok := this.m[msg.Target]
	if !ok {
		return errors.New("脚本不存在: " + msg.Target)
	}
	_, err := this.js.Exec(s)
	return err
}
