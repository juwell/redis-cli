package client

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	// Cli redis连接
	Cli *cli
)

func init() {
	Cli = newClient()
}

type cli struct {
	*redis.Client
}

func newClient() *cli {
	return &cli{}
}

func (c *cli) Connect() error {
	opt := redis.Options{}
	if len(Cfg.HostSocket) <= 0 {
		opt.Addr = Cfg.HostIP + `:` + Cfg.HostPort
	} else {
		opt.Addr = Cfg.HostSocket
	}
	c.Client = redis.NewClient(&opt)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	state := c.Client.Ping(ctx)
	if state == nil {
		return fmt.Errorf(`connect to %s err`, opt.Addr)
	} else if state.Err() != nil {
		return state.Err()
	}
	return nil
}

// Redirection 重定向
func (c *cli) Redirection() error {
	op := c.Client.Options()

	newAddr := ``
	if len(Cfg.HostSocket) <= 0 {
		newAddr = Cfg.HostIP + `:` + Cfg.HostPort
	} else {
		newAddr = Cfg.HostSocket
	}

	if op != nil && op.Addr != newAddr {
		c.Client.Close()

		return c.Connect()
	}

	return nil
}
