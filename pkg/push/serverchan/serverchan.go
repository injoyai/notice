package serverchan

import (
	"errors"
	"fmt"
	"github.com/injoyai/goutil/net/http"
	"github.com/injoyai/notice/pkg/push"
	"net/url"
	"regexp"
	"strings"
)

func New(sendKey string, client ...*http.Client) *ServerChan {
	s := &ServerChan{
		SendKey: sendKey,
		Api:     getApi(sendKey),
		client:  http.DefaultClient,
	}
	if len(client) > 0 && client[0] != nil {
		s.client = client[0]
	}
	return s
}

type ServerChan struct {
	SendKey string
	Api     string
	client  *http.Client
}

func (this *ServerChan) Name() string {
	return "Server酱"
}

func (this *ServerChan) Types() []string {
	return []string{push.TypeServerChan}
}

func (this *ServerChan) Push(msg *push.Message) error {
	if this.SendKey == "" {
		return errors.New("无效的Server酱推送SendKey")
	}
	if this.Api == "" {
		return errors.New("无效的Server酱推送SendKey")
	}
	data := url.Values{}
	data.Set("text", msg.Title)
	data.Set("desp", msg.Content)

	return this.client.Url(this.Api).
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetBody(data.Encode()).Post().Err()
}

func getApi(sendKey string) string {
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
