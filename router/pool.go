package router

import (
	"golang.org/x/net/websocket"
	"sync"
)

var (
	pool  []*websocket.Conn
	plock sync.Mutex
)

func WsIn(conn *websocket.Conn) {
	plock.Lock()
	defer plock.Unlock()
	pool = append(pool, conn)
}


