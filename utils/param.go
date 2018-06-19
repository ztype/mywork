package utils

type Param struct {
	Id   string
	Type string
	Time int
	Data string
}

type Message struct {
	Name  string
	Msgid string
	Param
}

type Respond struct {
	Msgid string
	Data  interface{}
}
