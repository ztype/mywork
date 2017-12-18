package mywork

import (
	"fmt"
	"net/http"
)

func defaultHandle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello world")
}

func main() {
	fmt.Println("application started.")
	http.HandleFunc("/", defaultHandle)
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		fmt.Println(err)
	}
}
