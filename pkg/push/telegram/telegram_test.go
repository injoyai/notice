package telegram

import (
	"github.com/injoyai/conv/cfg"
	"github.com/injoyai/notice/pkg/push"
	"testing"
)

func TestNew(t *testing.T) {
	cfg.Init(cfg.WithFile("../../../config/config.yaml"))
	token := cfg.GetString("telegram.token")
	proxy := cfg.GetString("telegram.proxy")
	chatID := cfg.GetString("telegram.chatID")

	te, err := New(token, proxy, chatID)
	if err != nil {
		t.Error(err)
		return
	}

	err = te.Push(&push.Message{Content: "content"})
	if err != nil {
		t.Error(err)
		return
	}
}
