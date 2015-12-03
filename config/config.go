package config

import "github.com/cghsystems/godata/env"

func RedisUrl() string {
	if value, err := env.Get("godata_redis_url", "127.0.0.1:6379"); err != nil {
		panic("An unexpected error has occurred")
	} else {
		return value
	}
}
