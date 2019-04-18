package chat

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func Serve(router *mux.Router, subscriber func(topic, channel string) chan bool) {
	flag.Parse()
	hub := newHub()
	go hub.run()

	router.HandleFunc("/ws/{name:[a-z]+}/{secret:[a-z]+}", func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		name := params["name"]
		secret := params["secret"]
		fmt.Println("Connected", name, secret)
		serveWs(hub, w, r, name, secret, subscriber)
	}).Methods("GET")

}
