package sms

import (
	"errors"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/injoyai/conv"
	"github.com/injoyai/notice/pkg/push"
	"strings"
)

func Push(cfg *AliyunConfig, code, param string) error {
	c, err := NewAliyun(cfg)
	if err != nil {
		return err
	}
	return c.Push(Message(code, param))
}

func Message(code, param string, phone ...string) *push.Message {
	return &push.Message{
		Target:  strings.Join(phone, ","),
		Title:   code,
		Content: param,
	}
}

func NewAliyun(cfg *AliyunConfig) (*Aliyun, error) {
	config := sdk.NewConfig()
	credential := credentials.NewAccessKeyCredential(cfg.SecretID, cfg.SecretKey)
	regionId := conv.Select[string](len(cfg.RegionID) == 0, "cn-hangzhou", cfg.RegionID)
	client, err := dysmsapi.NewClientWithOptions(regionId, config, credential)
	return &Aliyun{
		cfg:    cfg,
		Client: client,
	}, err
}

type AliyunConfig struct {
	SecretID  string `json:"secretID"`  //
	SecretKey string `json:"secretKey"` //
	SignName  string `json:"signName"`  //签名
	RegionID  string `json:"regionId"`  //地域,如cn-hangzhou

	DefaultPhones string `json:"defaultPhones"` //默认号码
}

type Aliyun struct {
	cfg *AliyunConfig
	*dysmsapi.Client
}

func (this *Aliyun) Name() string {
	return "阿里云短信"
}

func (this *Aliyun) Types() []string {
	return []string{push.TypeAliyunSMS}
}

func (this *Aliyun) Push(msg *push.Message) error {
	phones := conv.Select(msg.Target != "", msg.Target, this.cfg.DefaultPhones)
	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"
	request.SignName = this.cfg.SignName
	request.TemplateCode = msg.Title
	request.PhoneNumbers = phones
	request.TemplateParam = msg.Content
	response, err := this.Client.SendSms(request)
	if err != nil {
		return err
	}
	if !response.IsSuccess() {
		return errors.New(response.Message)
	}
	if response.Code != "OK" {
		return errors.New(response.Message)
	}
	return nil
}
