package client

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	ErrorReply  = '-'
	StatusReply = '+'
	IntReply    = ':'
	DoubleReply = ','
	NilReply    = '_'
	StringReply = '$'
	ArrayReply  = '*'
	MapReply    = '%'
	SetReply    = '~'
	BoolReply   = '#'
	VerbReply   = '='
	PushReply   = '>'
)

var (
	HugeEnuf float64 = 1e+300
	Infinity float64 = HugeEnuf * HugeEnuf
)

type RedisReply struct {
	Reply interface{}
	Type  rune
	Err   error
}

func (r *RedisReply) GetString() string {
	if r.Err != nil {
		return r.Err.Error()
	}
	if r.Type != StringReply && r.Type != ErrorReply &&
		r.Type != StatusReply && r.Type != VerbReply {
		return ``
	}

	return r.Reply.(string)
}

func (r *RedisReply) GetInt64() int64 {
	if r.Err != nil {
		return 0
	}
	if r.Type != IntReply {
		return 0
	}
	return r.Reply.(int64)
}

func (r *RedisReply) GetInt() int {
	if r.Err != nil {
		return 0
	}
	if r.Type != IntReply {
		return 0
	}
	return r.Reply.(int)
}

func (r *RedisReply) GetFloat64() float64 {
	if r.Err != nil {
		return 0.0
	}
	if r.Type != DoubleReply {
		return 0.0
	}
	return r.Reply.(float64)
}

func (r *RedisReply) GetFloat32() float32 {
	if r.Err != nil {
		return 0.0
	}
	if r.Type != DoubleReply {
		return 0.0
	}
	return r.Reply.(float32)
}

func (r *RedisReply) GetBool() bool {
	if r.Err != nil {
		return false
	}
	if r.Type != BoolReply {
		return false
	}
	return r.Reply.(bool)
}

func (r *RedisReply) GetArray() []RedisReply {
	if r.Err != nil {
		return nil
	}

	if r.Type != ArrayReply {
		return nil
	}

	return r.Reply.([]RedisReply)
}

// ***********************************************************

type commands struct {
	c        *SimpleClient
	replyBuf chan RedisReply
}

func (c *commands) init() {
	c.c.SetHandler(c.handler)
	c.replyBuf = make(chan RedisReply, 1)
}

func (c *commands) handler(data []byte) {
	// fmt.Printf("(test) recv:%v", string(data))
	if len(data) <= 0 {
		return
	}

	rep, _ := c.processItem(data)
	c.putReply(rep)
}

func (c *commands) putReply(r RedisReply) {
	defer func() {
		recover()
	}()

	select {
	case c.replyBuf <- r:
	default:
	}
}
func (c *commands) getReply() RedisReply {
	defer func() {
		recover()
	}()

	r, ok := <-c.replyBuf

	if !ok {
		return RedisReply{
			Err: fmt.Errorf(`chan close`),
		}
	}

	return r
}

// Do 单次来回
func (c *commands) Do(args ...string) RedisReply {
	err := c.c.Send([]byte(fmt.Sprintf("%s\r\n", strings.Join(args, ` `))))
	if err != nil {
		return RedisReply{
			Err: err,
		}
	}

	// todo 这里有可能会拿到别的消息, 除非服务端能保证不会主动下发消息
	return c.getReply()
}

// Doing 会阻塞, 持续读取服务器返回
func (c *commands) Doing(fn func(reply RedisReply), args ...string) error {
	if fn == nil {
		return errors.New(`Doing function is nil`)
	}

	err := c.c.Send([]byte(fmt.Sprintf("%s\r\n", strings.Join(args, ` `))))
	if err != nil {
		return err
	}

	for {
		select {
		case r, ok := <-c.replyBuf:
			if !ok {
				return nil
			}
			fn(r)
		}
	}
}

// 返回读取到的类型, 已经已读取的字节数
func (c *commands) processItem(data []byte) (RedisReply, int) {
	out := RedisReply{}
	readCount := 0

	out.Type = rune(data[0])
	switch data[0] {
	case ErrorReply:
		i := strings.Index(string(data), "\r\n")
		out.Err = errors.New(string(data[1:i]))
		readCount = i
	case StatusReply:
		i := strings.Index(string(data), "\r\n")
		out.Reply = string(data[1:i])
		readCount = i
	case StringReply:
		// todo 字符串不能这样直接赋值, 字符串还标出了长度
		fallthrough
	case VerbReply:
		/*
			例:
				$14
				123.123.123.123
		*/
		i := strings.Index(string(data), "\r\n")
		count, _ := strconv.Atoi(string(data[1:i]))
		if count < 0 {
			out.Type = NilReply
		} else {
			out.Reply = string(data[i+2 : i+2+count])
		}
		readCount = 2 + i + count
	case IntReply:
		i := strings.Index(string(data), "\r\n")
		v, e := strconv.ParseInt(string(data[1:i]), 10, 64)
		if e != nil {
			fmt.Printf("(debug) %v\n", e)
			out.Err = errors.New(`Bad integer value`)
		} else {
			out.Reply = v
		}
		readCount = i
	case DoubleReply:
		const inf = `,inf`
		const ninf = `,-inf`

		i := strings.Index(string(data), "\r\n")
		v := strings.ToLower(string(data[1:i]))
		if v == inf {
			out.Reply = Infinity
		} else if v == ninf {
			out.Reply = -Infinity
		} else {
			v, e := strconv.ParseFloat(v, 64)
			if e != nil {
				out.Err = e
			} else {
				out.Reply = v
			}
		}
		readCount = i
	case NilReply:
		readCount = 1
	case BoolReply:
		// c的redis中就是这样判断的
		if f := data[1:][0]; f == 't' || f == 'T' {
			out.Reply = true
		} else {
			out.Reply = false
		}
		readCount = strings.Index(string(data), "\r\n")
	case PushReply:
		i := strings.Index(string(data), "\r\n")
		fmt.Printf("(debug) Push:%v", string(data[1:i]))
		readCount = i
	case ArrayReply:
		i := strings.Index(string(data), "\r\n")
		readCount = i
		count, _ := strconv.Atoi(string(data[1:i]))
		if count <= 0 {
			out.Reply = nil
		} else {
			arry := make([]RedisReply, count)
			passCount := 0
			for t := 0; t < count; t++ {
				arry[t], passCount = c.processItem(data[readCount+2:])
				readCount += passCount
			}
			out.Reply = arry
		}
	case MapReply:
		i := strings.Index(string(data), "\r\n")
		fmt.Printf("(debug) Map:%v", string(data[1:i]))
		readCount = i
	case SetReply:
		i := strings.Index(string(data), "\r\n")
		fmt.Printf("(debug) Set:%v", string(data[1:i]))
		readCount = i
	default:
		out.Err = fmt.Errorf(`Protocol error, got %s as reply type byte`, data)
	}

	// out.Reply = data[1:]
	// 最后+2是因为结束符号为"\r\n"
	return out, readCount + 2
}
