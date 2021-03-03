package test

import (
	"encoding/json"
	"github.com/Esbiya/gnet/guds"
	"github.com/Esbiya/loguru"
	"io"
	"net"
	"os"
	"sync"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	addr, err := net.ResolveUnixAddr("unix", "/tmp/us.socket")
	if err != nil {
		panic("cannot resolve unix addr: " + err.Error())
	}
	// 拔号
	c, err := net.DialUnix("unix", nil, addr)
	if err != nil {
		panic("DialUnix failed: " + err.Error())
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	// 读结果
	go func() {
		defer wg.Done()
		for {
			// 读取消息长度
			s := make([]byte, 4)
			_, err := c.Read(s)
			if err == io.EOF {
				return
			}
			if err != nil {
				time.Sleep(1 * time.Second)
				continue
			}
			size := guds.BytesToInt(s)

			// 读取消息主体
			d := make([]byte, size)
			l, err := c.Read(d)
			if err != nil {
				loguru.Error("read data error: %v", err)
				os.Exit(1)
			}
			loguru.Debug(string(d[:l]))
		}
	}()

	// 写入
	data, _ := json.Marshal(map[string]interface{}{
		"api": "session.login",
	})
	_, err = c.Write(guds.MergeBytes(guds.IntToBytes(len(data)), data))
	if err != nil {
		panic("Writes failed.")
	}

	wg.Wait()
}
