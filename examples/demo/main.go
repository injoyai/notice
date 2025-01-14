package main

import (
	"github.com/injoyai/goutil/oss"
	"github.com/injoyai/notice"
)

func main() {
	dir := oss.UserInjoyDir("notice/server/")
	dir = "./"
	notice.Default(dir)
}
