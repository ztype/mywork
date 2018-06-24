package services

import (
	"fmt"
	"log"
	"mywork/base"
	"mywork/database"
	"mywork/utils"
	"sync"
)

const (
	JionRoom   = "join_room"
	CreateRoom = "create_room"
)

type room struct {
	lock  sync.Mutex
	Id    string
	Users map[string]*base.User
}

type RoomService struct {
	lock  sync.Mutex
	Rooms map[string]*room
	db    *database.DB
}

func NewRoomService() *RoomService {
	rs := new(RoomService)
	rs.Rooms = make(map[string]*room, 0)
	db, err := database.NewDatabase()
	if err != nil {
		log.Fatal(err)
	}
	rs.db = db
	return rs
}

func (s *RoomService) Name() string {
	return "room"
}

func (s *RoomService) Serve(p utils.Param) (interface{}, error) {
	switch p.Type {
	case JionRoom:
		return s.JionRoom(p)
	case CreateRoom:
		return s.CreateRoom(p)
	}
	return nil, fmt.Errorf("[%s] not found in %s", p.Type, s.Name())
}

func (s *RoomService) ObserveChannel() chan<- utils.Param {
	return nil
}

func (s *RoomService) JionRoom(p utils.Param) (interface{}, error) {
	r := new(room)
	_ = r
	// todo
	return nil, nil
}

func (s *RoomService) CreateRoom(p utils.Param) (interface{}, error) {
	return nil, nil
}
