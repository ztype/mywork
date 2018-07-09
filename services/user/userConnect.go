package user

import (
	"mywork/base"
	"sync"
)

var (
	mapUser map[string]*base.User
	ulock   sync.Mutex
)
