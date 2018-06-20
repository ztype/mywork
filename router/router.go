package router

import (
	"fmt"
	"mywork/services"
	"mywork/utils"
	"sync"
	"time"
)

const NotifyTimeoutTime = time.Duration(time.Millisecond * 100)

type Service interface {
	Name() string
	Serve(param utils.Param) (interface{}, error)
	// if return value is not nil,all messages will notify to service
	ObserveChannel()chan<-utils.Param
}

type Router struct {
	lock     sync.Mutex
	services map[string]Service
	chans    []chan<- utils.Param
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
	return r
}

func (r *Router) Regist(s Service) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.services[s.Name()] = s
}

func (r *Router) Observe(c chan<- utils.Param) {
	if c != nil {
		r.chans = append(r.chans, c)
	}
}

func (r *Router) Stop(c chan<- utils.Param) {
	for i := 0; i < len(r.chans); i ++ {
		ch := r.chans[i]
		if ch == c {
			r.chans = append(r.chans[:i], r.chans[:i+1]...)
		}
	}
}

func (r *Router) notify(msg utils.Param) {
	for _, c := range r.chans {
		select {
		// send but do not block for it
		case c <- msg:
		default:
		}
	}
}

func (r *Router) Handle(msg *utils.Message) (interface{}, error) {
	r.notify(msg.Param)
	if s, ok := r.services[msg.Name]; ok {
		ret, err := s.Serve(msg.Param)
		res := utils.Respond{}
		res.Msgid = msg.Msgid
		res.Data = ret
		return res, err
	}
	return nil, fmt.Errorf("service [%s] not found", msg.Name)
}
