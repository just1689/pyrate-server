package chat

import (
	"flag"
	"github.com/gorilla/mux"
	"net/http"
)

func Serve(router *mux.Router) {
	flag.Parse()
	hub := newHub()
	go hub.run()

	router.HandleFunc("/ws/{name:[a-z]+}/{secret:[a-z]+}", func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		name := params["name"]
		secret := params["secret"]
		serveWs(hub, w, r, name, secret)
	}).Methods("GET")

}
