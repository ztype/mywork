package handler

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/segmentio/ksuid"
)

const (
	resource = "template"
)

const (
	FieldId = "uid"
)

func newCookie(domain string) *http.Cookie {
	id := ksuid.New().String()
	cookie := new(http.Cookie)
	cookie.Path = "/"
	cookie.Domain = domain
	cookie.Expires = time.Now().Add(time.Hour * 2)
	cookie.Value = id
	cookie.Name = FieldId
	return cookie
}

func errHandle(w http.ResponseWriter, r *http.Request, err error) {
	fmt.Fprintf(w, "server internal error!"+fmt.Sprintf("%v", err))
}

func DefaultHandle(w http.ResponseWriter, r *http.Request) {
	log.Println("request from:", r.RemoteAddr, r.URL.String(), r.URL.Path)
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
	w.Header().Add("Pragma", "no-cache")
	tpl.ExecuteTemplate(w, "index.html", 0)
}

func gameHandle(w http.ResponseWriter, r *http.Request) {
	p := `F:\work\src\github.com\deck-of-cards\example\index.html`
	http.ServeFile(w, r, p)
}
