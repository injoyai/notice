package local

import (
	"errors"
	"github.com/injoyai/notice/pkg/push"
)

func New() *Local {
	return &Local{}
}

type Local struct{}

func (this *Local) Types() []string {
	return []string{}
}

func (this *Local) Name() string {
	return "本机"
}

func (this *Local) Push(msg *push.Message) error {
	return errors.New("暂不支持")
}
