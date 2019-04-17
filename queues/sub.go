package queues

import (
	"flag"
	"github.com/nsqio/go-nsq"
	"log"
)

var (
	maxInFlight   = flag.Int("max-in-flight", 200, "max number of messages to allow in flight")
	totalMessages = flag.Int("n", 0, "total messages to show (will wait if starved)")
)

type StringArray []string

type Config struct {
	Address       string
	LookupAddress string
	Topic         string
	Channel       string
	F             ReceiveFunction
	RemoteStopper chan bool
}

type tailHandler struct {
	topicName     string
	totalMessages int
	messagesShown int
	config        Config
}

func (th *tailHandler) HandleMessage(m *nsq.Message) error {
	th.messagesShown++
	th.config.F(&th.config, m.Body)
	return nil
}

type ReceiveFunction func(c *Config, b []byte)

func Subscribe(c Config) {
	cfg := nsq.NewConfig()
	var nsqdTCPAddrs = []string{}
	nsqdTCPAddrs = append(nsqdTCPAddrs, c.Address)

	flag.Var(&nsq.ConfigFlag{cfg}, "consumer-opt", "http://godoc.org/github.com/nsqio/go-nsq#Config")
	cfg.MaxInFlight = *maxInFlight

	consumer, err := nsq.NewConsumer(c.Topic, c.Channel, cfg)
	if err != nil {
		log.Println(err)
	}
	consumer.AddHandler(&tailHandler{topicName: c.Topic, totalMessages: *totalMessages, config: c})

	err = consumer.ConnectToNSQDs(nsqdTCPAddrs)
	if err != nil {
		log.Println(err)
	}

	arr := []string{}
	arr = append(arr, c.LookupAddress)
	err = consumer.ConnectToNSQLookupds(arr)
	if err != nil {
		log.Println(err)
	}

	//Block until told to stop
	<-c.RemoteStopper

	consumer.Stop()
	<-consumer.StopChan

}
