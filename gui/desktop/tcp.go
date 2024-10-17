package main

import (
	"context"
	"encoding/json"
	"github.com/injoyai/base/safe"
	"github.com/injoyai/goutil/cache"
	"github.com/injoyai/goutil/notice"
	"github.com/injoyai/ios"
	"github.com/injoyai/ios/client"
	"github.com/injoyai/logs"
	"github.com/injoyai/notice/oem"
	"github.com/injoyai/notice/output"
	"github.com/injoyai/notice/user"
	"net"
	"time"
)

var _ safe.Dialer = (*tcp)(nil)

func NewTCP() *tcp {
	return &tcp{
		Rerun: safe.NewRerun(),
		Cache: cache.NewFile("user"),
	}
}

type tcp struct {
	Rerun safe.Rerun
	*client.Client
	Cache *cache.File
}

func (this *tcp) Update(address, username, password string) error {
	this.Cache.Set("address", address)
	this.Cache.Set("username", username)
	this.Cache.Set("password", password)
	this.Cache.Save()
	return this.Rerun.DialRun(this)
}

func (this *tcp) Dial(ctx context.Context) (err error) {
	//获取服务地址账号密码信息
	address := this.Cache.GetString("address")
	username := this.Cache.GetString("username")
	password := this.Cache.GetString("password")
	this.Client, err = client.DialWithContext(ctx,
		func(ctx context.Context) (ios.ReadWriteCloser, string, error) {
			c, err := net.DialTimeout("tcp", address, time.Second)
			return c, address, err
		}, func(c *client.Client) {
			c.Event.OnDealMessage = func(c *client.Client, msg ios.Acker) {
				bs := msg.Payload()
				data := new(output.Details)
				if err := json.Unmarshal(bs, data); err != nil {
					logs.Err(err)
					return
				}
				switch data.Param["type"] {
				case output.WinTypeVoice:
					err = notice.DefaultVoice.Speak(data.Content)
				case output.WinTypePopup:
					err = notice.DefaultWindows.Publish(&notice.Message{
						Target:  notice.TargetPopup,
						Title:   data.Title,
						Content: data.Content,
					})
				default:
					err = notice.DefaultWindows.Publish(&notice.Message{
						Title:   data.Title,
						Content: data.Content,
					})
				}
				logs.PrintErr(err)
			}
			t := time.Now()
			c.WriteAny(user.LoginReq{
				Username:  username,
				Signal:    oem.Signal(username, password, t),
				Timestamp: t.Unix(),
			})
		})
	logs.PrintErr(err)
	return
}

func (this *tcp) Run(ctx context.Context) error {
	this.Client.Ctx = ctx
	return this.Client.Run()
}

func (this *tcp) Close() error {
	if this.Client != nil {
		return this.Client.Close()
	}
	return nil
}

func (this *tcp) Closed() bool {
	return this.Client == nil || this.Client.Closed()
}
