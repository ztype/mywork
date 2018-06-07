package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	resource = "resource"
)

func defaultHandle(w http.ResponseWriter, r *http.Request) {
	log.Println("request from:", r.RemoteAddr,r.URL.String())
	tpl, err := template.ParseFiles(filepath.Join(resource, "head.html"),
		filepath.Join(resource, "body.html"),
		filepath.Join(resource, "index.html"))
	if err != nil {
		errHandle(w, r, err)
		return
	}
	infos := []os.FileInfo{}
	walkfunc := func(path string, info os.FileInfo, err error) error {
		if err != nil || info == nil {
			return err
		}
		if strings.Contains(path, ".git") {
			return nil
		}
		infos = append(infos, info)
		return nil
	}
	filepath.Walk("./", walkfunc)
	w.Header().Add("Pragma","no-cache")
	tpl.ExecuteTemplate(w, "index.html", infos)

	//fmt.Fprintln(w, "hello world")
}

func faviconHandle(w http.ResponseWriter, r *http.Request) {
	bs, err := ioutil.ReadFile("." + "/favicon.ico")
	if err != nil {
		fmt.Println(err)
		return
	}
	w.Write(bs)
	return
}

func gameHandle(w http.ResponseWriter, r *http.Request) {
	p := `F:\work\src\github.com\deck-of-cards\example\index.html`
	http.ServeFile(w, r, p)
}

func websocketHandle(ws *websocket.Conn) {
	client := ws.Request().RemoteAddr
	log.Println("Client connected:", client)
	serveWs(ws)
}

func serveWs(ws *websocket.Conn) {
	i := 0
	for {
		i++
		var msg string
		err := websocket.Message.Receive(ws,&msg)
		if err != nil {
			log.Println(ws.Request().RemoteAddr,err)
			break
		}
		log.Println("from:",ws.Request().RemoteAddr,":",msg)
		err = websocket.Message.Send(ws,fmt.Sprintf("redten %d",i))
		if err != nil {
			log.Println(ws.Request().RemoteAddr,err)
			break
		}
	}
}

func errHandle(w http.ResponseWriter, r *http.Request, err error) {
	fmt.Fprintf(w, "server internal error!"+fmt.Sprintf("%v", err))
}

func main() {
	fmt.Println("application started.")
	go listenWebsocket()
	http.HandleFunc("/", defaultHandle)
	http.HandleFunc("/game", gameHandle)
	http.HandleFunc("/favicon.ico", faviconHandle)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
	}
}

func listenWebsocket(){
	http.Handle("/ws", websocket.Handler(websocketHandle))
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		fmt.Println(err)
	}
}
