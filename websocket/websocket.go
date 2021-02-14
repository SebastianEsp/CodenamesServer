package websocket

import (
	"net/http"

	"github.com/gorilla/mux"
)

func SetupRoutes() *mux.Router {
	hub := newHub()
	go hub.run()

	r := mux.NewRouter()
	r.HandleFunc("/ws2", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	r.HandleFunc("/ws2/{game}", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	return r
}
