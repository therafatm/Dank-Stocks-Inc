package dumplog

import (
	"common/utils"
	"encoding/xml"
	"io/ioutil"
	"logger_service/queries"
	"log"
	"bufio"
	"os"
	"io"

	"github.com/lestrrat/go-libxml2"
	"github.com/lestrrat/go-libxml2/xsd"
)

const prefix = ""
const indent = "\t"
const logfile = "log.xml"
const schemaFile = "schema.xsd"

func validateSchema(s *xsd.Schema, fread io.Reader) {
	d, err := libxml2.ParseReader(fread)
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

func Dumplog(host string, port string, filename string, username string) {
	db := queries.NewLogDBConnection(host, port)
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

	f, err := os.Create(filename)
	if err != nil {
		utils.LogErr(err, "Failed to open log file.")
	}
	defer f.Close()

	f.Write([]byte("<log>\n"))

	cmds, err := env.QueryUserCommand()
	if err != nil {
		utils.LogErr(err, "Failed to query UserCommand")
	}

	for _, cmd := range cmds {
		output, err := xml.MarshalIndent(cmd, prefix, indent)
		if err != nil {
			utils.LogErr(err , "Failed to marshal UserCommand")
		}else{
			f.Write(output)
			f.Write([]byte("\n"))
		}
	}

	quotes, err := env.QueryQuoteServer()
	if err != nil {
		utils.LogErr(err, "Failed to query Quotes")
	}

	for _, quote := range quotes {
		output, err := xml.MarshalIndent(quote, prefix, indent)
		if err != nil {
			utils.LogErr(err , "Failed to marshal Quote")
		}else{
			f.Write(output)
			f.Write([]byte("\n"))
		}
	}

	f.Write([]byte("\n</log>"))

	_, err = f.Seek(0, 0)
	reader := bufio.NewReader(f)
	validateSchema(s, reader)
	log.Println("Log succesfully valdiated " + filename)
}