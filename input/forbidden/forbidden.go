package forbidden

import (
	"fmt"
	"github.com/injoyai/conv/cfg"
	"strings"
)

var Forbidden *forbidden

func Init() {
	Forbidden = &forbidden{
		Words: cfg.GetStrings("input.forbidden"),
	}
}

type forbidden struct {
	Words []string
}

func (this *forbidden) Check(content string) error {
	for _, v := range this.Words {
		if strings.Contains(content, v) {
			return fmt.Errorf("包含违禁词:%s", v)
		}
	}
	return nil
}
