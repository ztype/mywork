package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"mywork/router"
	"net/http"

	"mywork/utils"

	"golang.org/x/net/websocket"
)

func ListenWebsocket(pattern, addr string) {
	http.Handle(pattern, websocket.Handler(websocketHandle))
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Println(err)
	}
}

func websocketHandle(ws *websocket.Conn) {
	serveWs(ws)
}

func serveWs(ws *websocket.Conn) {
	go doRecv(ws)
	doSend(ws)
}

func doSend(ws *websocket.Conn) {
	hout := router.DefaultRouter.ChannelOut()
	for {
		select {
		case m, ok := <-hout:
			if !ok {
				log.Println("hout closed")
				return
			}
			j, err := json.Marshal(m)
			if err != nil {
				log.Println(err)
				break
			}
			err = websocket.Message.Send(ws, fmt.Sprintf("%s", string(j)))
			if err != nil {
				log.Println(ws.Request().RemoteAddr, err)
				break
			}
		}
	}
}

func doRecv(ws *websocket.Conn) error {
	hin := router.DefaultRouter.ChannelIn()
	defer close(hin)
	for {
		m := new(utils.Message)
		var buf = make([]byte, 0)
		err := websocket.Message.Receive(ws, &buf)
		if err != nil {
			log.Println(ws.Request().RemoteAddr, err)
			return err
		}
		if err := json.Unmarshal(buf, m); err != nil {
			log.Println(err)
			continue
		}
		select {
		case hin <- m:
		default:
			log.Println("hin blocked")
		}
	}
	return nil
}
