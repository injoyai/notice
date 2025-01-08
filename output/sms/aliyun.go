package sms

import (
	"github.com/injoyai/goutil/notice"
	"github.com/injoyai/notice/output"
)

type AliyunConfig = notice.AliyunConfig

func NewAliyun(config *AliyunConfig) (*Aliyun, error) {
	i, err := notice.NewAliyunSMS(config)
	if err != nil {
		return nil, err
	}
	return &Aliyun{Interface: i}, nil
}

type Aliyun struct {
	notice.Interface
}

func (this *Aliyun) Types() []string {
	return []string{output.TypeAliyunSMS}
}

func (this *Aliyun) Push(msg *output.Message) error {
	return this.Publish(&notice.Message{
		Target: msg.Target,
		Param: map[string]interface{}{
			"TemplateID": msg.Target,
			"Param":      msg.Content,
		},
	})
}
