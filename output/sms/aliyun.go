package sms

import (
	"context"
	"github.com/injoyai/conv/cfg"
	"github.com/injoyai/goutil/notice"
	"github.com/injoyai/logs"
	"github.com/injoyai/notice/output"
	"strings"
)

func Init() {
	Aliyun, err := notice.NewAliyunSMS(&notice.AliyunConfig{
		SecretID:  cfg.GetString("sms.aliyun.secretID"),
		SecretKey: cfg.GetString("sms.aliyun.secretKey"),
		SignName:  cfg.GetString("sms.aliyun.signName"),
		RegionID:  cfg.GetString("sms.aliyun.regionID"),
	})
	logs.PanicErr(err)

	output.Trunk.Subscribe(func(ctx context.Context, msg *output.Message) {
		msg.On(output.TypeAliyunSMS, func(name string, msg *output.Message) {

			list := strings.SplitN(name, ":", 2)
			if len(list) == 2 {
				err := Aliyun.Publish(&notice.Message{
					Target: list[1],
					Param: map[string]interface{}{
						"TemplateID": list[0],
						"Param":      msg.Content,
					},
				})
				logs.PrintErr(err)
			}

		})
	})
}
