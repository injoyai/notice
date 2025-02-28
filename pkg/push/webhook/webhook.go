package webhook

import (
	"github.com/injoyai/goutil/net/http"
	"github.com/injoyai/goutil/script/js"
	"github.com/injoyai/notice/pkg/push"
	"strings"
)

func New(m map[string]*Config) *Webhook {
	return &Webhook{
		m:  m,
		js: js.New(),
	}
}

type Webhook struct {
	m      map[string]*Config
	js     *js.Client
	Client *http.Client //
}

func (this *Webhook) Name() string {
	return "Webhook"
}

func (this *Webhook) Types() []string {
	return []string{push.TypeWebhook}
}

func (this *Webhook) Push(msg *push.Message) error {
	w, ok := this.m[msg.Target]
	if !ok {
		if this.Client == nil {
			this.Client = http.NewClient()
		}
		err := this.Client.Url(msg.Target).SetBody(msg.Content).Post().Err()
		return err
	}
	s := strings.ReplaceAll(w.Body, "${title}", msg.Title)
	s = strings.ReplaceAll(s, "${content}", msg.Content)

	//if after, ok := strings.CutPrefix(s, "//js"); ok {
	//	result, err := this.js.Exec(after, func(client script.Client) {
	//		client.Set("title", msg.Title)
	//		client.Set("content", msg.Content)
	//		client.Set("method", msg.Method)
	//		client.Set("type", msg.Type)
	//		client.Set("target", msg.Target)
	//		client.Set("time", msg.Time)
	//	})
	//	if err != nil {
	//		return err
	//	}
	//	s = conv.String(result)
	//}

	_, err := w.do(s)
	return err
}

type Config struct {
	Url    string            //
	Method string            //
	Header map[string]string //
	Body   string            //内容 {"title":${title},"content":${content}}
	Retry  uint              //重试次数
	Client *http.Client      //
}

func (this Config) do(content string) (string, error) {
	if this.Client == nil {
		this.Client = http.NewClient()
	}
	if len(this.Method) == 0 {
		this.Method = http.MethodPost
	}
	x := this.Client.Url(this.Url)
	for k, v := range this.Header {
		x.SetHeader(k, v)
	}
	resp := x.SetBody(content).Retry(this.Retry).SetMethod(this.Method).Do()
	if resp.Err() != nil {
		return "", resp.Err()
	}
	return resp.GetBodyString(), nil
}
