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

type RoomService struct {
	lock    sync.Mutex
	Rooms   map[string]*Room
	db      *database.DB
	channel chan *utils.Param
}

func NewRoomService() *RoomService {
	rs := new(RoomService)
	rs.Rooms = make(map[string]*Room, 0)
	db, err := database.NewDatabase()
	if err != nil {
		log.Fatal(err)
	}
	rs.db = db
	rs.channel = make(chan *utils.Param, 2)
	return rs
}

func (s *RoomService) Name() string {
	return "room"
}

func (s *RoomService) Serve(p *utils.Param) (interface{}, error) {
	switch p.Type {
	case JoinRoom:
		return s.JoinRoom(p)
	case CreateRoom:
		return s.CreateRoom(p)
	}
	return nil, fmt.Errorf("[%s] not found in %s", p.Type, s.Name())
}

func (s *RoomService) ObserveChannel() chan *utils.Param {
	return s.channel
	return nil
}

func (s *RoomService) listen() {
	for {
		select {
		case p, ok := <-s.channel:
			if !ok {
				log.Println(s.Name(), "chan be closed")
				break
			}
			//todo
			s.Serve(p)
		}
	}
}

type roomMsg struct {
	RoomId string `json:"RoomId"`
}

func toRoomMsg(data string) (*roomMsg, error) {
	jm := new(roomMsg)
	if err := json.Unmarshal([]byte(data), jm); err != nil {
		return nil, err
	}
	return jm, nil
}

func (s *RoomService) JoinRoom(p *utils.Param) (interface{}, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	uid := p.Uid
	if s.GetUserRoom(uid) != 0 {
		return "you already are in a room", nil
	}

	user, err := s.db.GetUserById(uid)
	if err != nil {
		return nil, err
	}
	jm, err := toRoomMsg(p.Data)
	if err != nil {
		return nil, err
	}
	err = s.joinRoom(user, jm.RoomId)
	if err != nil {
		return nil, err
	}
	return "success", nil
}

func (s *RoomService) joinRoom(user *base.User, id string) error {
	//log.Println(s.Rooms)
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

func (s *RoomService) AllRooms() []string {
	ret := make([]string, 0)
	for _, r := range s.Rooms {
		ret = append(ret, r.ID().String())
	}
	return ret
}

func (s *RoomService) CreateRoom(p *utils.Param) (interface{}, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	user, err := s.db.GetUserById(p.Uid)
	if err != nil {
		return nil, err
	}
	r := newRoom()
	i := r.ID()
	s.Rooms[RoomId(i).String()] = r

	//add user to room
	if err := s.joinRoom(user, i.String()); err != nil {
		return nil, err
	}
	msg := make(map[string]interface{})
	msg["room_id"] = i
	msg["now_room_count"] = len(s.AllRooms())
	return msg, nil
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

func (s *RoomService) GetUserRoom(uid string) RoomId {
	for _, r := range s.Rooms {
		if r.HasUser(uid) {
			return r.ID()
		}
	}
	return 0
}

//func (s *RoomService)
