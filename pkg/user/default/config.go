package user

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/injoyai/conv"
	"github.com/redis/go-redis/v9"
	"time"
)

// Signal 签名算法，可在这里自定义
func Signal(username, password string, timestamp time.Time) string {
	h := sha256.New()
	h.Write([]byte(username + password + conv.String(timestamp.Unix())))
	bs := h.Sum(nil)
	return hex.EncodeToString(bs)
}

type Config struct {
	Signal func(username, password string, timestamp time.Time) string //签名算法
	Type   string                                                      //数据库类型
	DSN    string                                                      //数据库配置
	Auth   AuthConfig                                                  //鉴权配置
}

type AuthConfig struct {
	Enable     bool           //启用
	Type       string         //类型,redis,memory,db 等方式
	Redis      *redis.Options //redis配置
	SuperToken []string       //超级token
}
