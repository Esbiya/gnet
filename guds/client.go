package guds

import (
	"github.com/Esbiya/loguru"
	"io"
	"net"
	"os"
	"sync"
	"time"
)

type (
	Callback struct {
		Err  error
		Done bool
		Body chan interface{}
	}
	Client struct {
		wg             *sync.WaitGroup
		addr           *net.UnixAddr
		conn           *net.UnixConn
		retMap         *sync.Map
		connectTimeout time.Duration
	}
)

func (c *Callback) Then(callback func(b interface{})) *Callback {
	callback(<-c.Body)
	return c
}

func (c *Callback) Catch(callback func(err error)) *Callback {
	callback(c.Err)
	return c
}

func (c *Callback) Close() {
	close(c.Body)
}

func NewClient(address string, connectTimeout time.Duration) (*Client, error) {
	var err error
	c := Client{wg: &sync.WaitGroup{}, retMap: &sync.Map{}, connectTimeout: connectTimeout}
	c.addr, err = net.ResolveUnixAddr("unix", address)
	if err != nil {
		return nil, err
	}
	c.conn, err = net.DialUnix("unix", nil, c.addr)
	if err != nil {
		return nil, err
	}
	c.retMap.Store("client.init", make(chan interface{}))
	return &c, nil
}

func DefaultClient() *Client {
	c, _ := NewClient("/tmp/us.socket", 5*time.Second)
	return c
}

func (c *Client) loopRead() {
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		for {
			// 读取消息长度
			s := make([]byte, 4)
			_, err := c.conn.Read(s)
			if err == io.EOF {
				return
			}
			if err != nil {
				continue
			}
			size := BytesToInt(s)

			// 读取消息主体
			d := make([]byte, size)
			l, err := c.conn.Read(d)
			if err != nil {
				loguru.Error("read data error: %v", err)
				continue
			}

			msg := Message{}
			err = msg.Parse(d[:l])
			if err != nil {
				loguru.Error(err)
			}

			if v, ok := c.retMap.Load(msg.Api); ok {
				v.(chan interface{}) <- msg.Data
			}
		}
	}()
}

func (c *Client) read(done chan struct{}) <-chan []byte {
	context := make(chan []byte)
	c.wg.Add(1)
	go func() {
		defer close(context)
		for {
			// 读取消息长度
			s := make([]byte, 4)
			_, err := c.conn.Read(s)
			if err == io.EOF {
				return
			}
			if err != nil {
				time.Sleep(1 * time.Second)
				continue
			}
			size := BytesToInt(s)

			// 读取消息主体
			d := make([]byte, size)
			l, err := c.conn.Read(d)
			if err != nil {
				loguru.Error("read data error: %v", err)
				continue
			}
			select {
			case <-done:
				return
			case context <- d[:l]:
			}
		}
	}()
	return context
}

func (c *Client) Send(api string, context interface{}) *Callback {
	message := NewMessage(api, context)
	_, err := c.conn.Write(MergeBytes(IntToBytes(message.length), message.bytes))
	if err != nil {
		return &Callback{
			Err: err,
		}
	}
	out := make(chan interface{})
	c.retMap.Store(api, out)
	return &Callback{
		Body: out,
	}
}

func (c *Client) Start() {
	c.loopRead()
	ch, _ := c.retMap.Load("client.init")
	select {
	case <-(ch.(chan interface{})):
		loguru.Info("connect success! ")
	case <-time.After(c.connectTimeout):
		loguru.Error("connect timeout! ")
		os.Exit(1)
	}
}

func (c *Client) Wait() {
	c.wg.Wait()
}
