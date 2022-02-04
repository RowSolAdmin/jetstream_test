# jetstream_test

## Description

The following is some sample go code to illustrate what may be a problem in nats, with 'interest' streams.

It appears we can have an stream of type interest, add consumers to take all messages from stream and acknowledge. The consumer states there are no further messages to process however
if we look at stream stats we see following

    Every 2.0s: nats stream report                                                                                                                                                                                                           johns-mbp.lan: Fri Feb  4 18:38:54 2022

    Obtaining Stream stats

    ╭──────────────────────────────────────────────────────────────────────────────────────╮
    │                                    Stream Report                                     │
    ├───────────────┬─────────┬───────────┬──────────┬─────────┬──────┬─────────┬──────────┤
    │ Stream        │ Storage │ Consumers │ Messages │ Bytes   │ Lost │ Deleted │ Replicas │
    ├───────────────┼─────────┼───────────┼──────────┼─────────┼──────┼─────────┼──────────┤
    │ test_interest │ Memory  │ 4         │ 30,016   │ 2.2 MiB │ 0    │ 9984    │          │
    ╰───────────────┴─────────┴───────────┴──────────┴─────────┴──────┴─────────┴──────────╯

In this and the other examples we are publishing 40,000 messages to the stream but sometimes these messages are not deleted. I was first aware of this when checking file storage consumption and was not expecting
any space to be used but it was - so looked a bit deeper.

This is a simplified example, this could use a Workqueue stream however have taken this approach to simplify what appears to be a sporadic issue with messages not being deleted.

I am using 2.7.1 nats server running in a k8s cluster as per provided helm charts, but I also witnessed this behaviour - at least on my environment with 2.7.0

nats cli 0.0.28

## Steps

NOTE::
Also tried with nightly build from 3rd February 22

Get latest docker image

    >> docker pull nats:latest

    >> docker run -p 4222:4222 -ti nats:latest -js

From root of this project via command prompt execute

    >> go run ./create/create.go

This will create the interest stream and four associated consumers as follows

    Every 2.0s: nats stream report                                                                                                                                                                                                           johns-mbp.lan: Fri Feb  4 17:56:10 2022

    Obtaining Stream stats

    ╭────────────────────────────────────────────────────────────────────────────────────╮
    │                                   Stream Report                                    │
    ├───────────────┬─────────┬───────────┬──────────┬───────┬──────┬─────────┬──────────┤
    │ Stream        │ Storage │ Consumers │ Messages │ Bytes │ Lost │ Deleted │ Replicas │
    ├───────────────┼─────────┼───────────┼──────────┼───────┼──────┼─────────┼──────────┤
    │ test_interest │ Memory  │ 4         │ 0        │ 0 B   │ 0    │ 0       │          │
    ╰───────────────┴─────────┴───────────┴──────────┴───────┴──────┴─────────┴──────────╯




    Every 2.0s: nats consumer report test_interest                                                                                                                                                                                           johns-mbp.lan: Fri Feb  4 17:56:32 2022

    ╭──────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────╮
    │                                                        Consumer report for test_interest with 4 consumers                                                        │
    ├─────────────────┬──────┬────────────┬──────────┬─────────────┬─────────────┬─────────────┬───────────┬───────────────────────────────────────────────────────────┤
    │ Consumer        │ Mode │ Ack Policy │ Ack Wait │ Ack Pending │ Redelivered │ Unprocessed │ Ack Floor │ Cluster                                                   │
    ├─────────────────┼──────┼────────────┼──────────┼─────────────┼─────────────┼─────────────┼───────────┼───────────────────────────────────────────────────────────┤
    │ TEST_CONSUMER_0 │ Pull │ Explicit   │ 15.00s   │ 0           │ 0           │ 0           │ 0         │ NCEZBHRD4ZTBV3UQH7BGTWKMUPA72HX2BYY4N5EYTUBDINMJNA7Z24H4* │
    │ TEST_CONSUMER_1 │ Pull │ Explicit   │ 15.00s   │ 0           │ 0           │ 0           │ 0         │ NCEZBHRD4ZTBV3UQH7BGTWKMUPA72HX2BYY4N5EYTUBDINMJNA7Z24H4* │
    │ TEST_CONSUMER_2 │ Pull │ Explicit   │ 15.00s   │ 0           │ 0           │ 0           │ 0         │ NCEZBHRD4ZTBV3UQH7BGTWKMUPA72HX2BYY4N5EYTUBDINMJNA7Z24H4* │
    │ TEST_CONSUMER_3 │ Pull │ Explicit   │ 15.00s   │ 0           │ 0           │ 0           │ 0         │ NCEZBHRD4ZTBV3UQH7BGTWKMUPA72HX2BYY4N5EYTUBDINMJNA7Z24H4* │
    ╰─────────────────┴──────┴────────────┴──────────┴─────────────┴─────────────┴─────────────┴───────────┴───────────────────────────────────────────────────────────╯

The steam 'test_interest' is configured to subject TEST.>

Each consumer is associated with a specific subject partition 0-3 as follows

    TEST_CONSUMER_0     TEST.a.b.c.0
    TEST_CONSUMER_1     TEST.a.b.c.1
    TEST_CONSUMER_2     TEST.a.b.c.2
    TEST_CONSUMER_3     TEST.a.b.c.3

Next, run the following

    >> go run ./publish/publish.go

This will publish 10,000 messages to each partition 0,1,2,3 so a total of 40,000 messages will be generated

    Every 2.0s: nats stream report                                                                                                                                                                                                           johns-mbp.lan: Fri Feb  4 17:57:44 2022

    Obtaining Stream stats

    ╭──────────────────────────────────────────────────────────────────────────────────────╮
    │                                    Stream Report                                     │
    ├───────────────┬─────────┬───────────┬──────────┬─────────┬──────┬─────────┬──────────┤
    │ Stream        │ Storage │ Consumers │ Messages │ Bytes   │ Lost │ Deleted │ Replicas │
    ├───────────────┼─────────┼───────────┼──────────┼─────────┼──────┼─────────┼──────────┤
    │ test_interest │ Memory  │ 4         │ 40,000   │ 2.9 MiB │ 0    │ 0       │          │
    ╰───────────────┴─────────┴───────────┴──────────┴─────────┴──────┴─────────┴──────────╯


    Every 2.0s: nats consumer report test_interest                                                                                                                                                                                           johns-mbp.lan: Fri Feb  4 17:58:06 2022

    ╭───────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────╮
    │                                                        Consumer report for test_interest with 4 consumers                                                         │
    ├─────────────────┬──────┬────────────┬──────────┬─────────────┬─────────────┬──────────────┬───────────┬───────────────────────────────────────────────────────────┤
    │ Consumer        │ Mode │ Ack Policy │ Ack Wait │ Ack Pending │ Redelivered │ Unprocessed  │ Ack Floor │ Cluster                                                   │
    ├─────────────────┼──────┼────────────┼──────────┼─────────────┼─────────────┼──────────────┼───────────┼───────────────────────────────────────────────────────────┤
    │ TEST_CONSUMER_0 │ Pull │ Explicit   │ 15.00s   │ 0           │ 0           │ 10,000 / 25% │ 0         │ NCEZBHRD4ZTBV3UQH7BGTWKMUPA72HX2BYY4N5EYTUBDINMJNA7Z24H4* │
    │ TEST_CONSUMER_1 │ Pull │ Explicit   │ 15.00s   │ 0           │ 0           │ 10,000 / 25% │ 0         │ NCEZBHRD4ZTBV3UQH7BGTWKMUPA72HX2BYY4N5EYTUBDINMJNA7Z24H4* │
    │ TEST_CONSUMER_2 │ Pull │ Explicit   │ 15.00s   │ 0           │ 0           │ 10,000 / 25% │ 0         │ NCEZBHRD4ZTBV3UQH7BGTWKMUPA72HX2BYY4N5EYTUBDINMJNA7Z24H4* │
    │ TEST_CONSUMER_3 │ Pull │ Explicit   │ 15.00s   │ 0           │ 0           │ 10,000 / 25% │ 0         │ NCEZBHRD4ZTBV3UQH7BGTWKMUPA72HX2BYY4N5EYTUBDINMJNA7Z24H4* │
    ╰─────────────────┴──────┴────────────┴──────────┴─────────────┴─────────────┴──────────────┴───────────┴───────────────────────────────────────────────────────────╯

As we can, as expected total of 40,000 messages split across the 4 partitions - each partition containing 10,000 messages

Finally execute the following

    >> go run ./consume/consume.go

This will output message received; if any error encountered it should stop - but at the end of the run we get
This will setup 4 subscription; one for each partition using the matching durable name / subject filter defined for the durable consumers above

The end of the output is shown below:

    99 46 48] 0xc00009e180 <nil> <nil> 1}
    &{TEST.a.b.c.0 $JS.ACK.test_interest.TEST_CONSUMER_0.1.40000.10000.1643997452070995600.0 map[] [112 97 114 116 105 116 105 111 110 58 58 32 91 48 93 32 99 111 117 110 116 101 114 32 58 58 32 91 54 50 55 57 93 32 45 45 32 84 69 83 84 46 97 46 98 46 99 46 48] 0xc00009e180 <nil> <nil> 1}

    >>> Counter:: 40000 [[Correct number of messages received as we would expect...]]

    durable 'TEST_CONSUMER_0' subject 'TEST.a.b.c.0'
    durable 'TEST_CONSUMER_1' subject 'TEST.a.b.c.1'
    durable 'TEST_CONSUMER_2' subject 'TEST.a.b.c.2'
    durable 'TEST_CONSUMER_3' subject 'TEST.a.b.c.3'

So - we have received the 40000 message as we would expect however if we look at stream stat report we have 30,016 messages outstanding on this particular run and, this number can vary...
However - we look at the consumer stats and they are all as we would expect, coupled with the check via consume.go at end we got the 40k messages as expected

    Every 2.0s: nats stream report                                                                                                                                                                                                           johns-mbp.lan: Fri Feb  4 18:00:30 2022

    Obtaining Stream stats

    ╭──────────────────────────────────────────────────────────────────────────────────────╮
    │                                    Stream Report                                     │
    ├───────────────┬─────────┬───────────┬──────────┬─────────┬──────┬─────────┬──────────┤
    │ Stream        │ Storage │ Consumers │ Messages │ Bytes   │ Lost │ Deleted │ Replicas │
    ├───────────────┼─────────┼───────────┼──────────┼─────────┼──────┼─────────┼──────────┤
    │ test_interest │ Memory  │ 4         │ 30,016   │ 2.2 MiB │ 0    │ 9984    │          │
    ╰───────────────┴─────────┴───────────┴──────────┴─────────┴──────┴─────────┴──────────╯
