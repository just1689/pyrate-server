package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/just1689/pyrate-server/chat"
	"github.com/just1689/pyrate-server/queues"
	"io/ioutil"
	"log"
	"net/http"
)

var addr = flag.String("address", ":8000", "")
var nAddr = flag.String("nsqaddress", "192.168.88.26:30000", "")
var lnAddr = flag.String("lnsqaddress", "http://192.168.88.26:30004", "")

func main() {

	flag.Parse()

	fmt.Println("Starting Pirate Server on", *addr)
	router := mux.NewRouter()
	chat.Serve(router, Subscriber)
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

func Subscriber(topic, channel string) (stopper chan bool) {
	stopper = make(chan bool)
	go func() {
		config := queues.Config{
			Topic:         topic,
			Channel:       channel,
			LookupAddress: *lnAddr,
			Address:       *nAddr,
			RemoteStopper: stopper,
			F: func(c *queues.Config, b []byte) {
				//TODO: implement
			},
		}
		queues.Subscribe(config)
	}()
	return
}
