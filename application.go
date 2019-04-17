package main

import (
	"flag"
	"fmt"
	"github.com/just1689/pyrate-server/chat"
	"log"
	"net/http"
)

var addr = flag.String("address", ":8080", "")
var ws = "/ws"

func main() {

	fmt.Println("Starting Pirate Server on", *addr)
	chat.Serve(ws)
	log.Fatal(http.ListenAndServe(*addr, nil))

}
