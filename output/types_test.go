package output

import "testing"

func TestMessage_Check(t *testing.T) {
	m := &Message{
		Output: []string{"wechat:group:1"},
	}
	t.Log(m.Check(map[string]struct{}{}))
}
