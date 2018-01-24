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

func checkBalance(e *httpexpect.Expect, username string, expected int) (obj *httpexpect.Object) {
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

func commitBuy(e *httpexpect.Expect, username string, status int) (obj *httpexpect.Object) {
	cmd := commands.Command{ Name: "COMMIT_BUY", Username: username}
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

	checkBalance(e, username, amount)

	//update
	addAmount :=  20023
	obj = add(e, username, addAmount, http.StatusOK)
	obj.Keys().ContainsOnly("id", "username", "money")
	obj.ValueEqual("username", username)
	obj.ValueEqual("money", amount + addAmount)

	checkBalance(e, username, amount + addAmount)

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
	cmd := commands.Command{ Name: "ADD", Username: username, Amount: amount }
	endpoint := commands.FormatCommandEndpoint(cmd)
	e.GET(endpoint).
		Expect().
		Status(http.StatusOK)

	// buy make sure remainder is discarded
	actualShares := sharesToBuy - 1
	obj := buy(e, username, testSymbol, amount - 1, http.StatusOK)
	obj.Keys().ContainsOnly("id", "username", "symbol", "type", "shares", "amount", "time")
	obj.ValueEqual("username", username)
	obj.ValueEqual("symbol", testSymbol)
	obj.ValueEqual("type", "BUY")
	obj.ValueEqual("shares", actualShares)
	obj.ValueEqual("amount", actualShares * testPrice)

	checkBalance(e, username, (sharesToBuy - actualShares) * testPrice)

	// commit buy order
	obj = commitBuy(e, username, http.StatusOK)
	obj.Keys().ContainsOnly("id", "username", "symbol", "shares")
	obj.ValueEqual("username", username)
	obj.ValueEqual("symbol", testSymbol)
	obj.ValueEqual("shares", actualShares)

	checkBalance(e, username, amount  - (testPrice * actualShares))

	// not enough money
	obj = buy(e, username, testSymbol, testPrice * sharesToBuy, http.StatusInternalServerError)
	obj.Keys().ContainsOnly("error", "message")

}


