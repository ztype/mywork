package room

import (
	"testing"
	"log"
)

func Test_roomid(t *testing.T) {
	id := RoomId(52131)
	t.Log(id.String())

	for i := 0; i < 10; i++ {
		t.Log(newRoomId().String())
	}
}

type A struct {
}

func (a *A) T() {
	log.Println("T from a")
}

type B struct {
}

func (b *B) T() {
	log.Println("T from B")
}

func TT() {
	log.Println("TT")
}

func Test_func(t *testing.T) {
	a := A{}
	//b := B{}
	a.T = TT
	a.T()
}
