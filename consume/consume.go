package main

import (
	"fmt"
	"os"
	"sync/atomic"
	"time"

	"github.com/nats-io/nats.go"
	"test.com/constants"
)

var messagesReceived int32

// main -  consume
func main() {

	// connect to nats
	nc, err := nats.Connect("127.0.0.1:4222")
	if err != nil {
		fmt.Printf("unable to connect to nats, returned following error '%v'\n", err.Error())
		os.Exit(1)
	}

	// get jetstream connection
	jsc, err := nc.JetStream()
	if err != nil {
		fmt.Printf("unable to get jetstream client, returned following error '%v'\n", err)
		os.Exit(2)
	}

	items := make([]string, 0)

	for partition := 0; partition < constants.NumberPartitions; partition++ {

		// setup subject
		subject := constants.GetFilterSubject(partition)

		// construct know durable name
		durable := constants.GetDurableName(partition)

		// setup the pull subscription
		subscription, err := jsc.PullSubscribe(subject, durable)
		if err != nil {
			fmt.Printf("unable to get subscription on '%s', returned following error '%v'\n", subject, err)
			os.Exit(3)
		}

		items = append(items, fmt.Sprintf("durable '%s' subject '%s'", durable, subject))

		go getData(subscription, subject, durable)
	}

	for {
		time.Sleep(3 * time.Second)
		count := atomic.LoadInt32(&messagesReceived)
		fmt.Printf("Counter:: %d\n", count)

		// expect 20,000
		if count == constants.NumberPartitions*constants.NumberMessages {

			break
		} else {
			fmt.Printf("expected %d messages to be created, actual %d\n\n\n", constants.NumberPartitions*constants.NumberMessages, count)
		}
	}

	for _, item := range items {
		fmt.Println(item)
	}
}

// getData
func getData(subscription *nats.Subscription, subject, durable string) {

	for {
		messages, err := subscription.Fetch(10)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		for _, message := range messages {
			err := message.Ack()
			if err != nil {
				fmt.Printf("Acknowledge error for durable '%s' subject '%s' returned error '%v'\n\n\n", durable, subject, err)
				os.Exit(4)
			}
			fmt.Println(message)
			atomic.AddInt32(&messagesReceived, 1)
		}
	}
}
