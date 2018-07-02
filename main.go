package main

import (
	"fmt"
	"log"
	"mywork/handler"
	"net/http"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	log.Println("application started.")

	go handler.ListenWebsocket("/ws", ":8081")

	http.HandleFunc("/", handler.DefaultHandle)
	http.HandleFunc("/favicon.ico", handler.FaviconHandle)
	err := http.ListenAndServe("", nil)
	if err != nil {
		fmt.Println(err)
	}
}
