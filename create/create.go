package main

import (
	"fmt"
	"os"
	"time"

	"github.com/nats-io/nats.go"
	"test.com/constants"
)

// main -  responsible for creating an 'interest' stream on local nats / creating 2 durable consumers
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

	// create interest stream
	stream, err := jsc.AddStream(&nats.StreamConfig{
		Name:      constants.StreamName,
		MaxAge:    0,
		Subjects:  []string{"TEST.>"},
		Storage:   nats.MemoryStorage,
		Replicas:  1,
		Retention: nats.InterestPolicy,
	})

	// check error
	if err != nil {
		fmt.Printf("unable to add stream , returned following error '%v'\n\n\n", err)
		os.Exit(3)
	} else {
		fmt.Printf("Steam Information:: '%v'\n\n\n", stream)
	}

	for partition := 0; partition < constants.NumberPartitions; partition++ {

		// setup subject
		subject := constants.GetActualSubject(partition)

		// construct know durable name
		durable := constants.GetDurableName(partition)

		consumer, err := jsc.AddConsumer(
			constants.StreamName,
			&nats.ConsumerConfig{
				Durable:       durable,
				AckPolicy:     nats.AckExplicitPolicy,
				FilterSubject: subject,
				AckWait:       time.Second * 15,
			},
		)

		if err != nil {
			fmt.Printf("unable to add durable %s, returned following error '%v'\n\n\n", durable, err)
			os.Exit(4)
		} else {
			fmt.Printf("consumer %s, '%v'\n\n\n", durable, consumer)
		}
	}

	fmt.Println("stream / durable consumers created successfully")
}
