package forbidden

import (
	"fmt"
	"strings"
)

func New(words ...string) *Forbidden {
	return &Forbidden{
		Words: words,
	}
}

type Forbidden struct {
	Words []string
}

func (this *Forbidden) Check(title, content string) error {
	for _, v := range this.Words {
		if strings.Contains(title, v) {
			return fmt.Errorf("包含违禁词:%s", v)
		}
		if strings.Contains(content, v) {
			return fmt.Errorf("包含违禁词:%s", v)
		}
	}
	return nil
}
