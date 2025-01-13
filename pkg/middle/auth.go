package middle

import (
	"errors"
	"github.com/injoyai/notice/pkg/push"
)

func NewAuth() *Auth {
	return &Auth{}
}

// Auth 权限校验
type Auth struct{}

func (this *Auth) Handler(u push.User, msg *push.Message, next func() error) error {
	if !u.Limits(msg.Method) {
		return errors.New("无权限")
	}
	return next()
}
