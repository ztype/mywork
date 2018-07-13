package router

import (
	"sync"
	"fmt"
	"log"
)

var (
	connMap = make(map[string]*Connect, 0)
	clock   sync.Mutex
)

func AddConnect(c *Connect) {
	clock.Lock()
	connMap[c.ID()] = c
	clock.Unlock()

	log.Println(c.ID(), "connect")

	c.OnClose(onClose)
	c.OnMessage(OnConnMsg)
	if OnConnect != nil {
		OnConnect(c)
	}
	c.Active()
}

var OnConnect = func(c *Connect) {

}

func SendTo(cid string, data []byte) error {
	if c, ok := connMap[cid]; ok {
		return c.Send(data)
	}
	return fmt.Errorf("connect %s not found")
}

func onClose(id string) {
	clock.Lock()
	delete(connMap, id)
	clock.Unlock()
	log.Println(id, "dis connected", len(connMap))
}
