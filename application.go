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
	router.HandleFunc("/", handleHome)
	setupStaticHost(router)

	http.Handle("/", router)
	log.Fatal(http.ListenAndServe(*addr, nil))

}

func setupStaticHost(router *mux.Router) {
	dir := "web/static/"
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(dir))))

}

func handleHome(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadFile("web/index.html")
	if err != nil {
		fmt.Println(err)
		return
	}
	w.Write(b)
}
