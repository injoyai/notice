package notice

import (
	"crypto/tls"
	"github.com/injoyai/conv"
	"github.com/injoyai/goutil/net/http"
	gohttp "net/http"
	"time"
)

func NewClient(timeout ...time.Duration) *http.Client {
	return &http.Client{
		Client: &gohttp.Client{
			Transport: &gohttp.Transport{
				DisableKeepAlives: true,
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
			Timeout: conv.Default[time.Duration](0, timeout...),
		},
	}
}
