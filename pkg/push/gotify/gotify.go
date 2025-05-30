package gotify

import (
	"errors"
	"fmt"
	"github.com/injoyai/goutil/g"
	"github.com/injoyai/goutil/net/http"
	"github.com/injoyai/notice/pkg/push"
)

func New(host, token string, priority int, client ...*http.Client) *Gotify {
	x := &Gotify{
		Host:     host,
		Token:    token,
		Priority: priority,
		client:   http.DefaultClient,
	}
	if len(client) > 0 && client[0] != nil {
		x.client = client[0]
	}
	return x
}

type Gotify struct {
	Host     string
	Token    string
	Priority int
	client   *http.Client
}

func (this *Gotify) Name() string {
	return "Gotify"
}

func (this *Gotify) Types() []string {
	return []string{push.TypeGotify}
}

func (this *Gotify) Push(msg *push.Message) error {
	if this.Host == "" {
		return errors.New("无效的Gotify推送Host")
	}
	if this.Token == "" {
		return errors.New("无效的Gotify推送Token")
	}
	return this.client.Url(fmt.Sprintf("%s/message?token=%s", this.Host, this.Token)).
		SetHeader("Content-Type", "application/json").
		SetBody(g.Map{
			"title":    msg.Title,
			"message":  msg.Content,
			"priority": this.Priority,
		}).Debug().Post().Err()
}
