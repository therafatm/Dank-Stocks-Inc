package logging

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"common/models"
	"common/utils"

	"github.com/streadway/amqp"
)

type Command string

const (
	ADD              = Command("ADD")
	QUOTE            = Command("QUOTE")
	BUY              = Command("BUY")
	COMMIT_BUY       = Command("COMMIT_BUY")
	CANCEL_BUY       = Command("CANCEL_BUY")
	SELL             = Command("SELL")
	COMMIT_SELL      = Command("COMMIT_SELL")
	CANCEL_SELL      = Command("CANCEL_SELL")
	SET_BUY_AMOUNT   = Command("SET_BUY_AMOUNT")
	CANCEL_SET_BUY   = Command("CANCEL_SET_BUY")
	SET_BUY_TRIGGER  = Command("SET_BUY_TRIGGER")
	SET_SELL_AMOUNT  = Command("SET_SELL_AMOUNT")
	SET_SELL_TRIGGER = Command("SET_SELL_TRIGGER")
	CANCEL_SET_SELL  = Command("CANCEL_SET_SELL")
	DUMPLOG          = Command("DUMPLOG")
	DISPLAY_SUMMARY  = Command("DISPLAY_SUMMARY")
)

var validCommands = map[Command]bool{
	ADD:              true,
	QUOTE:            true,
	BUY:              true,
	COMMIT_BUY:       true,
	CANCEL_BUY:       true,
	SELL:             true,
	COMMIT_SELL:      true,
	CANCEL_SELL:      true,
	SET_BUY_AMOUNT:   true,
	CANCEL_SET_BUY:   true,
	SET_BUY_TRIGGER:  true,
	SET_SELL_AMOUNT:  true,
	SET_SELL_TRIGGER: true,
	CANCEL_SET_SELL:  true,
	DUMPLOG:          true,
	DISPLAY_SUMMARY:  true}

type Message struct {
	UserCommand        *UserCommandType        `xml:"userCommand"`
	AccountTransaction *AccountTransactionType `xml:"accountTransaction"`
	SystemEvent        *SystemEventType        `xml:"systemEvent"`
	QuoteServer        *QuoteServerType        `xml:"quoteServer"`
	ErrorEvent         *ErrorEventType         `xml:"errorEvent"`
}

type UserCommandType struct {
	XMLName           string  `xml:"userCommand"`
	Timestamp         int64   `xml:"timestamp"`
	Server            string  `xml:"server"`
	TransactionNumber int64   `xml:"transactionNum"`
	Command           Command `xml:"command"`
	Username          string  `xml:"username,omitempty"`
	Symbol            string  `xml:"stockSymbol,omitempty"`
	Filename          string  `xml:"filename,omitempty"`
	Funds             string  `xml:"funds,omitempty"`
}

type AccountTransactionType struct {
	XMLName           string `xml:"accountTransaction"`
	Timestamp         int64  `xml:"timestamp"`
	Server            string `xml:"server"`
	TransactionNumber int64  `xml:"transactionNum"`
	Action            string `xml:"action"`
	Username          string `xml:"username"`
	Funds             string `xml:"funds"`
}

type SystemEventType struct {
	XMLName           string `xml:"systemEvent"`
	Timestamp         int64  `xml:"timestamp"`
	Server            string `xml:"server"`
	TransactionNumber int64  `xml:"transactionNum"`
	Command           string `xml:"command"`
	Username          string `xml:"username"`
	Symbol            string `xml:"stockSymbol"`
	Funds             string `xml:"funds"`
}

type QuoteServerType struct {
	XMLName           string `xml:"quoteServer"`
	Timestamp         int64  `xml:"timestamp"`
	Server            string `xml:"server"`
	TransactionNumber int64  `xml:"transactionNum"`
	QuoteServerTime   int64  `xml:"quoteServerTime"`
	Username          string `xml:"username"`
	Symbol            string `xml:"stockSymbol"`
	Price             string `xml:"price"`
	CryptoKey         string `xml:"cryptokey"`
}

type ErrorEventType struct {
	XMLName           string  `xml:"errorEvent"`
	Timestamp         int64   `xml:"timestamp"`
	Server            string  `xml:"server"`
	TransactionNumber int64   `xml:"transactionNum"`
	Command           Command `xml:"command"`
	Username          string  `xml:"username,omitempty"`
	Symbol            string  `xml:"stockSymbol,omitempty"`
	Funds             string  `xml:"funds,omitempty"`
	ErrorMessage      string  `xml:"errorMessage,omitempty"`
}

const server = "transaction"
const schemaFile = "logging/schema.xsd"

type Logger interface {
	LogCommand(command Command, vars map[string]string)
	LogQuoteServ(stockQuote *models.StockQuote, trans string)
	LogTransaction(action string, username string, amount int, trans string)
	LogErrorEvent(command Command, vars map[string]string, emessage string)
}

type LogConnection struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
	Queue      amqp.Queue
}

func failOnError(err error, msg string) {
	if err != nil {
		utils.LogErr(err, msg)
		panic(err)
	}
}

func NewLoggerConnection() (logconn *LogConnection) {
	rabbitUser := os.Getenv("RABBITMQ_DEFAULT_USER")
	rabbitPass := os.Getenv("RABBITMQ_DEFAULT_PASS")
	rabbitHost := os.Getenv("RABBITMQ_HOST")
	rabbitPort := os.Getenv("RABBITMQ_PORT")
	url := fmt.Sprintf("amqp://%s:%s@%s:%s/", rabbitUser, rabbitPass, rabbitHost, rabbitPort)

	connection, err := amqp.Dial(url)
	logconn = &LogConnection{Connection: connection}

	failOnError(err, fmt.Sprintf("Failed to connect to Rabbit %s", url))

	logconn.Channel, err = logconn.Connection.Channel()
	failOnError(err, "Failed to open a channel")

	logconn.Queue, err = logconn.Channel.QueueDeclare(
		"log", // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")
	return
}

func DecodeMessage(reader io.Reader) (message *Message) {
	dec := gob.NewDecoder(reader)
	err := dec.Decode(&message)
	if err != nil {
		utils.LogErr(err, "Failed to decode message.")
	}
	return
}

func (logconn *LogConnection) publishMessage(message Message) {
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	err := enc.Encode(message)
	if err != nil {
		utils.LogErr(err, "Failed to encode message.")
	}

	err = logconn.Channel.Publish(
		"",                 // exchange
		logconn.Queue.Name, // routing key
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

func PrintMessage(message Message) {
	if message.UserCommand != nil {
		log.Printf("%+v\n", message.UserCommand)
	}
	if message.AccountTransaction != nil {
		log.Printf("%+v\n", message.AccountTransaction)
	}
	if message.SystemEvent != nil {
		log.Printf("%+v\n", message.SystemEvent)
	}
	if message.QuoteServer != nil {
		log.Printf("%+v\n", message.QuoteServer)
	}
	if message.ErrorEvent != nil {
		log.Printf("%+v\n", message.ErrorEvent)
	}
}

func formatStrAmount(amount string) (str string, err error) {
	b, err := strconv.Atoi(amount)
	if err != nil {
		return "", err
	}
	str = fmt.Sprintf("%d.%d", b/100, b%100)
	return
}

func formatAmount(amount int) string {
	return fmt.Sprintf("%d.%d", amount/100, amount%100)
}

func getUnixTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func parseTransactionNumber(transactionNumber string) (tnum int64) {
	tnum, err := strconv.ParseInt(transactionNumber, 10, 64)
	if err != nil {
		utils.LogErr(err, "Failed to parse transaction number")
	}
	return
}

func (logconn *LogConnection) LogCommand(command Command, vars map[string]string) {
	if _, exist := validCommands[command]; exist {
		timestamp := getUnixTimestamp()
		userCommand := UserCommandType{Timestamp: timestamp, Server: server, Command: command}

		if val, exist := vars["trans"]; exist {
			userCommand.TransactionNumber = parseTransactionNumber(val)
		}
		if val, exist := vars["username"]; exist {
			userCommand.Username = val
		}
		if val, exist := vars["symbol"]; exist {
			userCommand.Symbol = val
		}
		if val, exist := vars["filename"]; exist {
			userCommand.Filename = val
		}
		if val, exist := vars["amount"]; exist {
			var err error
			userCommand.Funds, err = formatStrAmount(val)
			if err != nil {
				utils.LogErr(err, "Failed to format amount")
				return
			}
		}

		msg := Message{UserCommand: &userCommand}
		logconn.publishMessage(msg)
	}
}

func (logconn *LogConnection) LogQuoteServ(stockQuote *models.StockQuote, trans string) {
	timestamp := getUnixTimestamp()
	quoteTimeInt, err := strconv.ParseInt(stockQuote.QuoteTimestamp, 10, 64)
	if err != nil {
		utils.LogErr(err, "Failed to parse quote server timestamp")
	}
	tnum := parseTransactionNumber(trans)

	quoteServer := QuoteServerType{Timestamp: timestamp,
		Server:            server,
		QuoteServerTime:   quoteTimeInt,
		Username:          stockQuote.Username,
		Symbol:            stockQuote.Symbol,
		Price:             stockQuote.Value,
		CryptoKey:         stockQuote.CrytpoKey,
		TransactionNumber: tnum}

	msg := Message{QuoteServer: &quoteServer}
	logconn.publishMessage(msg)
}

func (logconn *LogConnection) LogTransaction(action string, username string, amount int, trans string) {
	timestamp := getUnixTimestamp()
	tnum := parseTransactionNumber(trans)

	accountTransaction := AccountTransactionType{
		Timestamp:         timestamp,
		Server:            server,
		TransactionNumber: tnum,
		Username:          username,
		Action:            action,
		Funds:             formatAmount(amount),
	}

	msg := Message{AccountTransaction: &accountTransaction}
	logconn.publishMessage(msg)
}

// func LogSystemEvent(command string, username string, stocksymbol string, funds string) {

// 	file, err := os.OpenFile("log.xsd", os.O_APPEND|os.O_WRONLY, 0600)
// 	if err != nil {
// 		panic(err)
// 	}

// 	v := &SystemEvent{Timestamp: strconv.FormatInt(time.Now().UTC().UnixNano(), 10), Server: 1, Command: command, Username: username, StockSymbol: stocksymbol, Funds: funds}

// 	output, err := xml.MarshalIndent(v, "  ", "    ")

// 	if err != nil {

// 		fmt.Printf("error: %v\n", err)

// 	}

// 	file.Write(output)

// }

func (logconn *LogConnection) LogErrorEvent(command Command, vars map[string]string, emessage string) {
	timestamp := getUnixTimestamp()

	errorEvent := ErrorEventType{
		Timestamp:    timestamp,
		Server:       server,
		Command:      command,
		ErrorMessage: emessage}

	if val, exist := vars["trans"]; exist {
		errorEvent.TransactionNumber = parseTransactionNumber(val)
	}
	if val, exist := vars["username"]; exist {
		errorEvent.Username = val
	}
	if val, exist := vars["symbol"]; exist {
		errorEvent.Symbol = val
	}
	if val, exist := vars["amount"]; exist {
		var err error
		errorEvent.Funds, err = formatStrAmount(val)
		if err != nil {
			utils.LogErr(err, "Failed to format amount")
			return
		}
	}

	msg := Message{ErrorEvent: &errorEvent}
	logconn.publishMessage(msg)
}
