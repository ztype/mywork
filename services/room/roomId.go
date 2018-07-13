package room

import (
	"fmt"
	"sync"
)

type RoomId int

var (
	idgen = []RoomId{}
	ilock = sync.Mutex{}
)

func (i RoomId) String() string {
	return fmt.Sprintf("%d", i)
}

func newRoomId() RoomId {
	i := 1
	for i = 1; ; i++ {
		found := false
		for _, s := range idgen {
			if s == RoomId(i) {
				found = true
				break
			}
		}
		if !found {
			ilock.Lock()
			idgen = append(idgen, RoomId(i))
			ilock.Unlock()
			return RoomId(i)
		}
	}
	return RoomId(i)
}

func recycleRoomId(id RoomId) {
	index := -1
	for i, s := range idgen {
		if s == id {
			index = i
			break
		}
	}
	if index != -1 {
		ilock.Lock()
		idgen = idgen[index : index+1]
		ilock.Unlock()
	}
}
