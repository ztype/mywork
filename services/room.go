package services

import (
	"encoding/json"
	"fmt"
	"log"
	"mywork/base"
	"mywork/database"
	"mywork/utils"
	"sync"
)

const (
	JoinRoom   = "join_room"
	CreateRoom = "create_room"
	LeaveRoom  = "leave_room"
)

type roomid int

var idgen = []int{}

type Room struct {
	lock  sync.Mutex
	Id    roomid
	Users map[string]*base.User
}

type RoomService struct {
	lock  sync.Mutex
	Rooms map[string]*Room
	db    *database.DB
}

func (i roomid) String() string {
	return fmt.Sprintf("%04d", i)
}

func NewRoomService() *RoomService {
	rs := new(RoomService)
	rs.Rooms = make(map[string]*Room, 0)
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
	case JoinRoom:
		return s.JoinRoom(p)
	case CreateRoom:
		return s.CreateRoom(p)
	}
	return nil, fmt.Errorf("[%s] not found in %s", p.Type, s.Name())
}

func (s *RoomService) ObserveChannel() chan<- utils.Param {
	return nil
}

type roomMsg struct {
	RoomId string `json:"room_id"`
}

func toRoomMsg(data string) (*roomMsg, error) {
	jm := new(roomMsg)
	if err := json.Unmarshal([]byte(data), jm); err != nil {
		return nil, err
	}
	return jm, nil
}

func (s *RoomService) JoinRoom(p utils.Param) (interface{}, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	user, err := s.db.GetUserById(p.Uid)
	if err != nil {
		return nil, err
	}
	jm, err := toRoomMsg(p.Data)
	if err != nil {
		return nil, err
	}
	s.joinRoom(user, jm.RoomId)
	return nil, nil
}

func (s *RoomService) joinRoom(user *base.User, id string) error {
	if r, ok := s.Rooms[id]; ok {
		return r.UserIn(user)
	}
	return fmt.Errorf("room [%s] not found", id)
}

func (s *RoomService) getRoom(id string) *Room {
	if r, ok := s.Rooms[id]; ok {
		return r
	}
	return nil
}

func (s *RoomService) AllRooms()[]string{
	ret := make([]string,0)
	for _,r := range s.Rooms{
		ret = append(ret,r.ID().String())
	}
	return ret
}

func (s *RoomService) CreateRoom(p utils.Param) (interface{}, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	user, err := s.db.GetUserById(p.Uid)
	if err != nil {
		return nil, err
	}
	r := newRoom()
	i := r.ID()
	s.Rooms[roomid(i).String()] = r

	//add user to room
	if err := s.joinRoom(user, i.String()); err != nil {
		return nil, err
	}
	msg := make(map[string]interface{})
	msg["room_id"] = i
	msg["now_room_count"] = len(s.AllRooms())
	return msg, nil
}

func newroomid() roomid {
	i := 0
	for i = 1; ; i++ {
		found := false
		for _, s := range idgen {
			if s == i {
				found = true
				break
			}
		}
		if !found {
			idgen = append(idgen, i)
			return roomid(i)
		}
	}
	return roomid(i)
}

func (s *RoomService) LeaveRoom(p utils.Param) (interface{}, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	user, err := s.db.GetUserById(p.Uid)
	if err != nil {
		return nil, err
	}
	rm, err := toRoomMsg(p.Data)
	if err != nil {
		return nil, err
	}
	err = s.leaveRoom(user, rm.RoomId)
	if err != nil {
		return nil, err
	}
	return "ok", nil
}

func (s *RoomService) leaveRoom(user *base.User, rid string) error {
	r := s.getRoom(rid)
	if r == nil {
		return fmt.Errorf("room [%s] not found", rid)
	}
	if err := r.UserOut(user); err != nil {
		return err
	}
	if r.UserCount() == 0 {
		s.lock.Lock()
		delete(s.Rooms, r.Id.String())
		s.lock.Unlock()
	}
	return nil
}

//========= room area =============//
func newRoom() *Room {
	i := newroomid()
	r := new(Room)
	r.Id = roomid(i)
	r.Users = make(map[string]*base.User, 0)
	return r
}

func (r *Room) ID() roomid {
	return r.Id
}

func (r *Room) UserIn(user *base.User) error {
	ul := r.UserCount()
	if ul >= 4 {
		return fmt.Errorf("room [%s] is full", r.Id.String())
	}
	if _, ok := r.Users[user.Id()]; ok {
		return nil
	}
	r.Users[user.Id()] = user
	return nil
}

func (r *Room) UserCount() int {
	return len(r.Users)
}

func (r *Room) UserOut(user *base.User) error {
	if r.UserCount() == 0 {
		return nil
	}
	u, ok := r.Users[user.Id()]
	if ok {
		delete(r.Users, u.Id())
	}
	return nil
}
