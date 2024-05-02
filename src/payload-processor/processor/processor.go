package processor

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Processor struct {
	ch        *amqp.Channel
	queueName string
}

func NewProcessor(ch *amqp.Channel, queueName string) Processor {
	log.Println("Created processor")
	return Processor{
		ch:        ch,
		queueName: queueName,
	}
}

func (p Processor) ProcessMessages() {
	msgs, err := p.ch.Consume(
		p.queueName, // queue
		"",          // consumer
		true,        // auto-ack
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)
	if err != nil {
		log.Fatal("Failed to register a consumer:", err)
	}

	for msg := range msgs {
		log.Printf("Received a message: %s", msg.Body)
	}
}
