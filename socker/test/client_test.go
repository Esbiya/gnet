package test

import (
	"testing"

	"github.com/Esbiya/gnet/socker"
	"github.com/Esbiya/loguru"
)

func TestClient(t *testing.T) {
	c := socker.DefaultTCPClient()
	c.Start()
	c.Send("session.login", map[string]interface{}{
		"xxx": "111",
	}).Then(func(b interface{}) {
		loguru.Debug(b)
	}).Then(func(b interface{}) {
		loguru.Debug(b)
	}).Close()
	c.Wait()
}
