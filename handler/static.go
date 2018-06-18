package handler

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func FaviconHandle(w http.ResponseWriter, r *http.Request) {
	bs, err := ioutil.ReadFile("./res" + "/favicon.ico")
	if err != nil {
		fmt.Println(err)
		return
	}
	w.Write(bs)
	return
}
