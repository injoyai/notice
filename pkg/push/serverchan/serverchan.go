package serverchan

import (
	"errors"
	"fmt"
	"github.com/injoyai/conv"
	"github.com/injoyai/goutil/net/http"
	"github.com/injoyai/notice/pkg/push"
	"net/url"
	"regexp"
	"strings"
)

func Push(sendKey string, title, content string) error {
	return New(sendKey).Push(push.NewMessage(title, content))
}

func New(sendKey string, client ...*http.Client) *ServerChan {
	return &ServerChan{
		DefaultSendKey: sendKey,
		client:         conv.Default[*http.Client](http.DefaultClient, client...),
	}
}

type ServerChan struct {
	DefaultSendKey string
	client         *http.Client
}

func (this *ServerChan) Name() string {
	return "Server酱"
}

func (this *ServerChan) Types() []string {
	return []string{push.TypeServerChan}
}

func (this *ServerChan) Push(msg *push.Message) error {
	sendKey := conv.Select[string](msg.Target != "", msg.Target, this.DefaultSendKey)
	if sendKey == "" {
		return errors.New("无效的Server酱推送SendKey")
	}
	data := url.Values{}
	data.Set("text", msg.Title)
	data.Set("desp", msg.Content)

	resp := this.client.Url(this.getApi(sendKey)).
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetBody(data.Encode()).Post()
	if resp.Err() != nil {
		return resp.Err()
	}

	m := resp.GetBodyDMap()
	if m.GetInt("code") != 0 {
		return errors.New(m.GetString("error"))
	}

	return nil
}

func (*ServerChan) getApi(sendKey string) string {
	// 根据 sendkey 是否以 "sctp" 开头决定 API 的 URL
	var apiUrl string
	if strings.HasPrefix(sendKey, "sctp") {
		// 使用正则表达式提取数字部分
		re := regexp.MustCompile(`sctp(\d+)t`)
		matches := re.FindStringSubmatch(sendKey)
		if len(matches) > 1 {
			num := matches[1]
			apiUrl = fmt.Sprintf("https://%s.push.ft07.com/send/%s.send", num, sendKey)
		} else {
			return "Invalid sendkey format for sctp"
		}
	} else {
		apiUrl = fmt.Sprintf("https://sctapi.ftqq.com/%s.send", sendKey)
	}
	return apiUrl
}
