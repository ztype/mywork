package room

import (
	"fmt"
	"mywork/base"
	"sync"
)

const (
	JoinRoom   = "join_room"
	CreateRoom = "create_room"
	LeaveRoom  = "leave_room"
)


type Room struct {
	lock  sync.Mutex
	Id    RoomId
	Users map[string]*base.User
}

func newRoom() *Room {
	i := newRoomId()
	r := new(Room)
	r.Id = RoomId(i)
	r.Users = make(map[string]*base.User, 0)
	return r
}

func (r *Room) ID() RoomId {
	return r.Id
}

func (r *Room) UserIn(user *base.User) error {
	ul := r.UserCount()
	if ul >= 4 {
		return fmt.Errorf("room [%s] is full", r.Id.String())
	}
	if _, ok := r.Users[user.ID()]; ok {
		return nil
	}
	r.Users[user.ID()] = user
	return nil
}

func (r *Room) UserCount() int {
	return len(r.Users)
}

func (r *Room) UserOut(user *base.User) error {
	if r.UserCount() == 0 {
		return nil
	}
	u, ok := r.Users[user.ID()]
	if ok {
		delete(r.Users, u.ID())
	}
	return nil
}

func (r *Room) HasUser(uid string) bool {
	return r.GetUser(uid) != nil
}

func (r *Room) GetUser(uid string) *base.User {
	for id, u := range r.Users {
		if id == uid {
			return u
		}
	}
	return nil
}

func (r *Room) Close() {
	recycleRoomId(r.Id)
}
