package main

type config struct {
	hostIP     string
	hostPort   string
	hostSocket string
}

func newConfig() *config {
	return &config{
		hostIP:     `127.0.0.1`,
		hostPort:   `6379`,
		hostSocket: ``,
	}
}

// Version 返回当前版本号
func Version() string {
	return `0.0.1`
}

// RedisVersion 返回支持的最新的redis版本号
func RedisVersion() string {
	return `6.0.10`
}
