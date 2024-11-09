package main

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"net/http"
	"os"
	"time"
)

var rabbit_host = os.Getenv("RABBIT_HOST")
var rabbit_port = os.Getenv("RABBIT_PORT")
var rabbit_user = os.Getenv("RABBIT_USER")
var rabbit_password = os.Getenv("RABBIT_PASSWORD")

func main() {

	mux := http.NewServeMux()

	server := &http.Server{
		Addr:    ":80",
		Handler: mux,
	}

	mux.HandleFunc("POST /publish/{message}", submit)

	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("%s", err)
	}
	time.Sleep(time.Second)
}

func submit(writer http.ResponseWriter, request *http.Request) {
	message := request.PathValue("message")
	fmt.Println("Received messsage: " + message)

	conn, err := amqp.Dial("amqp://" + rabbit_user + ":" + rabbit_password + "@" + rabbit_host + ":" + rabbit_port + "/")

	if err != nil {
		log.Fatalf("%s: %s", "Failed to connnect to RabbitMQ", err)
	}

	defer conn.Close()

	ch, err := conn.Channel()

	if err != nil {
		log.Fatalf("%s: %s", "Failed to open a channel", err)
	}

	defer ch.Close()

	q, err := ch.QueueDeclare(
		"publisher", //name
		false,       //durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)

	if err != nil {
		log.Fatalf("%s: %s", "Failed to declare a queue", err)
	}

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})

	if err != nil {
		log.Fatalf("%s: %s", "Failed to publish a message", err)
	}

	fmt.Println("Published message successfully")
}
