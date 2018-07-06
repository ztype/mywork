package router

import (
	"sync"
)

type Server interface{
	Name() string
	Channel() chan []byte
}

var (
	slock sync.Mutex
	mapServer map[string]Server
)

func Regist(s Server){

}