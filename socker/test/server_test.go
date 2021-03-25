package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/Esbiya/gnet/socker"
)

func testLogin(qr chan string) bool {
	<-time.After(1 * time.Second)
	qr <- "xxx"
	<-time.After(1 * time.Second)
	return true
}

func TestServer(t *testing.T) {
	server := socker.DefaultTCPServer()

	server.Router().Register("session.login", func(msg socker.Data, c chan socker.Reply) {
		qr := make(chan string, 1)
		go func() {
			c <- socker.Reply{
				Status: socker.Continue,
				Body:   fmt.Sprintf(`{"code":200,"msg":"success","data":{"qr":"%s"}}`, <-qr),
			}
		}()
		if testLogin(qr) {
			c <- socker.Reply{
				Status: socker.Done,
				Body: map[string]interface{}{
					"code": 200,
					"msg":  "success",
					"data": map[string]interface{}{
						"session": "哈哈哈",
					},
				},
			}
		}
	})

	server.Run()
}
