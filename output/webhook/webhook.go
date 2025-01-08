package webhook

import (
	"errors"
	"github.com/injoyai/goutil/net/http"
	"github.com/injoyai/notice/output"
	"strings"
	"time"
)

func New(m map[string]*Config) *Webhook {
	return &Webhook{m: m}
}

type Webhook struct {
	m map[string]*Config
}

func (this *Webhook) Types() []string {
	return []string{output.TypeWebhook}
}

func (this *Webhook) Push(msg *output.Message) error {
	w, ok := this.m[msg.Target]
	if !ok {
		return errors.New("webhook不存在: " + msg.Target)
	}
	s := strings.ReplaceAll(w.Body, "{title}", msg.Title)
	s = strings.ReplaceAll(s, "{content}", msg.Content)

	_, err := w.do()
	return err
}

type Config struct {
	Url     string            //
	Method  string            //
	Header  map[string]string //
	Body    string            //内容 {"title":{title},"content":{content}}
	Timeout time.Duration     //超时时间
	Proxy   string            //代理
	Retry   uint              //重试次数
	client  *http.Client
}

func (this Config) do() (string, error) {
	if this.client == nil {
		this.client = http.NewClient()
	}
	if this.Timeout > 0 {
		this.client.SetTimeout(this.Timeout)
	}
	if this.Proxy != "" {
		this.client.SetProxy(this.Proxy)
	}

	x := this.client.Url(this.Url)
	for k, v := range this.Header {
		x.SetHeader(k, v)
	}
	resp := x.SetBody(this.Body).Retry(this.Retry).SetMethod(this.Method).Do()
	if resp.Err() != nil {
		return "", resp.Err()
	}
	return resp.GetBodyString(), nil
}
