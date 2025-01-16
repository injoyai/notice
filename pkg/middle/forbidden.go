package middle

import (
	"fmt"
	"github.com/injoyai/notice/pkg/push"
	"strings"
)

func NewForbidden(words ...string) *Forbidden {
	return &Forbidden{
		Words: words,
	}
}

// Forbidden 简易违禁词中间件
type Forbidden struct {
	Words []string
}

func (this *Forbidden) Handler(msg *push.Message, u push.User, f func() error) error {
	for _, v := range this.Words {
		if strings.Contains(msg.Title, v) {
			return fmt.Errorf("标题包含违禁词:%s", v)
		}
		if strings.Contains(msg.Content, v) {
			return fmt.Errorf("内容包含违禁词:%s", v)
		}
	}
	return f()
}
