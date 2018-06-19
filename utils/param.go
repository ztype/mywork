package utils

type Param struct {
	Id   string
	Type string
	Time int
	Data string
}

type Message struct {
	Name string
	MsgId string
	Param
}

type Respond struct{
	MsgId string
	Data interface{}
}
