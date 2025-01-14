package middle

import (
	"github.com/injoyai/notice/pkg/push"
	"strings"
)

func NewSplit(sep string) *Split {
	if len(sep) == 0 {
		sep = ","
	}
	return &Split{
		Sep: sep,
	}
}

type Split struct {
	Sep string
}

func (this *Split) Handler(u push.User, msg *push.Message, next func() error) error {
	method := strings.Split(msg.Method, this.Sep)
	for _, v := range method {
		msg.Method = v
		if err := next(); err != nil {
			return err
		}
	}
	return nil
}
