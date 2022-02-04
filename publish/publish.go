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

// main -  responsible for publishing messages to 'test_interest'
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

	for partition := 0; partition < constants.NumberPartitions; partition++ {
		subject := constants.GetActualSubject(partition)

		for counter := 0; counter < constants.NumberMessages; counter++ {

			go func(p, c int, s string) {

				ack, err := jsc.Publish(subject, []byte(fmt.Sprintf("partition:: [%d] counter :: [%d] -- %s", p, c, s)))
				if err != nil {
					fmt.Printf("Error publishing to counter %d on partition %d\n\n\n", c, p)
					os.Exit(1)
				}

				// inc message counter
				atomic.AddInt32(&messagesReceived, 1)

				fmt.Printf("published message, partition [%d] // counter [%d] on subject '%s':: received acknowledge of '%v'\n\n", p, c, s, ack)
			}(partition, counter, subject)
		}
	}
	fmt.Println("message published successfully")

	for {
		time.Sleep(5 * time.Second)
		count := atomic.LoadInt32(&messagesReceived)
		fmt.Printf("message count:: %d\n", count)

		// expect 20,000
		if count == constants.NumberPartitions*constants.NumberMessages {
			break
		} else {
			fmt.Printf("expected %d messages to be created, actual %d", constants.NumberPartitions*constants.NumberMessages, count)
		}
	}
}
