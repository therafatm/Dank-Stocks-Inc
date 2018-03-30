package workcom

import (
	"os"
	"fmt"
	"bytes"
	"encoding/gob"
	"github.com/streadway/amqp"
	"test/workload_generator/commands"
	"common/utils"
)

const userTopic = "usertopic"
const otherRoute = "other"
const userRoute = "useroute"

type Workload interface {
	PublishCommand(command commands.Command)
	PublishUserRoute(route UserRoute)
	NameQueue()
}

type WorkloadConnection struct {
	Connection *amqp.Connection
	Channel *amqp.Channel
}

type UserRoute struct {
	Username string
	Host string
	Port string
}

func failOnError(err error, msg string) {
	if err != nil {
		utils.LogErr(err, msg)
		panic(err)
	}
}

func NewWorkloadConnection() (wconn *WorkloadConnection) {
	rabbitUser := os.Getenv("RABBITMQ_DEFAULT_USER")
	rabbitPass := os.Getenv("RABBITMQ_DEFAULT_PASS")
	rabbitHost := os.Getenv("RABBITMQ_COORD_HOST")
	rabbitPort := os.Getenv("RABBITMQ_COORD_PORT")
	url := fmt.Sprintf("amqp://%s:%s@%s:%s/", rabbitUser, rabbitPass, rabbitHost, rabbitPort)

	conn, err := amqp.Dial(url)
	failOnError(err, fmt.Sprintf("Failed to connect to Rabbit %s", url))
	wconn = &WorkloadConnection{Connection: conn}

	wconn.Channel, err = conn.Channel()
	failOnError(err, "Failed to open a channel")

	err = wconn.Channel.ExchangeDeclare(
		userTopic, // name
		"topic",
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare an exchange.")
	return
}

func (wconn *WorkloadConnection) PublishCommand(command commands.Command) {
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	err := enc.Encode(command)
	if err != nil {
		utils.LogErr(err, "Failed to encode message.")
	}

	err = wconn.Channel.Publish(
		userTopic,       // exchange
		command.Username, 	// routing key
		false,              // mandatory
		false,              // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        buffer.Bytes(),
		})
	if err != nil {
		utils.LogErr(err, "Failed to publish log message.")
	}
}

func (wconn *WorkloadConnection) PublishOtherCommand(command commands.Command) {
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	err := enc.Encode(command)
	if err != nil {
		utils.LogErr(err, "Failed to encode message.")
	}

	err = wconn.Channel.Publish(
		usertopic,       // exchange
		otherRoute, 			// routing key
		false,              // mandatory
		false,              // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        buffer.Bytes(),
		})
	if err != nil {
		utils.LogErr(err, "Failed to publish log message.")
	}
}

func (wconn *WorkloadConnection) PublishUserRoute(route UserRoute) {
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	err := enc.Encode(route)
	if err != nil {
		utils.LogErr(err, "Failed to encode message.")
	}

	err = wconn.Channel.Publish(
		userTopic,       // exchange
		userRoute, 	// routing key
		false,              // mandatory
		false,              // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        buffer.Bytes(),
		})
	if err != nil {
		utils.LogErr(err, "Failed to publish log message.")
	}
}

func (wconn *WorkloadConnection) NameQueue() amqp.Queue {
     q, err := wconn.Channel.QueueDeclare(
            "",    // name
            false, // durablea
            false, // delete when usused
            true,  // exclusive
            false, // no-wait
            nil,   // arguments
    )
    failOnError(err, "Failed to declare a queue")
    err = wconn.Channel.QueueBind(
    		q.Name,
    		userTopic,
    		userRoute,
    		false,
    		nil)
    failOnError(err, "Failed to bind queue")
    return q
}


