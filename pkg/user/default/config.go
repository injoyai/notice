package user

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/injoyai/base/maps"
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
	Type string     //数据库类型
	DSN  string     //数据库配置
	Auth AuthConfig //鉴权配置
}

type AuthConfig struct {
	Enable     bool           //启用
	Type       string         //类型,redis,memory,db 等方式
	Redis      *redis.Options //redis配置
	SuperToken []string       //超级token
}

func initToken(cfg *Config) *auth {
	if cfg == nil {
		cfg = &Config{}
	}
	m := &auth{
		Enable: cfg.Auth.Enable,
		Cache: func() Cache {
			switch cfg.Auth.Type {
			case Redis:
				return &_redis{
					Client: redis.NewClient(cfg.Auth.Redis),
				}
			default:
				return &_memory{
					Safe: maps.NewSafe(),
				}
			}
		}(),
		SuperTokens: cfg.Auth.SuperToken,
	}
	return m
}

type auth struct {
	Enable      bool
	Cache       Cache
	SuperTokens []string
}

// IsSuper 超级token，可以免校验
func (this *auth) IsSuper(token string) bool {
	if !this.Enable {
		return true
	}
	for _, v := range this.SuperTokens {
		if token == v {
			return true
		}
	}
	return false
}
