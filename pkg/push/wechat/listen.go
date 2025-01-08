package wechat

import (
	"bufio"
	"github.com/eatmoreapple/openwechat"
	"github.com/injoyai/base/g"
	"github.com/injoyai/conv"
	"github.com/injoyai/goutil/net/http"
	"github.com/injoyai/logs"
	"io"
	"strings"
)

func (this *Wechat) DealMessage(msg *openwechat.Message) {
	if msg.IsText() {
		logs.Read(msg.Content)

		switch msg.Content {
		case "ping":
			msg.ReplyText("pong")
			return
		}

		after, ok := strings.CutPrefix(msg.Content, "@"+this.Self.NickName)
		if ok || msg.IsSendByFriend() {
			result, err := llama(after)
			if err != nil {
				msg.ReplyText("我遇到了一个错误：" + err.Error())
				logs.Warnf("llama处理消息错误： %v\n", err)
				return
			}
			msg.ReplyText(result)
		}

	}
}

func llama(req string) (string, error) {
	http.DefaultClient.SetTimeout(0)
	resp := http.Url("http://127.0.0.1:11434/api/generate").SetBody(g.Map{
		"model":  "llama3.2",
		"prompt": req,
	}).Debug().Post()
	if resp.Err() != nil {
		return "", resp.Err()
	}
	defer resp.Body.Close()
	buf := bufio.NewReader(resp.Body)
	result := ""
	bs := []byte(nil)
	for {
		b, err := buf.ReadByte()
		if err != nil {
			if err == io.EOF {
				return result, nil
			}
			return "", err
		}
		switch {
		case b == '{':
			bs = bs[:0]
			bs = append(bs, b)
		case b == '}':
			bs = append(bs, b)
			m := conv.NewMap(bs)
			result += m.GetString("response")
			if m.GetBool("done") {
				return result, nil
			}
		default:
			bs = append(bs, b)
		}
	}

}
