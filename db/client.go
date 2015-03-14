package db

import (
	"time"

	"github.com/fzzy/radix/redis"
)

type Client struct {
	redis *redis.Client
}

func New(url string) *Client {
	c, _ := redis.DialTimeout("tcp", url, time.Duration(10)*time.Second)
	return &Client{redis: c}
}

func (c *Client) Set(key, value string) error {
	c.redis.Cmd("select", 0)
	r := c.redis.Cmd("set", key, value)
	return r.Err
}

func (c *Client) Close() {
	c.redis.Close()
}

func (c *Client) Get(key string) string {
	str, _ := c.redis.Cmd("get", key).Str()
	return str
}
