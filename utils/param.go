package utils

type Param struct {
	Uid  string
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
