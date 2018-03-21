package main 

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"common/logging"
	"common/utils"
	"logger_service/queries"

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
	db :=  queries.NewLogDBConnection(host, port)
	env := queries.Env{DB: db}
	env.DB.SetMaxOpenConns(300)


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

}

		