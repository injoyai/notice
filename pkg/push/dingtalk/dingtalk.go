package dingtalk

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/injoyai/goutil/net/http"
	"github.com/injoyai/notice/pkg/push"
	"net/url"
	"strings"
	"time"
)

//"github.com/alibabacloud-go/dingtalk/robot_1_0"

func New(url, secret string, client ...*http.Client) *DingTalk {
	d := &DingTalk{
		URL:    url,
		Secret: secret,
		client: http.DefaultClient,
	}
	if len(client) > 0 && client[0] != nil {
		d.client = client[0]
	}
	return d
}

type DingTalk struct {
	URL    string       //webhook地址
	Secret string       //秘钥
	client *http.Client //客户端
}

func (this *DingTalk) Name() string { return "钉钉" }

func (this *DingTalk) Types() []string { return []string{push.TypeDingTalk} }

func (this *DingTalk) Push(msg *push.Message) error {
	// https://open.dingtalk.com/document/robots/custom-robot-access#title-72m-8ag-pqw
	messageRequest := request{
		MessageType: "text",
		Text: Text{
			Content: msg.Content,
		},
	}
	switch msg.Target {
	case "@all":
		messageRequest.At.IsAtAll = true
	case "":
	default:
		messageRequest.At.AtUserIds = strings.Split(msg.Target, "|")
	}

	timestamp := time.Now().UnixMilli()
	sign, err := this.sign(this.Secret, timestamp)
	if err != nil {
		return err
	}
	jsonData, err := json.Marshal(messageRequest)
	if err != nil {
		return err
	}

	return this.client.Url(fmt.Sprintf("%s&timestamp=%d&sign=%s", this.URL, timestamp, sign)).SetBody(jsonData).Post().Err()
}

func (this *DingTalk) sign(secret string, timestamp int64) (string, error) {
	// https://open.dingtalk.com/document/robots/customize-robot-security-settings
	// timestamp + key -> sha256 -> URL encode
	stringToSign := fmt.Sprintf("%d\n%s", timestamp, secret)
	h := hmac.New(sha256.New, []byte(secret))
	_, err := h.Write([]byte(stringToSign))
	if err != nil {
		return "", err
	}
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	signature = url.QueryEscape(signature)
	return signature, nil
}

type request struct {
	MessageType string   `json:"msgtype"`
	Text        Text     `json:"text"`
	Markdown    Markdown `json:"markdown"`
	At          At
}

type Text struct {
	Content string `json:"content"`
}

type Markdown struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

type At struct {
	AtUserIds []string `json:"atUserIds"`
	IsAtAll   bool     `json:"isAtAll"`
}

type response struct {
	Code    int    `json:"errcode"`
	Message string `json:"errmsg"`
}
