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

func TestAddUser(t *testing.T) {
	e := httpexpect.New(t, url)
	e.GET("/api/clearUsers").
		Expect().
		Status(http.StatusOK)

	amount := 8786332

	//new user
	cmd := commands.Command{ Name: "ADD", Username: username, Amount: amount }
	endpoint := commands.FormatCommandEndpoint(cmd)
	obj := e.GET(endpoint).
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	obj.Keys().ContainsOnly("id", "username", "money")
	obj.ValueEqual("username", username)
	obj.ValueEqual("money", amount)


	//update
	addAmount :=  20023
	cmd = commands.Command{ Name: "ADD", Username: username, Amount: addAmount }
	endpoint = commands.FormatCommandEndpoint(cmd)
	obj = e.GET(endpoint).
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	obj.Keys().ContainsOnly("id", "username", "money")
	obj.ValueEqual("username", username)
	obj.ValueEqual("money", amount + addAmount)


	//bad amount
	fAmount := 200.23
	endpoint = fmt.Sprintf("/api/add/%s/%f",username , fAmount)
	obj = e.GET(endpoint).
		Expect().
		Status(http.StatusInternalServerError).
		JSON().Object()

	obj.Keys().ContainsOnly("error", "message")

}


func TestBuy(t *testing.T) {
	e := httpexpect.New(t, url)

	e.GET("/api/clearUsers").
		Expect().
		Status(http.StatusOK)


	// initial buy
	sharesToBuy := 12

	amount := testPrice * sharesToBuy 
	cmd := commands.Command{ Name: "ADD", Username: username, Amount: amount }
	endpoint := commands.FormatCommandEndpoint(cmd)
	e.GET(endpoint).
		Expect().
		Status(http.StatusOK)

	// buy make sure remainder is discarded
	sharesToBuy = 10
	remainder := testPrice - 1
	cmd = commands.Command{ Name: "BUY", Username: username, Symbol: testSymbol, Amount: testPrice * sharesToBuy + remainder }
	endpoint = commands.FormatCommandEndpoint(cmd)
	obj := e.GET(endpoint).
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	obj.Keys().ContainsOnly("id", "username", "symbol", "type", "shares", "amount", "time")
	obj.ValueEqual("username", username)
	obj.ValueEqual("symbol", testSymbol)
	obj.ValueEqual("type", "BUY")
	obj.ValueEqual("shares", sharesToBuy)
	obj.ValueEqual("amount", testPrice * sharesToBuy )

	// not enough money
	sharesToBuy = 13
	cmd = commands.Command{ Name: "BUY", Username: username, Symbol: testSymbol, Amount: testPrice * sharesToBuy }
	endpoint = commands.FormatCommandEndpoint(cmd)
	obj = e.GET(endpoint).
		Expect().
		Status(http.StatusInternalServerError).
		JSON().Object()

	obj.Keys().ContainsOnly("error", "message")


}