package main

import (
	"github.com/go-redis/redis/v8"
)

type client struct {
	rdb *redis.Client
}

func newClient() *client {
	return &client{}
}

func (c *client) Connect() {
	opt := redis.Options{}
	if len(cfg.hostSocket) <= 0 {
		opt.Addr = cfg.hostIP + `:` + cfg.hostPort
	} else {
		opt.Addr = cfg.hostSocket
	}
	c.rdb = redis.NewClient(&opt)
}
