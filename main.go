package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	resource = "resource"
)

func defaultHandle(w http.ResponseWriter, r *http.Request) {
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
	tpl.ExecuteTemplate(w, "index.html", infos)
	//fmt.Fprintln(w, "hello world")
}

func faviconHandle(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, filepath.Join(".", "favicon.ico"))
	return
	bs, err := ioutil.ReadFile("." + "/favicon.ico")
	if err != nil {
		fmt.Println(err)
		return
	}
	w.Write(bs)
	return
}

func errHandle(w http.ResponseWriter, r *http.Request, err error) {
	fmt.Fprintf(w, "server internal error!"+fmt.Sprintf("%v", err))
}

func main() {
	fmt.Println("application started.")
	http.HandleFunc("/", defaultHandle)

	http.HandleFunc("/favicon.ico", faviconHandle)
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		fmt.Println(err)
	}
}
