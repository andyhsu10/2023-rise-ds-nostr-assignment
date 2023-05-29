package aggregator

import (
	"context"
	"distrise/internal/client"
	"distrise/internal/configs"
	"distrise/internal/databases"
	"distrise/internal/models"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/nbd-wtf/go-nostr"
	amqp "github.com/rabbitmq/amqp091-go"
)

type DistriseMessage struct {
	RelayURL string      `json:"relay_url"`
	Event    nostr.Event `json:"event"`
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func consumer(ch *amqp.Channel, name string) {
	// Connect to database
	db, err := databases.GetDB()
	if err != nil {
		log.Fatal(err)
	}

	// Consume message from RabbitMQ
	msgs, err := ch.Consume(
		name,  // queue
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	failOnError(err, "Failed to register a consumer")

	for d := range msgs {
		var msg DistriseMessage
		err := json.Unmarshal(d.Body, &msg)

		if err == nil {
			event := models.CoreEvent{
				RelayURL: msg.RelayURL,
				Data:     msg.Event.String(),
			}
			err := db.Create(&event).Error

			if err != nil {
				log.Println("Error:", err)
			} else {
				log.Printf("Event saved to database: %v\n\n", event)
			}
		}
	}
}

func publisher(ch *amqp.Channel, relayUrl string, name string) {
	// Initiate nostr client
	ctx := context.Background()
	relay, err := client.GetClient(ctx, relayUrl)
	if err != nil {
		panic(err)
	}

	// Create filters
	var filters = []nostr.Filter{{
		Kinds: []int{1}, // type 1 event (note)
	}}

	sub, err := relay.Subscribe(ctx, filters)
	if err != nil {
		panic(err)
	}

	for ev := range sub.Events {
		// handle returned event.
		// channel will stay open until the ctx is cancelled
		msg := DistriseMessage{
			RelayURL: relay.URL,
			Event:    *ev,
		}
		marshalMsg, err := json.Marshal(msg)

		// Publish message to RabbitMQ
		if err == nil {
			err = ch.PublishWithContext(ctx,
				"",    // exchange
				name,  // routing key
				false, // mandatory
				false, // immediate
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        marshalMsg,
				})
			failOnError(err, "Failed to publish a message")
		}
	}
}

func Aggregate() {
	queueName := "distrise_events"
	config := configs.GetConfig()

	// Connect to RabbitMQ, and initiate a queue
	conn, err := amqp.Dial(config.RabbitMqURL)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	_, err = ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	failOnError(err, "Failed to declare a queue named: "+queueName)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	// Start consumer
	go consumer(ch, queueName)

	// Start publisher
	for _, relayUrl := range config.RelayUrls {
		log.Printf("Starting publisher for relay: %v\n", relayUrl)
		go publisher(ch, relayUrl, queueName)
	}

	wg.Wait()
}

func View() {
	config := configs.GetConfig()

	// Connect to database
	db, err := databases.GetDB()
	if err != nil {
		log.Fatal(err)
	}

	filter := models.CoreEvent{
		RelayURL: config.RelayUrl,
	}

	var events []models.CoreEvent
	res := db.Where(&filter).Order("created_at DESC").Find(&events)
	if res.Error != nil {
		log.Fatal(res.Error)
	}

	for _, event := range events {
		fmt.Printf("Event: %v\n\n", event)
	}
}
