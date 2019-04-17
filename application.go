package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/just1689/pyrate-server/chat"
	"log"
	"net/http"
)

var addr = flag.String("address", ":8000", "")

func main() {

	flag.Parse()

	fmt.Println("Starting Pirate Server on", *addr)
	router := mux.NewRouter()
	chat.Serve(router)
	http.Handle("/", router)
	log.Fatal(http.ListenAndServe(*addr, nil))

}
