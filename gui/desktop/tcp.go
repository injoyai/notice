package main

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/injoyai/base/maps/wait/v2"
	"github.com/injoyai/base/safe"
	"github.com/injoyai/goutil/cache"
	"github.com/injoyai/goutil/notice"
	"github.com/injoyai/ios"
	"github.com/injoyai/ios/client"
	"github.com/injoyai/logs"
	"github.com/injoyai/notice/output"
	"github.com/injoyai/notice/user"
	"github.com/injoyai/notice/util"
	"net"
	"time"
)

var _ safe.Dialer = (*tcp)(nil)

func NewTCP() *tcp {
	t := &tcp{
		Rerun: util.NewRerun(),
		Cache: cache.NewFile("user"),
		wait:  wait.New(time.Second * 2),
	}
	return t
}

type tcp struct {
	Rerun *util.Rerun
	*client.Client
	Cache   *cache.File
	onLogin func()
	onClose func(err error)
	wait    *wait.Entity
}

func (this *tcp) Update(address, username, password string) error {
	this.Cache.Set("address", address)
	this.Cache.Set("username", username)
	this.Cache.Set("password", password)
	this.Cache.Save()
	if this.Client != nil {
		this.Client.Close()
	}
	return this.Rerun.DialRun(this)
}

func (this *tcp) Dial(ctx context.Context) (err error) {
	//获取服务地址账号密码信息
	address := this.Cache.GetString("address")
	username := this.Cache.GetString("username")
	password := this.Cache.GetString("password")
	this.Client, err = client.DialWithContext(
		ctx,
		this.dial(address),
		func(c *client.Client) {
			c.Event.OnDealMessage = this.dealMessage
			c.Event.OnDisconnect = func(c *client.Client, err error) {
				if this.onClose != nil {
					this.onClose(err)
				}
			}
			t := time.Now()
			c.WriteAny(user.LoginReq{
				ID:        t.String(),
				Username:  username,
				Signal:    user.Signal(username, password, t),
				Timestamp: t.Unix(),
			})
			if _, err := this.wait.Wait(t.String()); err != nil {
				logs.Err(err)
				return
			}
			if this.onLogin != nil {
				this.onLogin()
			}

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

func (this *tcp) dial(address string) func(ctx context.Context) (ios.ReadWriteCloser, string, error) {
	return func(ctx context.Context) (ios.ReadWriteCloser, string, error) {
		c, err := net.DialTimeout("tcp", address, time.Second)
		return c, address, err
	}
}

func (this *tcp) dealMessage(c *client.Client, msg ios.Acker) {
	bs := msg.Payload()

	resp := new(output.Resp)
	if err := json.Unmarshal(bs, resp); err == nil && len(resp.ID) > 0 {
		if resp.Success {
			this.wait.Done(resp.ID, nil)
		} else {
			this.wait.Done(resp.ID, errors.New(resp.Message))
		}
		return
	}

	data := new(output.Details)
	if err := json.Unmarshal(bs, data); err != nil {
		logs.Err(err)
		return
	}
	var err error
	switch data.Type {
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
