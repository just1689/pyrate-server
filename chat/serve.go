package chat

import (
	"flag"
	"net/http"
)

func Serve(wsContext string) {
	flag.Parse()
	hub := newHub()
	go hub.run()
	http.HandleFunc(wsContext, func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
}
