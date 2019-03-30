package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Jhamm1/OutboundCommunication/ccmServiceClient"
	. "github.com/Jhamm1/OutboundCommunication/db"
	. "github.com/Jhamm1/OutboundCommunication/models"
	"github.com/streadway/amqp"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var dao = CommunicationsDAO{}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func SendingQueue(msg []byte) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"Outbound_communication_sending", // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	body := msg
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	log.Printf(" [x] Sent %s", body)
	failOnError(err, "Failed to publish a message")
}

//ToDO: Write function for consumingCommunicationOffQueue
func consumeOfftheQueue() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"Outbound_communication_service", // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue"+q.Name)

	msgs, err := ch.Consume(
		q.Name,    // queue
		"Sending", // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)

			var comms Communication
			json.Unmarshal([]byte(d.Body), &comms)
			ccmServiceClient.SendMail(string(comms.Email), string(comms.FirstName))
			SendingQueue(d.Body)

			//Update Outbound DB
			session, err := mgo.Dial("mongodb://localhost:27017")
			if err != nil {
				log.Fatal(err)
			}
			c := session.DB("outbound-communications_db").C("communications")
			selector := bson.M{"_id": comms.ID}
			updator := bson.M{"$set": bson.M{"status": "completedRequest"}}
			err = c.Update(selector, updator) //&communication.Email)
			if err := c.Update(selector, updator); err != nil {
				panic(err)
			}
			fmt.Printf("%+v\n", comms.ID)
		}

	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func main() {
	consumeOfftheQueue()
}
