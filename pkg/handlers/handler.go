package handlers

import "net/http"

func Test(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is my home page"))
}
