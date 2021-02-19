package client

import (
	"bufio"
	"context"
	"fmt"
	"net"

	"github.com/go-redis/redis/v8"
)

type Buffer struct {
	Buff []byte
}

type MonitorClient struct {
	conn   net.Conn
	option redis.Options

	buff chan Buffer
	exit chan bool
}

func NewMonitorClient(opt redis.Options) *MonitorClient {
	return &MonitorClient{
		option: opt,
		buff:   make(chan Buffer, 1000),
		exit:   make(chan bool),
	}
}

func (m *MonitorClient) Connect() error {
	conn, err := m.option.Dialer(context.Background(), m.option.Network, m.option.Addr)
	if err != nil {
		return err
	}

	m.conn = conn
	return nil
}

func (m *MonitorClient) Do(fn func(d []byte)) {
	if fn == nil {
		return
	}

	go m.readGoroutine()

	// 发送monitor命令
	// w := bufio.NewWriter(m.conn)
	// _, err := w.WriteString("monitor\r\n")
	_, err := m.conn.Write([]byte("monitor\r\n"))
	if err != nil {
		fn([]byte(err.Error()))
		defer recover()
		close(m.exit)
		close(m.buff)
		return
	}

	for {
		select {
		case _, ok := <-m.exit:
			if !ok {
				return
			}
		case data, _ := <-m.buff:
			// fmt.Println(`[test] here`)
			fn(data.Buff)
		}
	}
}

func (m *MonitorClient) Close() {
	defer recover()

	if m.conn != nil {
		m.conn.Close()
		m.conn = nil
	}

	close(m.exit)
	close(m.buff)
}

func (m *MonitorClient) readGoroutine() {
	// data := make([]byte, 1024*16)
	r := bufio.NewReader(m.conn)
	for {
		select {
		case _, ok := <-m.exit:
			if !ok {
				return
			}
		default:
		}
		data, _, err := r.ReadLine()
		// fmt.Println(`[test]`, string(data))
		// _, err := r.Read(data)
		if err != nil {
			fmt.Println(err)
			defer recover()
			close(m.exit)
			return
		}

		m.buff <- Buffer{
			Buff: data,
		}
	}
}
