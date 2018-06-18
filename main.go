package main

import (
	"fmt"
	"mywork/handler"
	"net/http"
)

func main() {
	fmt.Println("application started.")
	go handler.ListenWebsocket("/ws", ":8081")
	http.HandleFunc("/", handler.DefaultHandle)
	//http.HandleFunc("/game", gameHandle)
	http.HandleFunc("/favicon.ico", handler.FaviconHandle)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
	}
}
