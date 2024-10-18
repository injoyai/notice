package sms

import (
	"context"
	"github.com/injoyai/conv"
	"github.com/injoyai/conv/cfg"
	"github.com/injoyai/goutil/notice"
	"github.com/injoyai/logs"
	"github.com/injoyai/notice/output"
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
			Aliyun.Publish(&notice.Message{
				Target:  name,
				Content: msg.Content,
				Param: map[string]interface{}{
					"TemplateID": conv.String(msg.Param["TemplateID"]),
					"Param":      conv.String(msg.Param["Param"]),
				},
			})
		})
	})
}
