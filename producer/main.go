package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/streadway/amqp"
)

func welcome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

func helloWorld(ch *amqp.Channel, q amqp.Queue) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body := "hellos"
		err := ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			})
		log.Printf(" [x] Sent %s", body)
		failOnError(err, "Failed to publish a message")

		fmt.Fprintf(w, "Hello World!")
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {

	conn, err := amqp.Dial("amqp://rabbitmq:rabbitmq@rabbit1:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	http.HandleFunc("/", welcome)
	http.HandleFunc("/hello", helloWorld(ch, q))

	fmt.Println("Started listening...")
	err = http.ListenAndServe(":4444", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
