package services

import "testing"

func Test_roomid(t *testing.T) {
	id := RoomId(52131)
	t.Log(id.String())

	for i := 0; i < 10; i++ {
		t.Log(newRoomId().String())
	}
}
