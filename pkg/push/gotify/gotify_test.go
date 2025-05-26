package gotify

import (
	"github.com/injoyai/conv/cfg"
	"github.com/injoyai/notice/pkg/push"
	"testing"
)

func TestNew(t *testing.T) {
	cfg.Init(cfg.WithFile("../../../config/config_real.yaml"))
	host := cfg.GetString("gotify.host")
	t.Log(host)
	token := cfg.GetString("gotify.token")
	t.Log(token)
	priority := cfg.GetInt("gotify.priority")
	t.Log(priority)
	x := New(host, token, priority)
	err := x.Push(&push.Message{Title: "title", Content: "content"})
	t.Log(err)
}
