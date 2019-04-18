package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/just1689/pyrate-server/chat"
	"io/ioutil"
	"log"
	"net/http"
)

var addr = flag.String("address", ":8000", "")

func main() {

	flag.Parse()

	fmt.Println("Starting Pirate Server on", *addr)
	router := mux.NewRouter()
	chat.Serve(router)
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadFile("web/index.html")
		if err != nil {
			fmt.Println(err)
			return
		}
		w.Write(b)
	})
	router.HandleFunc("/py.js", func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadFile("web/py.js")
		if err != nil {
			fmt.Println(err)
			return
		}
		w.Write(b)
	})

	http.Handle("/", router)
	log.Fatal(http.ListenAndServe(*addr, nil))

}
