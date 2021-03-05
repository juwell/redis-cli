package client

import (
	"bufio"
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"runtime"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

type Buffer struct {
	Buff []byte
}

// SimpleClient 重写了redis连接客户端, 名字没有想好, 先随便用这个了
type SimpleClient struct {
	commands
	conn   net.Conn
	option redis.Options

	readBuff  chan Buffer
	writeBuff chan Buffer
	exit      chan bool
	handler   func(d []byte)
}

func NewSimpleClient(opt redis.Options) *SimpleClient {
	if opt.Addr == "" {
		opt.Addr = "localhost:6379"
	}
	if opt.Network == "" {
		if strings.HasPrefix(opt.Addr, "/") {
			opt.Network = "unix"
		} else {
			opt.Network = "tcp"
		}
	}
	if opt.DialTimeout == 0 {
		opt.DialTimeout = 5 * time.Second
	}
	if opt.Dialer == nil {
		opt.Dialer = func(ctx context.Context, network, addr string) (net.Conn, error) {
			netDialer := &net.Dialer{
				Timeout:   opt.DialTimeout,
				KeepAlive: 5 * time.Minute,
			}
			if opt.TLSConfig == nil {
				return netDialer.DialContext(ctx, network, addr)
			}
			return tls.DialWithDialer(netDialer, network, addr, opt.TLSConfig)
		}
	}
	if opt.PoolSize == 0 {
		opt.PoolSize = 10 * runtime.NumCPU()
	}
	switch opt.ReadTimeout {
	case -1:
		opt.ReadTimeout = 0
	case 0:
		opt.ReadTimeout = 3 * time.Second
	}
	switch opt.WriteTimeout {
	case -1:
		opt.WriteTimeout = 0
	case 0:
		opt.WriteTimeout = opt.ReadTimeout
	}
	if opt.PoolTimeout == 0 {
		opt.PoolTimeout = opt.ReadTimeout + time.Second
	}
	if opt.IdleTimeout == 0 {
		opt.IdleTimeout = 5 * time.Minute
	}
	if opt.IdleCheckFrequency == 0 {
		opt.IdleCheckFrequency = time.Minute
	}

	if opt.MaxRetries == -1 {
		opt.MaxRetries = 0
	} else if opt.MaxRetries == 0 {
		opt.MaxRetries = 3
	}
	switch opt.MinRetryBackoff {
	case -1:
		opt.MinRetryBackoff = 0
	case 0:
		opt.MinRetryBackoff = 8 * time.Millisecond
	}
	switch opt.MaxRetryBackoff {
	case -1:
		opt.MaxRetryBackoff = 0
	case 0:
		opt.MaxRetryBackoff = 512 * time.Millisecond
	}

	c := &SimpleClient{
		option:    opt,
		readBuff:  make(chan Buffer, 100000),
		writeBuff: make(chan Buffer, 100000),
		exit:      make(chan bool),
	}
	c.commands.c = c
	c.init()
	return c
}

func (m *SimpleClient) SetHandler(fn func(d []byte)) {
	m.handler = fn
}

func (m *SimpleClient) Connect() error {
	conn, err := m.option.Dialer(context.Background(), m.option.Network, m.option.Addr)
	if err != nil {
		return err
	}

	m.conn = conn

	go m.readGoroutine()
	go m.writeGoroutine()
	go m.handleGoroutine()

	return nil
}

func (m *SimpleClient) Close() {
	defer func() {
		recover()
	}()

	if m.conn != nil {
		m.conn.Close()
		m.conn = nil
	}

	close(m.exit)
	close(m.readBuff)
	close(m.writeBuff)
}

func (m *SimpleClient) readGoroutine() {
	defer func() {
		recover()
	}()
	data := make([]byte, 1024*16)
	var total []byte
	r := bufio.NewReader(m.conn)
	for {
		select {
		case _, ok := <-m.exit:
			if !ok {
				return
			}
		default:
		}
		n, err := r.Read(data)
		// log.Println("[debug]", string(data[:n]))
		// fmt.Println(`[debug]`, string(data))
		// _, err := r.Read(data)
		if err != nil {
			// fmt.Println("(debug)", err)
			close(m.exit)
			return
		}

		if string(data[n-2:n]) == "\r\n" {
			if total == nil {
				// log.Println("[debug1]", string(data[:n]))
				m.readBuff <- Buffer{
					Buff: data[:n],
				}
			} else {
				// 合并
				total = append(total, data[:n]...)
				// log.Println("[debug2]", string(total))
				m.readBuff <- Buffer{
					Buff: total,
				}
				total = nil
			}
		} else {
			// 等下一个包
			if total == nil {
				total = make([]byte, 0, n)
				total = append(total, data[:n]...)
			} else {
				// 合并
				total = append(total, data[:n]...)
			}
		}
	}
}

func (m *SimpleClient) writeGoroutine() {
	for {
		select {
		case _, ok := <-m.exit:
			if !ok {
				return
			}
		case b := <-m.writeBuff:
			if len(b.Buff) > 0 {
				m.conn.Write(b.Buff)
			}
		}
	}
}

func (m *SimpleClient) send(data []byte) error {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("(debug)", err)
		}
	}()

	select {
	case _, ok := <-m.exit:
		if !ok {
			return fmt.Errorf(`Connection is closed`)
		}
	default:
	}

	m.writeBuff <- Buffer{
		Buff: data,
	}
	return nil
}

func (m *SimpleClient) handleGoroutine() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("(debug)", err)
		}
	}()

	for {
		select {
		case _, ok := <-m.exit:
			if !ok {
				return
			}
		case data, _ := <-m.readBuff:
			// fmt.Println(`[test] here`)
			m.handler(data.Buff)
		}
	}
}

// func (m *SimpleClient) Redirection(opt redis.Options) error {
// 	if m.option.Addr != opt.Addr || m.option.Network != opt.Network {
// 		m.option = opt

// 		m.Close()
// 		return m.Connect()
// 	}
// 	return nil
// }
