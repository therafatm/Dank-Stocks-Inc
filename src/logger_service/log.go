package main

import (
	"bytes"
	"fmt"
	"common/logging"
	"common/utils"
	"time"
	"logger_service/queries"
)

func main() {
	logger := logging.NewLoggerConnection()
	msgs, err := logger.Channel.Consume(
		logger.Queue.Name, // queue
		"",                // consumer
		true,              // auto-ack
		false,             // exclusive
		false,             // no-local
		false,             // no-wait
		nil,               // args
	)
	if err != nil {
		utils.LogErr(err, "Failed to make consumer channel")
	}

	host := "logdb"
	port := "5432"
	db :=  queries.NewLogDBConnection(host, port)
	env := queries.Env{DB: db}
	env.DB.SetMaxOpenConns(300)


	go func() {
		fmt.Print("Logger service is running")

		for d := range msgs {
			fmt.Print("test")
			reader := bytes.NewReader(d.Body)
			message := logging.DecodeMessage(reader)
			logging.PrintMessage(*message)
			_, err := env.StoreMessage(*message)
			if err != nil {
				utils.LogErr(err, "Failed to store message")
			}
		}
	}()

	for {
		time.Sleep(10 * time.Second)
	}
}
