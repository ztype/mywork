package main

import (
	"fmt"
	"github.com/segmentio/ksuid"
	"golang.org/x/net/websocket"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"time"
	"encoding/json"
)

const (
	resource = "template"
)

const (
	FieldId = "uid"
)

var manager = NewManager()

func newCookie(domain string) *http.Cookie {
	id := ksuid.New().String()
	cookie := new(http.Cookie)
	cookie.Path = "/"
	cookie.Domain = domain
	//cookie.HttpOnly = true
	cookie.Expires = time.Now().Add(time.Hour * 2)
	cookie.Value = id
	cookie.Name = FieldId
	return cookie
}

func defaultHandle(w http.ResponseWriter, r *http.Request) {
	log.Println("request from:", r.RemoteAddr, r.URL.String())
	tpl, err := template.ParseFiles(filepath.Join(resource, "head.html"),
		filepath.Join(resource, "body.html"),
		filepath.Join(resource, "index.html"))
	if err != nil {
		errHandle(w, r, err)
		return
	}

	c, err := r.Cookie(FieldId)
	if err != nil {
		if err == http.ErrNoCookie {
			c = newCookie(r.Host)
			http.SetCookie(w, c)
		} else {
			errHandle(w, r, err)
			return
		}
	}

	if !manager.IsExist(c.Value) {
		c = newCookie(r.Host)
		http.SetCookie(w, c)
	}
	manager.UserConnect(c.Value)
	w.Header().Add("Pragma", "no-cache")
	total := manager.UserCount()
	tpl.ExecuteTemplate(w, "index.html", total)
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
	serveWs(ws)
}

type message struct {
	Id   string
	Type string
	Time int
	Data string
}

func handleMessage(m *message) {
	if user := manager.GetUser(m.Id); user != nil {
		user.HeartBeat()
	}
}

func serveWs(ws *websocket.Conn) {
	i := 0
	m := new(message)
	for {
		i++
		var buf = make([]byte, 0)
		err := websocket.Message.Receive(ws, &buf)
		if err != nil {
			log.Println(ws.Request().RemoteAddr, err)
			break
		}
		if err := json.Unmarshal(buf, m); err == nil {
			handleMessage(m)
		}
		err = websocket.Message.Send(ws, fmt.Sprintf("redten %d", i))
		if err != nil {
			log.Println(ws.Request().RemoteAddr, err)
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
	//http.HandleFunc("/favicon.ico", faviconHandle)
	http.FileServer(http.Dir("./res"))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
	}
}

func listenWebsocket() {
	http.Handle("/ws", websocket.Handler(websocketHandle))
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		fmt.Println(err)
	}
}
