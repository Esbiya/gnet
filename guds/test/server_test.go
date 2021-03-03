package test

import (
	"fmt"
	"github.com/Esbiya/gnet/guds"
	"testing"
	"time"
)

func testLogin(qr chan string) bool {
	<-time.After(1 * time.Second)
	qr <- "xxx"
	<-time.After(1 * time.Second)
	return true
}

func TestServer(t *testing.T) {
	server := guds.DefaultServer()

	server.Router().Register("session.login", func(msg guds.Data, c chan guds.Reply) {
		qr := make(chan string, 1)
		go func() {
			c <- guds.Reply{
				Status: guds.Continue,
				Body:   fmt.Sprintf(`{"code":200,"msg":"success","data":{"qr":"%s"}}`, <-qr),
			}
		}()
		if testLogin(qr) {
			c <- guds.Reply{
				Status: guds.Done,
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
