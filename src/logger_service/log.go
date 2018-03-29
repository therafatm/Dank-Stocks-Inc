package main

import (
	"bytes"
	"common/logging"
	"common/utils"
	"logger_service/queries"
	"time"
	"log"
	"os"
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

	host := os.Getenv("LOG_DB_HOST")
	port := os.Getenv("LOG_DB_PORT")
	db := queries.NewLogDBConnection(host, port)
	env := queries.Env{DB: db}

	go func() {
		buffer := map[string][][]interface{}{}
		writeTime := time.Now()

		for d := range msgs {
			reader := bytes.NewReader(d.Body)
			message := logging.DecodeMessage(reader)
			if message.DumpLog == nil{
				logging.PrintMessage(*message)
				buffer = queries.StoreMessage(buffer, *message)
				buffer, writeTime, err = env.CommitMessages(buffer, writeTime, false)
				if err != nil {
					utils.LogErr(err, "Failed to commit message")
				}
			}else {
				log.Println(len(buffer["usercommand"]))
				buffer, writeTime, err = env.CommitMessages(buffer, writeTime, true)
				if err != nil {
					utils.LogErr(err, "Failed to commit message")
				}
				log.Println("Dumping log.")
				dumplog.Dumplog(host, port, message.DumpLog.Filename, message.DumpLog.Username)
			}
		}
	}()

	for {
		time.Sleep(10 * time.Second)
	}
}
