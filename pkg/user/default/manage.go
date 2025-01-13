package user

import (
	"github.com/injoyai/base/maps"
	"github.com/redis/go-redis/v9"
	"time"
)

const (
	Redis  = "redis"
	Memory = "memory"
)

func NewManage(cfg *Config) *Manage {
	return &Manage{
		Config: cfg,
		Cache: func() Cache2 {
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
		Signal: nil,
	}
}

type Manage struct {
	*Config
	Cache       Cache2
	Signal      func(username, password string, timestamp time.Time) string
	CheckToken  bool
	SuperTokens []string
}

// IsSuper 超级token，可以免校验
func (this *Manage) IsSuper(token string) bool {
	if !this.CheckToken {
		return true
	}
	for _, v := range this.SuperTokens {
		if token == v {
			return true
		}
	}
	return false
}
