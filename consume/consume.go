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
var stop int32

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
	subs := make([]*nats.Subscription, 0)

	for partition := 0; partition < constants.NumberPartitions; partition++ {

		// setup subject
		subject := constants.GetActualSubject(partition)

		// construct know durable name
		durable := constants.GetDurableName(partition)

		// setup the pull subscription
		subscription, err := jsc.PullSubscribe(subject, durable, nats.Bind(constants.StreamName, durable))
		if err != nil {
			fmt.Printf("unable to get subscription on '%s', returned following error '%v'\n", subject, err)
			os.Exit(3)
		}

		items = append(items, fmt.Sprintf("durable '%s' subject '%s'", durable, subject))

		subs = append(subs, subscription)

		go getData(subscription, subject, durable)
	}

	for {
		time.Sleep(3 * time.Second)
		count := atomic.LoadInt32(&messagesReceived)
		fmt.Printf("Counter:: %d\n", count)

		// expect 40,000
		if count == constants.NumberPartitions*constants.NumberMessages {

			atomic.AddInt32(&stop, 1)

			break
		} else {
			fmt.Printf("expected %d messages to be created, actual %d\n\n\n", constants.NumberPartitions*constants.NumberMessages, count)
		}
	}

	for _, item := range items {
		fmt.Println(item)
	}

	for _, sub := range subs {
		err := sub.Unsubscribe()
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

// getData
func getData(subscription *nats.Subscription, subject, durable string) {

	for {
		check := atomic.LoadInt32(&stop)

		if check == 0 {
			messages, err := subscription.Fetch(1000)
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
		} else {
			break
		}
	}
}
