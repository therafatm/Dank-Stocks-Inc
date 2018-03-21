package queries

import (
	"database/sql"
	"fmt"
	"os"
	"common/utils"
	"common/logging"


	_ "github.com/lib/pq"
)

type Env struct {
	DB     *sql.DB
}

func NewLogDBConnection(host string, port string) (db *sql.DB) {
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := "logs"
	config := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", config)
	if err != nil {
		utils.LogErr(err, "Error connecting to DB.")
		panic(err)
	}
	return 
}

func (env Env) InsertUserCommand(data logging.UserCommandType) (res sql.Result, err error) {
	query := "INSERT INTO UserCommand(timestamp, server, transactionNum, command, username, stocksymbol, funds) VALUES($1,$2,$3,$4,$5,$6,$7)"
	res, err = env.DB.Exec(query, data.Timestamp, data.Server, data.TransactionNumber, data.Command, data.Username, data.Symbol, data.Funds)
	return
}

func (env Env) InsertAccountTransaction(data logging.AccountTransactionType) (res sql.Result, err error) {
	query := "INSERT INTO AccountTransaction(timestamp, server, transactionNum, command, username, stocksymbol, funds) VALUES($1,$2,$3,$4,$5,$6,$7)"
	res, err = env.DB.Exec(query, data.Timestamp, data.Server, data.TransactionNumber, data.Action, data.Username, data.Funds)
	return
}

func (env Env) InsertSystemEvent(data logging.SystemEventType) (res sql.Result, err error) {
	query := "INSERT INTO SystemEvent(timestamp, server, transactionNum, command, username, stocksymbol, funds) VALUES($1,$2,$3,$4,$5,$6,$7)"
	res, err = env.DB.Exec(query, data.Timestamp, data.Server, data.TransactionNumber, data.Command, data.Username,  data.Symbol, data.Funds)
	return
}

func (env Env) InsertQuoteServer(data logging.QuoteServerType) (res sql.Result, err error) {
	query := "INSERT INTO QuoteServer(timestamp, server, transactionNum, quoteServerTime, username, stocksymbol, money, cryptokey) VALUES($1,$2,$3,$4,$5,$6,$7,$8)"
	res, err = env.DB.Exec(query, data.Timestamp, data.Server, data.TransactionNumber, data.QuoteServerTime, data.Username,  data.Symbol, data.Price, data.CryptoKey)
	return
}

func (env Env) InsertErrorEvent(data logging.ErrorEventType) (res sql.Result, err error) {
	query := "INSERT INTO Errors(timestamp, server, transactionNum, command, username, funds, errorMessage) VALUES($1,$2,$3,$4,$5,$6,$7)"
	res, err = env.DB.Exec(query, data.Timestamp, data.Server, data.TransactionNumber, data.Command, data.Username, data.Funds, data.ErrorMessage)
	return
}

func (env Env) StoreMessage(message logging.Message) (result sql.Result, err error) {
	if message.UserCommand != nil {
		result, err = env.InsertUserCommand(*message.UserCommand)
		if err != nil {
			return 
		}
	}
	if message.AccountTransaction != nil {
		result, err = env.InsertAccountTransaction(*message.AccountTransaction)
		if err != nil {
			return
		}
	}
	if message.SystemEvent != nil {
		result, err = env.InsertSystemEvent(*message.SystemEvent)
		if err != nil {
			return
		}
	}
	if message.QuoteServer != nil {
		result, err = env.InsertQuoteServer(*message.QuoteServer)
		if err != nil {
			return
		}
	}
	if message.ErrorEvent != nil {
		result, err = env.InsertErrorEvent(*message.ErrorEvent)
		if err != nil {
			return
		}
	}

	return
}

