package test

import (
	"github.com/panjf2000/gnet/guds"
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
	server := guds.Default()

	server.Router().Register("session.login", func(msg guds.Data, c chan guds.Reply) {
		qr := make(chan string, 1)
		go func() {
			c <- guds.Reply{
				Status: guds.Continue,
				Body: map[string]interface{}{
					"code": 200,
					"msg":  "success",
					"data": map[string]interface{}{
						"qr": <-qr,
					},
				},
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
