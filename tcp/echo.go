package tcp

import (
	"bufio"
	"context"
	"go-redis/lib/logger"
	"go-redis/lib/sync/atomic"
	"go-redis/lib/sync/wait" //实现超时的功能
	"io"
	"net"
	"sync"
	_ "sync/atomic"
	"time"
)

type EchoClient struct {
	Conn    net.Conn
	Waiting wait.Wait
}

// MakeEchoHandler creates EchoHandler
func MakeHandler() *EchoHandler {
	return &EchoHandler{}
}

type EchoHandler struct {
	activeConn sync.Map
	closing    atomic.Boolean
}

//实现接口
func (h *EchoHandler) Handle(ctx context.Context, conn net.Conn) {
	if h.closing.Get() {
		// 关闭并拒绝新连接
		_ = conn.Close()
	}
	client := &EchoClient{
		Conn: conn,
	}
	h.activeConn.Store(client, struct{}{}) //hashset,只需key
	reader := bufio.NewReader(conn)
	//简单实现服务端回应用户端请求
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				logger.Info("connection close")
				h.activeConn.Delete(client)
			} else {
				logger.Warn(err)
			}
			return
		}
		client.Waiting.Add(1)
		b := []byte(msg)
		_, _ = conn.Write(b)
		client.Waiting.Done()
	}
}

func (c *EchoClient) Close() error {
	c.Waiting.WaitWithTimeout(10 * time.Second) //等待客户端业务做完，但如果超时也关闭
	_ = c.Conn.Close()
	return nil
}

// Close stops echo handler
func (h *EchoHandler) Close() error {
	logger.Info("handler shutting down...")
	h.closing.Set(true) //关闭业务引擎，拒绝新客户端连接
	h.activeConn.Range(func(key interface{}, val interface{}) bool {
		client := key.(*EchoClient)
		_ = client.Close()
		return true
	})
	return nil
}
