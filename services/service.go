package service

import (
	"mywork/router"
	"mywork/base"
	"github.com/segmentio/ksuid"
	"mywork/services/room"
)

var sRoom = room.NewRoomService()

func init() {
	router.OnConnect = onConnect
}

func onConnect(c *router.Connect) {
	u := base.NewUser(ksuid.New().String(), c)
	sRoom.AddUser(u)
}
