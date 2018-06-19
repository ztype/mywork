package router

import (
	"fmt"
	"mywork/services"
	"mywork/utils"
	"sync"
)

type Service interface {
	Name() string
	Serve(param utils.Param) (interface{}, error)
}

type Router struct {
	lock     sync.Mutex
	services map[string]Service
}

var DefaultRouter *Router

func init() {
	DefaultRouter = NewRouter()
	DefaultRouter.Regist(services.NewManager())
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

func (r *Router) Handle(msg *utils.Message) (interface{}, error) {
	if s, ok := r.services[msg.Name]; ok {
		ret, err := s.Serve(msg.Param)
		res := utils.Respond{}
		res.Msgid = msg.Msgid
		res.Data = ret
		return res, err
	}
	return nil, fmt.Errorf("service [%s] not found", msg.Name)
}
