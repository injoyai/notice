package plugin

import "github.com/injoyai/notice/output"

func New() *Plugin {
	return &Plugin{}
}

type Plugin struct{}

func (*Plugin) Types() []string {
	return []string{output.TypePlugin}
}

func (*Plugin) Push(msg *output.Message) error { return nil }
