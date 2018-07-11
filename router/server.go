package router

import (
	"log"
)

type Server interface {
	Name() string
	Channel() chan *Msg
}

type Msg struct {
	Cid  string
	Data []byte
}

var (
	mapServer map[string]Server
)

func Regist(s Server) {
	mapServer[s.Name()] = s
	go serve(s)
}

func serve(s Server) {
	c := s.Channel()
	if c != nil {
		for {
			select {
			case m, ok := <-c:
				if !ok {
					log.Println("service", s.Name(), "chan closed")
					return
				}
				SendTo(m.Cid, m.Data)
			}
		}
	}
}

func OnConnMsg(cid string, data []byte) {
	m :=& Msg{cid,data}
	for _,s := range mapServer{
		if c := s.Channel();c != nil {
			select{
			case c <- m:
			default:

			}
		}
	}
}
