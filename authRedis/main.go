package authRedis

import (
	"github.com/go-redis/redis"
)

type Container struct {
	Client *redis.Client
}

func (c Container) Valid(l, p string) bool {
	val, err := c.Client.HMGet("stocking:auth:"+l, "password").Result()

	if err != nil {
		return false
	}
	
	if val[0].(string) != p {
		return false
	}

	return true
	
	
}
