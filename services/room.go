package services

import (
	"mywork/utils"
	"sync"
)

const (
	JionRoom = "join_room"
	CreateRoom = "create_room"
)

type room struct {
	lock sync.Mutex
	Id    string
	Users map[string]*User
}

type RoomService struct {
	lock sync.Mutex
	Hall  map[string]*User
	Rooms map[string]*room
}

func NewRoomService() *RoomService {
	rs := new(RoomService)
	rs.Rooms = make(map[string]*room, 0)
	return rs
}

func (s *RoomService) Name() string {
	return "room"
}

func (s *RoomService) Serve(p utils.Param) (interface{}, error) {
	switch p.Type {
	case JionRoom:
	case CreateRoom:

	}
	return nil, nil
}

func (s *RoomService)ObserveChannel()chan<-utils.Param{
	return nil
}

func (s *RoomService) JionRoom(p utils.Param) (interface{},error){
	r := new(room)
	_ = r
	// todo
	return nil,nil
}
