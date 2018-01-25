package test

import (
	"net/http"
	"testing"
	"fmt"
	"os"

	"test/workload_generator/commands"
	"github.com/gavv/httpexpect"
)


var (
	url string
	username string
)

const testPrice = 20000
const testSymbol = "TEST"

func init() {
   	host := os.Getenv("TRANS_HOST")
	port := os.Getenv("TRANS_PORT")
	url = fmt.Sprintf("http://%s:%s", host, port)
	username = "testuser"
}

func initTest(t *testing.T) (e *httpexpect.Expect){
	e = httpexpect.New(t, url)
	e.GET("/api/clearUsers").
		Expect().
		Status(http.StatusOK)
	return e
}

func checkAvailableBalance(e *httpexpect.Expect, username string, expected int) (obj *httpexpect.Object) {
	cmd := commands.Command{ Name: "BALANCE", Username: username }
	endpoint := commands.FormatCommandEndpoint(cmd)
	obj = e.GET(endpoint).
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	obj.Keys().ContainsOnly("balance")
	obj.ValueEqual("balance", expected)
	return
}

func checkAvailableShares(e *httpexpect.Expect, username string, symbol string, expected int) (obj *httpexpect.Object) {
	cmd := commands.Command{ Name: "SHARES", Username: username, Symbol: symbol }
	endpoint := commands.FormatCommandEndpoint(cmd)
	obj = e.GET(endpoint).
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	obj.Keys().ContainsOnly("shares")
	obj.ValueEqual("shares", expected)
	return
}


func add(e *httpexpect.Expect, username string, amount int,  status int) (obj *httpexpect.Object) {
	cmd := commands.Command{ Name: "ADD", Username: username, Amount: amount }
	endpoint := commands.FormatCommandEndpoint(cmd)
	obj = e.GET(endpoint).
		Expect().
		Status(status).
		JSON().Object()
	return
}

func buy(e *httpexpect.Expect, username string, symbol string, amount int, status int) (obj *httpexpect.Object) {
	cmd := commands.Command{ Name: "BUY", Username: username, Symbol: testSymbol, Amount: amount }
	endpoint := commands.FormatCommandEndpoint(cmd)
	obj = e.GET(endpoint).
		Expect().
		Status(status).
		JSON().Object()
	return
}

func sell(e *httpexpect.Expect, username string, symbol string, amount int, status int) (obj *httpexpect.Object) {
	cmd := commands.Command{ Name: "SELL", Username: username, Symbol: testSymbol, Amount: amount }
	endpoint := commands.FormatCommandEndpoint(cmd)
	obj = e.GET(endpoint).
		Expect().
		Status(status).
		JSON().Object()
	return
}

func commitBuy(e *httpexpect.Expect, username string, status int) (obj *httpexpect.Object) {
	cmd := commands.Command{ Name: "COMMIT_BUY", Username: username}
	endpoint := commands.FormatCommandEndpoint(cmd)
	obj = e.GET(endpoint).
		Expect().
		Status(status).
		JSON().Object()
	return
}

func cancelBuy(e *httpexpect.Expect, username string, status int) (obj *httpexpect.Object) {
	cmd := commands.Command{ Name: "CANCEL_BUY", Username: username}
	endpoint := commands.FormatCommandEndpoint(cmd)
	obj = e.GET(endpoint).
		Expect().
		Status(status).
		JSON().Object()
	return
}


func commitSell(e *httpexpect.Expect, username string, status int) (obj *httpexpect.Object) {
	cmd := commands.Command{ Name: "COMMIT_SELL", Username: username}
	endpoint := commands.FormatCommandEndpoint(cmd)
	obj = e.GET(endpoint).
		Expect().
		Status(status).
		JSON().Object()
	return
}

func cancelSell(e *httpexpect.Expect, username string, status int) (obj *httpexpect.Object) {
	cmd := commands.Command{ Name: "CANCEL_SELL", Username: username}
	endpoint := commands.FormatCommandEndpoint(cmd)
	obj = e.GET(endpoint).
		Expect().
		Status(status).
		JSON().Object()
	return
}

func TestAddUser(t *testing.T) {
	e := initTest(t)

	amount := 8786332

	//new user
	obj := add(e, username, amount, http.StatusOK)
	obj.Keys().ContainsOnly("id", "username", "money")
	obj.ValueEqual("username", username)
	obj.ValueEqual("money", amount)

	checkAvailableBalance(e, username, amount)

	//update
	addAmount :=  20023
	obj = add(e, username, addAmount, http.StatusOK)
	obj.Keys().ContainsOnly("id", "username", "money")
	obj.ValueEqual("username", username)
	obj.ValueEqual("money", amount + addAmount)

	checkAvailableBalance(e, username, amount + addAmount)

	//bad amount
	fAmount := 200.23
	endpoint := fmt.Sprintf("/api/add/%s/%f",username , fAmount)
	obj = e.GET(endpoint).
		Expect().
		Status(http.StatusInternalServerError).
		JSON().Object()

	obj.Keys().ContainsOnly("error", "message")
}


func TestBuy(t *testing.T) {
	e := initTest(t)

	// initial buy
	sharesToBuy := 10

	amount := testPrice * sharesToBuy 
	obj := add(e, username, amount, http.StatusOK)

	// buy make sure remainder is discarded
	actualShares := sharesToBuy - 1
	obj = buy(e, username, testSymbol, amount - 1, http.StatusOK)
	obj.Keys().ContainsOnly("id", "username", "symbol", "type", "shares", "amount", "time")
	obj.ValueEqual("username", username)
	obj.ValueEqual("symbol", testSymbol)
	obj.ValueEqual("type", "BUY")
	obj.ValueEqual("shares", actualShares)
	obj.ValueEqual("amount", actualShares * testPrice)

	checkAvailableBalance(e, username, (sharesToBuy - actualShares) * testPrice)

	// commit buy order
	obj = commitBuy(e, username, http.StatusOK)
	obj.Keys().ContainsOnly("id", "username", "symbol", "shares")
	obj.ValueEqual("username", username)
	obj.ValueEqual("symbol", testSymbol)
	obj.ValueEqual("shares", actualShares)

	checkAvailableBalance(e, username, amount  - (testPrice * actualShares))
	checkAvailableShares(e, username, testSymbol, actualShares)

	// not enough money
	obj = buy(e, username, testSymbol, testPrice * sharesToBuy, http.StatusInternalServerError)
	obj.Keys().ContainsOnly("error", "message")
}

func TestCancelBuy(t *testing.T) {
	e := initTest(t)

	// initial buy
	sharesForMoney := 100
	numShares := 3

	amount := testPrice * sharesForMoney
	obj := add(e, username, amount, http.StatusOK)

	obj = buy(e, username, testSymbol, 1 * testPrice, http.StatusOK)
	obj = buy(e, username, testSymbol, 2 * testPrice, http.StatusOK)
	obj = buy(e, username, testSymbol, numShares * testPrice, http.StatusOK)

	obj = cancelBuy(e, username, http.StatusOK)
	obj.Keys().ContainsOnly("id", "username", "symbol", "type", "shares", "amount", "time")
	obj.ValueEqual("username", username)
	obj.ValueEqual("symbol", testSymbol)
	obj.ValueEqual("shares", numShares)
	obj.ValueEqual("amount", numShares * testPrice)
	obj.ValueEqual("type", "BUY")


	obj = cancelBuy(e, username, http.StatusInternalServerError)
	obj.Keys().ContainsOnly("error", "message")

	checkAvailableBalance(e, username, amount)
	checkAvailableShares(e, username, testSymbol, 0)


}	

func TestSell(t *testing.T) {
	e := initTest(t)

	sharesToBuy := 10
	sharesToSell := sharesToBuy - 2

	amount := testPrice * sharesToBuy 
	add(e, username, amount, http.StatusOK)
	buy(e, username, testSymbol, amount, http.StatusOK)
	commitBuy(e, username, http.StatusOK)
	checkAvailableBalance(e, username, amount  - (testPrice * sharesToBuy))
	checkAvailableShares(e, username, testSymbol, sharesToBuy)

	obj := sell(e, username, testSymbol, sharesToSell * testPrice, http.StatusOK)
	obj.Keys().ContainsOnly("id", "username", "symbol", "type", "shares", "amount", "time")
	obj.ValueEqual("username", username)
	obj.ValueEqual("symbol", testSymbol)
	obj.ValueEqual("shares", sharesToSell)
	obj.ValueEqual("amount", sharesToSell * testPrice)
	obj.ValueEqual("type", "SELL")

	checkAvailableShares(e, username, testSymbol, sharesToBuy - sharesToSell)
	checkAvailableBalance(e, username, amount  - (testPrice * sharesToBuy))

	// commit sell order
	obj = commitSell(e, username, http.StatusOK)
	obj.Keys().ContainsOnly("id", "username", "symbol", "shares")
	obj.ValueEqual("username", username)
	obj.ValueEqual("symbol", testSymbol)
	obj.ValueEqual("shares", sharesToBuy - sharesToSell)

	checkAvailableShares(e, username, testSymbol, sharesToBuy - sharesToSell)
	checkAvailableBalance(e, username, amount - (testPrice * sharesToBuy) + (sharesToSell * testPrice) )

	// not enough shares
	obj = sell(e, username, testSymbol, testPrice * sharesToBuy, http.StatusInternalServerError)
	obj.Keys().ContainsOnly("error", "message")

	//sell remaining shares
	obj = sell(e, username, testSymbol, (sharesToBuy - sharesToSell) * testPrice, http.StatusOK)
	obj.Keys().ContainsOnly("id", "username", "symbol", "type", "shares", "amount", "time")
	obj.ValueEqual("username", username)
	obj.ValueEqual("symbol", testSymbol)
	obj.ValueEqual("shares", (sharesToBuy - sharesToSell))
	obj.ValueEqual("amount", (sharesToBuy - sharesToSell) * testPrice)
	obj.ValueEqual("type", "SELL")

	// commit sell order
	obj = commitSell(e, username, http.StatusOK)
	obj.Keys().ContainsOnly("id", "username", "symbol", "shares")
	obj.ValueEqual("username", username)
	obj.ValueEqual("symbol", testSymbol)
	obj.ValueEqual("shares", 0)

	checkAvailableShares(e, username, testSymbol, 0)
	checkAvailableBalance(e, username, amount)
}