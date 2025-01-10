package main

import (
	"github.com/injoyai/goutil/oss"
	"github.com/injoyai/notice"
)

func main() {
	notice.Default(oss.UserInjoyDir("notice/server/"))
}
