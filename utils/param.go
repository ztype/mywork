package utils

import "time"

type Param struct {
	Uid   string
	Type  string
	Error string
	Data  string
}

type Message struct {
	Name  string
	Msgid string
	Time  int
	Param
}

type Respond struct {
	Msgid string
	Time  int
	Data  interface{}
}

func RespondFromMsg(m *Message)*Respond{
	r := new(Respond)
	r.Msgid = m.Msgid
	r.Time = int(time.Now().Unix())
	return r
}