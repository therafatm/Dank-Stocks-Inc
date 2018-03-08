package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"common/logging"
	"common/utils"

	"github.com/lestrrat/go-libxml2"
	"github.com/lestrrat/go-libxml2/xsd"
)

const prefix = ""
const indent = "\t"
const logfile = "log.xml"
const schemaFile = "schema.xsd"

func XMLMarshalMessage(message logging.Message) (output []byte, err error) {
	if message.UserCommand != nil {
		output, err = xml.MarshalIndent(message.UserCommand, prefix, indent)
		if err != nil {
			return
		}
	}
	if message.AccountTransaction != nil {
		output, err = xml.MarshalIndent(message.AccountTransaction, prefix, indent)
		if err != nil {
			return
		}
	}
	if message.SystemEvent != nil {
		output, err = xml.MarshalIndent(message.SystemEvent, prefix, indent)
		if err != nil {
			return
		}
	}
	if message.QuoteServer != nil {
		output, err = xml.MarshalIndent(message.QuoteServer, prefix, indent)
		if err != nil {
			return
		}
	}
	if message.ErrorEvent != nil {
		output, err = xml.MarshalIndent(message.ErrorEvent, prefix, indent)
		if err != nil {
			return
		}
	}

	return
}

func validateSchema(s *xsd.Schema, ele []byte) {
	wrapper := []byte(fmt.Sprintf("<log>%s</log>", ele))

	d, err := libxml2.Parse(wrapper)
	if err != nil {
		utils.LogErr(err, "failed to parse XML")
		return
	}

	if err := s.Validate(d); err != nil {
		for _, err := range err.(xsd.SchemaValidationError).Errors() {
			if err != nil {
				utils.LogErr(err, "failed to validate XML.")
				return
			}
		}
	}
	if err != nil {
		utils.LogErr(err, "failed to validate XML.")
	}
}

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

	go func() {
		schema, err := os.Open(schemaFile)
		if err != nil {
			utils.LogErr(err, "failed to open file")
			return
		}
		defer schema.Close()

		schemabuf, err := ioutil.ReadAll(schema)
		if err != nil {
			utils.LogErr(err, "failed to read file")
			return
		}

		s, err := xsd.Parse(schemabuf)
		if err != nil {
			utils.LogErr(err, "failed to parse XSD")
			return
		}
		defer s.Free()

		f, err := os.Create(logfile)
		if err != nil {
			utils.LogErr(err, "Failed to open log file.")
		}
		defer f.Close()

		for d := range msgs {
			reader := bytes.NewReader(d.Body)
			message := logging.DecodeMessage(reader)
			output, err := XMLMarshalMessage(*message)
			if err != nil {
				utils.LogErr(err, "Failed to marshal message")
			}

			//validateSchema(s, output)

			_, err = f.Write(output)
			if err != nil {
				utils.LogErr(err, "Failed to write entry to file")
			}
			_, err = f.WriteString("\n")
			if err != nil {
				utils.LogErr(err, "Failed to write newline to file")
			}
		}
	}()

	for {
		log.Println("still alive")
		time.Sleep(10 * time.Second)
	}
}
