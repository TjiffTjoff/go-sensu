package main

import (
	"github.com/streadway/amqp"
	"log"
	"strconv"
	"strings"
)

/* Functions:
func Connect(server string, port int, user string, password string, vhost string) (*amqp.Connection, error)
Establishes an connection to the amqp server and returns amqp.Connection

func Channel(conn *amqp.Connection, msgType string) (*amqp.Channel, error)
Opens a channel and binds certain exchange to it. Returns amqp.Channel

func Send(channel *amqp.Channel,msgType string, msg []byte) error
Sends messages over channel and returns error
*/

// Function to establish connection to amqp server
func OpenConnection(config RabbitmqConf) (*amqp.Connection, error) {
	// Workaround to parse / in vhost name to %2F
	parsedVhost := strings.Replace(config.Vhost, "/", "%2F", -1)

	// Create a uri string from arguments
	uri := "amqp://" + config.User + ":" + config.Password + "@" + config.Host + ":" + strconv.Itoa(config.Port) + "/" + parsedVhost

	// Open a connection to the amqp server
	conn, err := amqp.Dial(uri)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

// Function to open channel on the amqp connection
func OpenChannel(conn *amqp.Connection, msgType string) (*amqp.Channel, error) {
	// Open a channel to communicate with the server
	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	// Declare the exchange to use when publishing
	if err := channel.ExchangeDeclare(
		msgType,
		"direct",
		false,
		false,
		false,
		false,
		nil,
	); err != nil {
		return nil, err
	}

	// Declare the queue to use when publishing
	channel.QueueDeclare(
		msgType,
		false,
		true,
		false,
		false,
		nil,
	)

	// Bind the queue to the exchange
	channel.QueueBind(
		msgType,
		"",
		msgType,
		false,
		nil,
	)

	return channel, nil
}

// Function to send keep alive message over specified channel
func Send(channel *amqp.Channel, msgType string, msgBody []byte) error {
	// Create the amqp message to publish
	msg := amqp.Publishing{
		ContentType:  "application/octet-stream",
		DeliveryMode: amqp.Persistent,
		Priority:     0,
		Body:         msgBody,
	}

	// Publish message to amqp server
	if err := channel.Publish(msgType, "", false, false, msg); err != nil {
		return err
	}

	// Returns nil as error if message was sent successfully
	return nil
}

func ListenAndSend(msgType string, goChannel chan []byte, amqpChannel *amqp.Channel) {
	for {
		var msgBody []byte
		msgBody = <-goChannel
		if err := Send(amqpChannel, msgType, msgBody); err != nil {
			log.Printf("Could not send: %s\n", err)
		}
	}
}
