package push

import (
	"github.com/injoyai/conv"
	"github.com/injoyai/goutil/net/http"
	"time"
)

type HTTPClientOption func(c *http.Client)

func WithTimeout(timeout ...time.Duration) HTTPClientOption {
	return func(c *http.Client) {
		c.SetTimeout(conv.DefaultDuration(0, timeout...))
	}
}

func WithProxy(proxy string) HTTPClientOption {
	return func(c *http.Client) {
		c.SetProxy(proxy)
	}
}
