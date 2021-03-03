package test

import (
	"github.com/Esbiya/gnet/guds"
	"github.com/Esbiya/loguru"
	"testing"
)

func TestClient(t *testing.T) {
	c := guds.DefaultClient()
	c.Start()
	c.Send("session.login", map[string]interface{}{}).Then(func(b interface{}) {
		loguru.Debug(b)
	}).Then(func(b interface{}) {
		loguru.Debug(b)
	}).Close()
	c.Wait()
}
