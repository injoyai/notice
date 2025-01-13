package sms

import (
	"github.com/injoyai/goutil/notice"
	"github.com/injoyai/notice/pkg/push"
)

type AliyunConfig = notice.AliyunConfig

func NewAliyun(config *AliyunConfig) (*Aliyun, error) {
	i, err := notice.NewAliyunSMS(config)
	if err != nil {
		return &Aliyun{Interface: i}, err
	}
	return &Aliyun{Interface: i}, nil
}

type Aliyun struct {
	notice.Interface
}

func (this *Aliyun) Name() string {
	return "阿里云短信"
}

func (this *Aliyun) Types() []string {
	return []string{push.TypeAliyunSMS}
}

func (this *Aliyun) Push(msg *push.Message) error {
	return this.Publish(&notice.Message{
		Target: msg.Target,
		Param: map[string]interface{}{
			"TemplateID": msg.Target,
			"Param":      msg.Content,
		},
	})
}
