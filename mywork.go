package main

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
)

const (
	resource = "resource"
)

func defaultHandle(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.ParseFiles(filepath.Join(resource, "head.html"),
		filepath.Join(resource, "body.html"),
		filepath.Join(resource, "index.html"))
	if err != nil {
		errHandle(w, r)
		return
	}
	tpl.ExecuteTemplate(w, "index.html", nil)
	//fmt.Fprintln(w, "hello world")
}

func errHandle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "server internal error!")
}

func main() {
	fmt.Println("application started.")
	http.HandleFunc("/", defaultHandle)
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		fmt.Println(err)
	}
}
