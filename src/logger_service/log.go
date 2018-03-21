package main

import (
	"bytes"
	"common/logging"
	"common/utils"
	"logger_service/queries"
	"time"
	"log"
	"logger_service/dumplog"
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
	db := queries.NewLogDBConnection(host, port)
	env := queries.Env{DB: db}
	env.DB.SetMaxOpenConns(300)

	go func() {
		for d := range msgs {
			reader := bytes.NewReader(d.Body)
			message := logging.DecodeMessage(reader)
			if message.DumpLog == nil{
				logging.PrintMessage(*message)
				_, err := env.StoreMessage(*message)
				if err != nil {
					utils.LogErr(err, "Failed to store message")
				}
			}else {
				log.Println("Dumping log.")
				dumplog.Dumplog(host, port, message.DumpLog.Filename, message.DumpLog.Username)
			}
		}
	}()

	for {
		time.Sleep(10 * time.Second)
	}
}
