module test.com

go 1.17

replace test.com/constants => ./jetstream_test/constants

require github.com/nats-io/nats.go v1.13.1-0.20220121202836-972a071d373d

require (
	github.com/nats-io/nats-server/v2 v2.7.1 // indirect
	github.com/nats-io/nkeys v0.3.0 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	golang.org/x/crypto v0.0.0-20220112180741-5e0467b6c7ce // indirect
)
