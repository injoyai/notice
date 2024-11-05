package tcp

import (
	"encoding/json"
	"github.com/injoyai/conv"
	"github.com/injoyai/ios"
	"github.com/injoyai/ios/client"
	"github.com/injoyai/logs"
	"github.com/injoyai/notice/input/forbidden"
	"github.com/injoyai/notice/output"
	"github.com/injoyai/notice/user"
)

func DealMessage(c *client.Client, msg ios.Acker) {

	data := new(output.Message)
	if err := json.Unmarshal(msg.Payload(), data); err != nil {
		logs.Warn("解析消息失败", err)
		return
	}

	err := func() error {
		//从缓存获取用户
		u, err := user.GetByCache(c.GetKey())
		if err != nil {
			return err
		}
		//校验发送权限
		if err := data.Check(u.LimitMap()); err != nil {
			return err
		}
		//检查违禁词
		if err := forbidden.Forbidden.Check(data.Content); err != nil {
			return err
		}
		//发送队列
		_, err = output.Trunk.Do(data)
		return err
	}()

	c.WriteAny(output.Resp{
		ID:      data.ID,
		Success: err == nil,
		Message: conv.String(err),
	})

}
