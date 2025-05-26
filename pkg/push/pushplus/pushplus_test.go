package pushplus

import (
	"github.com/injoyai/conv/cfg"
	"github.com/injoyai/notice/pkg/push"
	"testing"
)

func TestNew(t *testing.T) {
	cfg.Init(cfg.WithFile("../../../config/config_real.yaml"))
	token := cfg.GetString("pushPlus.token")
	t.Log(token)
	x := New(token)
	err := x.Push(&push.Message{Title: "title", Content: "content"})
	t.Log(err)
}
