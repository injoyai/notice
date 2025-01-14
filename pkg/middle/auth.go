package middle

import (
	"fmt"
	"github.com/injoyai/notice/pkg/push"
)

func NewAuth(enable ...bool) *Auth {
	return &Auth{
		enable: len(enable) == 0 || enable[0],
	}
}

// Auth 权限校验
type Auth struct {
	enable bool
}

func (this *Auth) Handler(u push.User, msg *push.Message, next func() error) error {
	if !u.Limits(msg.Method) {
		return fmt.Errorf("没有该推送权限")
	}
	return next()
}
