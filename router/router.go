package router

import (
	"encoding/json"
	"fmt"
	"log"
	"mywork/services"
	"mywork/utils"
	"sync"
	"time"

	"github.com/segmentio/ksuid"
)

const NotifyTimeoutTime = time.Duration(time.Millisecond * 100)

type Service interface {
	Name() string
	Serve(param *utils.Param) (interface{}, error)
	// if return value is not nil,all messages will notify to service
	ObserveChannel() chan *utils.Param
}

type Router struct {
	lock     sync.Mutex
	services map[string]Service
	chans    []chan *utils.Param
	chanIn   chan *utils.Message //chan<- *utils.Message
	chanOut  chan *utils.Message //<-chan *utils.Message
}

var DefaultRouter *Router

func init() {
	DefaultRouter = NewRouter()
	DefaultRouter.Regist(services.NewManager())
	DefaultRouter.Regist(services.NewRoomService())
}

func NewRouter() *Router {
	r := new(Router)
	r.services = make(map[string]Service, 0)
	r.chanIn = make(chan *utils.Message, 1)
	r.chanOut = make(chan *utils.Message, 1)
	go r.handleIn()
	return r
}

func (r *Router) Release() {
	for _, c := range r.chans {
		close(c)
	}
	close(r.chanOut)
}

func (r *Router) Regist(s Service) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.services[s.Name()] = s
	r.observe(s)
}

func (r *Router) observe(s Service) {
	c := s.ObserveChannel()
	if c != nil {
		r.chans = append(r.chans, c)
		r.serve(c, s.Name())
	}
}

func (r *Router) serve(c <-chan *utils.Param, name string) {
	for {
		select {
		case p, ok := <-c:
			if !ok {
				log.Println("service", name, "chan closed")
				break
			}
			m := new(utils.Message)
			m.Name = name
			m.Msgid = ksuid.New().String()
			m.Param = *p
			r.sendMsg(m)
		}
	}

}

func (r *Router) Stop(c chan<- *utils.Param) {
	for i := 0; i < len(r.chans); i++ {
		ch := r.chans[i]
		if ch == c {
			r.chans = append(r.chans[:i], r.chans[:i+1]...)
		}
	}
}

func (r *Router) notify(msg *utils.Param) {
	for _, c := range r.chans {
		select {
		// send but do not block for it
		case c <- msg:
		default:
		}
	}
}

func response(msg *utils.Message, v interface{}, err error) utils.Message {
	m := *msg
	// clean the data content for copy
	m.Data = ""
	m.Error = ""
	if v != nil {
		bs, e := json.Marshal(v)
		if e != nil {
			log.Println(e)
		}
		m.Data = string(bs)
	}
	if err != nil {
		m.Error = err.Error()
	}
	m.Time = int(time.Now().Unix())
	return m
}

func (r *Router) handle(msg *utils.Message) *utils.Message {
	r.notify(&msg.Param)
	if s, ok := r.services[msg.Name]; ok {
		ret, err := s.Serve(&msg.Param)
		res := response(msg, ret, err)
		return &res
	}
	res := response(msg, nil, fmt.Errorf("service [%s] not found", msg.Name))
	return &res
}

func (r *Router) ChannelIn() chan<- *utils.Message {
	return r.chanIn
}

func (r *Router) ChannelOut() <-chan *utils.Message {
	return r.chanOut
}

func (r *Router) sendMsg(m *utils.Message) {
	select {
	case r.chanOut <- m:
	default:
		log.Println("a message send fail")
	}
}

func (r *Router) handleIn() {
	for {
		select {
		case m, ok := <-r.chanIn:
			if !ok {
				log.Println("chan is closed")
				return
			}
			res := r.handle(m)
			r.sendMsg(res)
		}
	}
}
