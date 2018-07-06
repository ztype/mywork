package router

import (
	"golang.org/x/net/websocket"
	"github.com/segmentio/ksuid"
	"log"
	"sync"
)

type Connect struct {
	ws    *websocket.Conn
	id    string
	chin  chan []byte //channel to receive ws msg
	chout chan []byte //channel to send ws msg
}

var (
	connMap map[string]*Connect
	clock sync.Mutex
)

func ServeConnect(c *Connect){
	clock.Lock()
	connMap[c.ID()] = c
	clock.Unlock()
	c.active()
}

func NewConnect(ws *websocket.Conn) *Connect {
	return &Connect{
		ws:    ws,
		id:    ksuid.New().String(),
		chin:  make(chan []byte, 0),
		chout: make(chan []byte, 0),
	}
}

func (c *Connect) ID() string {
	return c.id
}

func (c *Connect) Close(){
	c.ws.Close()
	close(c.chin)
	close(c.chout) // fixme
}

func (c *Connect) active() {
	go c.wsRecv()
	c.wsSend()
}

func (c *Connect) wsRecv() {
	for {
		buf := make([]byte, 0)
		err := websocket.Message.Receive(c.ws, &buf)
		if err != nil {
			log.Println(c.ws.Request().RemoteAddr, err)
			break
		}
		select {
		case c.chin <- buf:
		default:
			log.Println(c.id, "ws chan in blocked")
		}
	}
	log.Println(c.id, "ws recv quit")
}

func (c *Connect) wsSend() {
	for {
		select {
		case buf, ok := <-c.chout:
			if !ok {
				log.Println(c.id, "ws closed")
				break
			}
			err := websocket.Message.Send(c.ws, buf)
			if err != nil {
				log.Println(c.ws.RemoteAddr(), err)
				break
			}
		}
	}
	log.Println(c.id, "ws send quit")
}
