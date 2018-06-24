package services

import "testing"

func Test_roomid(t *testing.T) {
	id := roomid(52131)
	t.Log(id.String())

	for i := 0; i < 10; i++ {
		t.Log(newroomid().String())
	}
}
