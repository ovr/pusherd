package main

import (
	"flag"
	"github.com/streadway/amqp"
	"log"
	"bytes"
	"time"
)

var (
	configFile string
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	flag.StringVar(&configFile, "config", "./config.json", "Config filepath")
	flag.Parse()

	configuration := &Configuration{}
	configuration.Init(configFile)

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open channel")

	err = ch.ExchangeDeclare(
		"interpals",
		"fanout",
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare exchange")

	queue, err := ch.QueueDeclare(
		"push-notifications",
		true,
		false,
		true,
		false,
		nil,
	);
	failOnError(err, "Failed to declare queue")

	err = ch.QueueBind(
		queue.Name,
		"push-notifications",
		"interpals",
		false,
		nil,
	);
	failOnError(err, "Failed to bind queue")

	receiveChannel, err := ch.Consume(
		queue.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	for {
		for d := range receiveChannel {
			log.Printf("Received a message: %s", d.Body)
			dot_count := bytes.Count(d.Body, []byte("."))
			t := time.Duration(dot_count)
			time.Sleep(t * time.Second)
			log.Printf("Done")
			d.Ack(false)
		}
	}
}
