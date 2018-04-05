package logging

import (
	"common/utils"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx"
)

type LogDB struct {
	DB *pgx.Conn
}

const num_records = 5000
const timeout = 10

const (
	USERCOMMAND        = "usercommand"
	ACCOUNTTRANSACTION = "accounttransaction"
	SYSTEMEVENT        = "systemevents"
	ERRORS             = "errors"
	QUOTESERVER        = "quoteserver"
)

var schema = map[string][]string{
	USERCOMMAND:        []string{"timestamp", "server", "transactionnum", "command", "username", "stocksymbol", "funds"},
	ACCOUNTTRANSACTION: []string{"timestamp", "server", "transactionnum", "action", "username", "funds"},
	SYSTEMEVENT:        []string{"timestamp", "server", "transactionnum", "command", "username", "stocksymbol", "funds"},
	ERRORS:             []string{"timestamp", "server", "transactionnum", "command", "username", "funds", "errormessage"},
	QUOTESERVER:        []string{"timestamp", "server", "transactionnum", "quoteservertime", "username", "stocksymbol", "money", "cryptokey"},
}

func NewLogDBConnection(host string, port string) (logdb LogDB) {
	user := os.Getenv("PGUSER")
	password := os.Getenv("PGPASSWORD")
	uport, err := strconv.ParseUint(port, 10, 16)
	if err != nil {
		utils.LogErr(err, "Error parsing port")
		panic(err)
	}
	u16port := uint16(uport)

	dbname := os.Getenv("LOG_DB")
	config := pgx.ConnConfig{
		Host:     host,
		Port:     u16port,
		Database: dbname,
		User:     user,
		Password: password,
	}

	db, err := pgx.Connect(config)
	if err != nil {
		utils.LogErr(err, "Error connecting to DB.")
		panic(err)
	}

	logdb = LogDB{DB: db}
	return
}

func (logdb LogDB) InsertUserCommand(data UserCommandType) (res pgx.CommandTag, err error) {
	query := "INSERT INTO UserCommand(timestamp, server, transactionNum, command, username, stocksymbol, funds) VALUES($1,$2,$3,$4,$5,$6,$7)"
	res, err = logdb.DB.Exec(query, data.Timestamp, data.Server, data.TransactionNumber, data.Command, data.Username, data.Symbol, data.Funds)
	return
}

func (logdb LogDB) InsertAccountTransaction(data AccountTransactionType) (res pgx.CommandTag, err error) {
	query := "INSERT INTO AccountTransaction(timestamp, server, transactionNum, action, username, funds) VALUES($1,$2,$3,$4,$5,$6)"
	res, err = logdb.DB.Exec(query, data.Timestamp, data.Server, data.TransactionNumber, data.Action, data.Username, data.Funds)
	return
}

func (logdb LogDB) InsertSystemEvent(data SystemEventType) (res pgx.CommandTag, err error) {
	query := "INSERT INTO SystemEvents(timestamp, server, transactionnum, command, username, stocksymbol, funds) VALUES($1,$2,$3,$4,$5,$6,$7)"
	res, err = logdb.DB.Exec(query, data.Timestamp, data.Server, data.TransactionNumber, data.Command, data.Username, data.Symbol, data.Funds)
	return
}

func (logdb LogDB) InsertQuoteServer(data QuoteServerType) (res pgx.CommandTag, err error) {
	query := "INSERT INTO QuoteServer(timestamp, server, transactionnum, quoteservertime, username, stocksymbol, money, cryptokey) VALUES($1,$2,$3,$4,$5,$6,$7,$8)"
	res, err = logdb.DB.Exec(query, data.Timestamp, data.Server, data.TransactionNumber, data.QuoteServerTime, data.Username, data.Symbol, data.Price, data.CryptoKey)
	return
}

func (logdb LogDB) InsertErrorEvent(data ErrorEventType) (res pgx.CommandTag, err error) {
	query := "INSERT INTO Errors(timestamp, server, transactionnum, command, username, funds, errorMessage) VALUES($1,$2,$3,$4,$5,$6,$7)"
	res, err = logdb.DB.Exec(query, data.Timestamp, data.Server, data.TransactionNumber, data.Command, data.Username, data.Funds, data.ErrorMessage)
	return
}

func (logdb LogDB) QueryUserCommand() (ret []UserCommandType, err error) {
	query := "SELECT timestamp, server, transactionnum, command, username, stocksymbol, funds FROM usercommand"
	rows, err := logdb.DB.Query(query)

	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		data := UserCommandType{}
		err = rows.Scan(&data.Timestamp, &data.Server, &data.TransactionNumber, &data.Command, &data.Username, &data.Symbol, &data.Funds)
		if err != nil {
			return
		}
		ret = append(ret, data)
	}
	return
}

func (logdb LogDB) QueryQuoteServer() (ret []QuoteServerType, err error) {
	query := "SELECT timestamp, server, transactionnum, quoteservertime, username, stocksymbol, money, cryptokey FROM quoteserver"
	rows, err := logdb.DB.Query(query)

	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		data := QuoteServerType{}
		err = rows.Scan(&data.Timestamp, &data.Server, &data.TransactionNumber, &data.QuoteServerTime, &data.Username, &data.Symbol, &data.Price, &data.CryptoKey)
		if err != nil {
			return
		}
		ret = append(ret, data)
	}
	return
}

func (logdb LogDB) QueryQuoteServer(user string) (ret []QuoteServerType, err error) {
	query := "SELECT timestamp, server, transactionnum, quoteservertime, username, stocksymbol, money, cryptokey FROM quoteserver"
	rows, err := logdb.DB.Query(query)

	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		data := QuoteServerType{}
		err = rows.Scan(&data.Timestamp, &data.Server, &data.TransactionNumber, &data.QuoteServerTime, &data.Username, &data.Symbol, &data.Price, &data.CryptoKey)
		if err != nil {
			return
		}
		ret = append(ret, data)
	}
	return
}

func ConvertUserCommand(data UserCommandType) (ret []interface{}) {
	ret = []interface{}{
		data.Timestamp,
		data.Server,
		data.TransactionNumber,
		data.Command,
		data.Username,
		data.Symbol,
		data.Funds,
	}
	return
}

func ConvertQuoteServer(data QuoteServerType) (ret []interface{}) {
	ret = []interface{}{
		data.Timestamp,
		data.Server,
		data.TransactionNumber,
		data.QuoteServerTime,
		data.Username,
		data.Symbol,
		data.Price,
		data.CryptoKey,
	}
	return
}

func ConvertAccountTransaction(data AccountTransactionType) (ret []interface{}) {
	ret = []interface{}{
		data.Timestamp,
		data.Server,
		data.TransactionNumber,
		data.Action,
		data.Username,
		data.Funds,
	}
	return
}

func ConvertSystemEvent(data SystemEventType) (ret []interface{}) {
	ret = []interface{}{
		data.Timestamp,
		data.Server,
		data.TransactionNumber,
		data.Command,
		data.Username,
		data.Symbol,
		data.Funds,
	}
	return
}

func ConvertErrorEvent(data ErrorEventType) (ret []interface{}) {
	ret = []interface{}{
		data.Timestamp,
		data.Server,
		data.TransactionNumber,
		data.Command,
		data.Username,
		data.Funds,
		data.ErrorMessage,
	}
	return
}

func (logdb LogDB) CommitMessages(buffer map[string][][]interface{}, writeTime time.Time, commitNow bool) (map[string][][]interface{}, time.Time, error) {
	curTime := time.Now()
	write := writeTime.Sub(curTime).Seconds() > timeout

	for k, _ := range buffer {
		if len(buffer[k])%num_records == 0 || write || commitNow {
			writeTime = curTime
			if len(buffer[k]) != 0 {
				_, err := logdb.DB.CopyFrom(
					pgx.Identifier{k},
					schema[k],
					pgx.CopyFromRows(buffer[k]),
				)
				if err != nil {
					return buffer, curTime, err
				}
			}
			buffer[k] = make([][]interface{}, 0)
		}
	}

	return buffer, writeTime, nil
}

func StoreMessage(buffer map[string][][]interface{}, message Message) map[string][][]interface{} {
	if message.UserCommand != nil {
		if _, ok := buffer[USERCOMMAND]; !ok {
			buffer[USERCOMMAND] = make([][]interface{}, 0)
		}

		buffer[USERCOMMAND] = append(buffer[USERCOMMAND], ConvertUserCommand(*message.UserCommand))
	}
	if message.AccountTransaction != nil {
		if _, ok := buffer[ACCOUNTTRANSACTION]; !ok {
			buffer[ACCOUNTTRANSACTION] = make([][]interface{}, 0)
		}

		buffer[ACCOUNTTRANSACTION] = append(buffer[ACCOUNTTRANSACTION], ConvertAccountTransaction(*message.AccountTransaction))
	}
	if message.SystemEvent != nil {
		if _, ok := buffer[SYSTEMEVENT]; !ok {
			buffer[SYSTEMEVENT] = make([][]interface{}, 0)
		}

		buffer[SYSTEMEVENT] = append(buffer[SYSTEMEVENT], ConvertSystemEvent(*message.SystemEvent))
	}
	if message.QuoteServer != nil {
		if _, ok := buffer[QUOTESERVER]; !ok {
			buffer[QUOTESERVER] = make([][]interface{}, 0)
		}

		buffer[QUOTESERVER] = append(buffer[QUOTESERVER], ConvertQuoteServer(*message.QuoteServer))
	}
	if message.ErrorEvent != nil {
		if _, ok := buffer[ERRORS]; !ok {
			buffer[ERRORS] = make([][]interface{}, 0)
		}
		buffer[ERRORS] = append(buffer[ERRORS], ConvertErrorEvent(*message.ErrorEvent))
	}

	return buffer
}
