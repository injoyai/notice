package tcp

import (
	"encoding/json"
	"github.com/injoyai/ios"
	"github.com/injoyai/ios/client"
	"github.com/injoyai/logs"
	"github.com/injoyai/notice/output"
)

func DealMessage(c *client.Client, msg ios.Acker) {
	data := new(output.Message)
	if err := json.Unmarshal(msg.Payload(), data); err != nil {
		logs.Warn("解析消息失败", err)
		return
	}
	output.Trunk.Do(data)
}
