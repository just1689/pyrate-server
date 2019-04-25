package queues

import (
	"fmt"
	"github.com/nats-io/go-nats"
)

func GetNATSPublisher() (f func(subject string, msg []byte) error, err error) {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		fmt.Println(err)
		return
	}

	f = func(subject string, msg []byte) error {
		return nc.Publish(subject, msg)
	}
	return

}

func SubscribeToNATS(subject string) (outgoing chan []byte, stopMe chan bool, err error) {

	stopMe = make(chan bool)
	outgoing = make(chan []byte)

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Simple Async Subscriber
	_, err = nc.Subscribe(subject, func(m *nats.Msg) {
		outgoing <- m.Data
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	go func() {
		<-stopMe
		nc.Close()
	}()

	return

}
