package router

import (
	"fmt"
	"log"

	"github.com/segmentio/ksuid"
	"golang.org/x/net/websocket"
	"bytes"
)

type OnMessage func(string,[]byte)
type OnClose func(id string)

var SignalClose = []byte{0x00}

var ErrBlock = fmt.Errorf("chan blocked")

type Connect struct {
	ws      *websocket.Conn
	id      string
	chin    chan []byte //channel to receive ws msg
	chout   chan []byte //channel to send ws msg
	onclose []OnClose
}

func NewConnect(ws *websocket.Conn) *Connect {
	c := &Connect{
		ws:      ws,
		id:      ksuid.New().String(),
		chin:    make(chan []byte, 0),
		chout:   make(chan []byte, 0),
		onclose: nil,
	}

	return c
}

func (c *Connect) ID() string {
	return c.id
}

func (c *Connect) Close() {
	c.ws.Close()
	close(c.chin)
	close(c.chout) // fixme
}

func (c *Connect) OnMessage(cb OnMessage) {
	go c.ob(cb)
}

func (c *Connect) ob(cb OnMessage) {
	for {
		select {
		case bs, ok := <-c.chin:
			if !ok {
				log.Println(c.ID(), "ob closed")
				return
			}
			if cb != nil {
				cb(c.ID(),bs)
			}
		}
	}
}

func (c *Connect) OnClose(cb OnClose) {
	c.onclose = append(c.onclose, cb)
}

func (c *Connect) Send(data []byte) error {
	select {
	case c.chout <- data:
		return nil
	default:
		return ErrBlock
	}
}

// active will block till connect close
func (c *Connect) Active() {
	go c.wsSend()
	c.wsRecv() //that will block
	for _, oc := range c.onclose {
		if oc != nil {
			oc(c.ID())
		}
	}
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
				return
			}
			if bytes.Compare(buf, SignalClose) == 0 {
				log.Println(c.ID(), "signal close")
				return
			}
			err := websocket.Message.Send(c.ws, buf)
			if err != nil {
				log.Println(c.ws.RemoteAddr(), err)
				return
			}
		}
	}
	log.Println(c.id, "ws send quit")
}
