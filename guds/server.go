package guds

import (
	"fmt"
	"github.com/Esbiya/gnet"
	"github.com/Esbiya/loguru"
	"github.com/panjf2000/ants/v2"
	"os"
	"sync/atomic"
	"time"
)

type (
	Conn   gnet.Conn
	Action gnet.Action
	Server struct {
		*gnet.EventServer
		addr         string
		router       *Router
		pool         *ants.Pool
		clients      int
		started      int32
		connected    int32
		disconnected int32
		timeout      time.Duration
		multicore    bool
	}
)

const (
	None     Action = iota
	Close           // 连接关闭
	ShutDown        // 服务关闭
	Continue        // 消息继续
	Done            // 消息结束
)

const (
	DefaultAntsPoolSize = 1 << 18

	ExpiryDuration = 10 * time.Second

	Nonblocking = true
)

func (u *Server) Router() *Router {
	return u.router
}

func (u *Server) OnInitComplete(srv gnet.Server) (action gnet.Action) {
	loguru.Info("server is listening on %s (multi-cores: %t, loops: %d)", srv.Addr.String(), srv.Multicore, srv.NumEventLoop)
	return
}

func (u *Server) OnOpened(c gnet.Conn) (out []byte, action gnet.Action) {
	c.SetContext(c)
	atomic.AddInt32(&u.connected, 1)
	msg := NewMessage("client.init", "hello world")
	out = MergeBytes(IntToBytes(msg.length), msg.bytes)
	if c.LocalAddr() == nil {
		panic("nil local addr")
	}
	if c.RemoteAddr() == nil {
		panic("nil local addr")
	}
	return
}

func (u *Server) OnClosed(c gnet.Conn, err error) (action gnet.Action) {
	if c.Context() != c {
		panic("invalid context")
	}

	atomic.AddInt32(&u.disconnected, 1)
	if atomic.LoadInt32(&u.connected) == atomic.LoadInt32(&u.disconnected) &&
		atomic.LoadInt32(&u.disconnected) == int32(u.clients) {
		action = gnet.Shutdown
	}

	return
}

func NewServer(addr string, multicore bool, poolSize int, timeout time.Duration) *Server {
	options := ants.Options{ExpiryDuration: ExpiryDuration, Nonblocking: Nonblocking}
	pool, _ := ants.NewPool(poolSize, ants.WithOptions(options))
	return &Server{
		addr: addr,
		router: &Router{
			Callbacks: map[string]func(msg Data, c chan Reply){},
		},
		pool:      pool,
		timeout:   timeout,
		multicore: multicore,
	}
}

// 记录日志
func (u *Server) RecordLog() {
	err := loguru.Enable(loguru.FileLog)
	if err != nil {
		loguru.Error("log write file error: %v", err)
		os.Exit(1)
	}
}

func (u *Server) React(frame []byte, c gnet.Conn, out chan gnet.Out) {
	var msg Message
	err := msg.Parse(frame[4:])

	reply := make(chan Reply)
	if err != nil {
		msg.reset(false, map[string]interface{}{
			"code": 400,
			"msg":  "message decode err! ",
		})
		out <- gnet.Out{
			Body: msg.out(),
		}
	} else {
		loguru.Debug("receive message - length: %d, body: %s", msg.length, msg.Stringify())
		_ = u.pool.Submit(func() {
			u.router.Get(msg.Api)(msg.ToData(), reply)
		})
	LOOP:
		select {
		case _reply := <-reply:
			msg.reset(_reply.Async, _reply.Body)
			loguru.Debug("reply   message - length: %d, body: %s", msg.length, msg.Stringify())
			if msg.async {
				_ = c.AsyncWrite(msg.out())
			} else {
				out <- gnet.Out{
					Body:   msg.out(),
					Action: gnet.Action(_reply.Status),
				}
			}
			switch _reply.Status {
			case Continue:
				goto LOOP
			default:
			}
		case <-time.After(u.timeout):
			msg.reset(false, map[string]interface{}{
				"code": 500,
				"msg":  "process message timeout! ",
			})
			out <- gnet.Out{
				Body: msg.out(),
			}
		}
	}
	return
}

func (u *Server) version() {
	fmt.Println(loguru.Fuchsia("                    __"))
	fmt.Println(loguru.Fuchsia("   ____ ___  ______/ /____"))
	fmt.Println(loguru.Fuchsia("  / __ `/ / / / __  / ___/"))
	fmt.Println(loguru.Fuchsia(" / /_/ / /_/ / /_/ (__  )"))
	fmt.Println(loguru.Fuchsia(" \\__, /\\__,_/\\__,_/____/"))
	fmt.Println(loguru.Fuchsia("/____/                      v1.0.0"))
}

func (u *Server) Run() {
	u.version()
	_ = gnet.Serve(u, fmt.Sprintf("unix://%s", u.addr), gnet.WithMulticore(u.multicore))
}

func DefaultServer() *Server {
	return NewServer("/tmp/us.socket", true, DefaultAntsPoolSize, ExpiryDuration)
}
