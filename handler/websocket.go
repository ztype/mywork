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
	i := 0
	m := new(utils.Message)
	for {
		i++
		var buf = make([]byte, 0)
		err := websocket.Message.Receive(ws, &buf)
		if err != nil {
			log.Println(ws.Request().RemoteAddr, err)
			break
		}
		if err := json.Unmarshal(buf, m); err != nil {
			log.Println(err)
			continue
		}
		ret, err := router.DefaultRouter.Handle(m)
		if err != nil {
			websocket.Message.Send(ws, fmt.Sprintf("ERROR:%v", err))
			continue
		}
		j, err := json.Marshal(ret)
		err = websocket.Message.Send(ws, fmt.Sprintf("%s", string(j)))
		if err != nil {
			log.Println(ws.Request().RemoteAddr, err)
			break
		}
	}
}
