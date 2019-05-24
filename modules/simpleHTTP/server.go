package simpleHTTP

import (
	"net/http"
)

func Server(port string) {
	http.HandleFunc("/write", RemoteWrtie)
	http.HandleFunc("/read", RemoteRead)
	http.ListenAndServe(":"+port, nil)
}
