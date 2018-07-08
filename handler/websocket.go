package handler

import (
	"fmt"
	"mywork/router"
	"net/http"

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
	c := router.NewConnect(ws)
	router.AddConnect(c)
}
