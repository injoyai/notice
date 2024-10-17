package oem

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/injoyai/base/maps"
	"time"
)

// Signal 签名算法，可在这里自定义
func Signal(username, password string, timestamp time.Time) string {
	h := sha256.New()
	h.Write([]byte(username + password + timestamp.String()))
	bs := h.Sum(nil)
	return hex.EncodeToString(bs)
}

var Token = maps.NewSafe()

func SetToken(token, username string, expiration time.Duration) error {
	Token.Set(token, username, expiration)
	return nil
}

func GetToken(token string) (string, error) {
	val, ok := Token.Get(token)
	if !ok {
		return "", nil
	}
	return val.(string), nil
}

// IsSuperToken 超级token，可以免登录
func IsSuperToken(token string) bool {
	switch token {
	case "123456789":
		return true
	}
	return false
}
