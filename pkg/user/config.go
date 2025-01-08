package user

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/injoyai/base/maps"
	"github.com/injoyai/conv"
	"github.com/injoyai/conv/cfg"
	"github.com/injoyai/goutil/database/redis"
	"time"
)

// Signal 签名算法，可在这里自定义
func Signal(username, password string, timestamp time.Time) string {
	h := sha256.New()
	h.Write([]byte(username + password + conv.String(timestamp.Unix())))
	bs := h.Sum(nil)
	return hex.EncodeToString(bs)
}

var Token *tokenCache

func initToken() {
	Token = &tokenCache{
		Type:   cfg.GetString("token.type"),
		Memory: maps.NewSafe(),
		Redis: redis.New(
			cfg.GetString("redis.address"),
			cfg.GetString("redis.password"),
			cfg.GetInt("redis.db"),
		),
		Super:   cfg.GetStrings("token.super"),
		Disable: cfg.GetBool("token.disable"),
	}
}

type tokenCache struct {
	Type    string
	Memory  *maps.Safe
	Redis   *redis.Client
	Super   []string
	Disable bool
}

func (this *tokenCache) Get(token string) (string, error) {
	switch this.Type {
	case "redis":
		username, err := this.Redis.Get(token)
		if err != nil {
			if err.Error() == redis.Nil.Error() {
				return "", nil
			}
			return "", err
		}
		return username, nil
	default:
		val, ok := this.Memory.Get(token)
		if !ok {
			return "", nil
		}
		return val.(string), nil
	}

}

func (this *tokenCache) Set(token, username string, expiration time.Duration) error {
	switch this.Type {
	case "redis":
		return this.Redis.Set(token, username, expiration)
	default:
		this.Memory.Set(token, username, expiration)
		return nil
	}
}

// IsSuper 超级token，可以免校验
func (this *tokenCache) IsSuper(token string) bool {
	if this.Disable {
		return true
	}
	for _, v := range this.Super {
		if token == v {
			return true
		}
	}
	return false
}
