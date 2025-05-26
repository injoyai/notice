package serverchan

import (
	"github.com/injoyai/conv/cfg"
	"github.com/injoyai/notice/pkg/push"
	"testing"
)

func TestNew(t *testing.T) {
	cfg.Init(cfg.WithFile("../../../config/config.yaml"))
	sendKey := cfg.GetString("serverChan.sendKey")
	t.Log(sendKey)
	x := New(sendKey)
	err := x.Push(&push.Message{Title: "title", Content: "content"})
	t.Log(err)
}
